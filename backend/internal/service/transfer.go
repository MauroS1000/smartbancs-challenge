/*
Package service concentra la lógica de negocio central de la aplicación.
Este archivo en particular gestiona el ciclo de vida de las transferencias financieras,
garantizando la integridad transaccional (ACID), previniendo condiciones de carrera
y mitigando bloqueos mutuos (deadlocks) en entornos de alta concurrencia.
*/
package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"smartbancs-backend/internal/model"
)

// Errores de dominio predefinidos para estandarizar las respuestas de la capa de servicio
// y evitar la exposición de detalles internos de la base de datos a capas superiores.
var (
	ErrInvalidAmount     = errors.New("el monto debe ser mayor a 0")
	ErrSameAccount       = errors.New("la cuenta de origen y destino no pueden ser iguales")
	ErrAccountNotFound   = errors.New("cuenta no encontrada")
	ErrInsufficientFunds = errors.New("fondos insuficientes en la cuenta origen")
)

// TransferService orquesta la ejecución de transferencias.
// Mantiene las dependencias necesarias: el pool de conexiones a la base de datos
// para persistencia transaccional y el cliente de IA para evaluación de fraude.
type TransferService struct {
	pool     *pgxpool.Pool
	aiClient *AIFraudClient
}

// NewTransferService es el constructor que inyecta las dependencias requeridas
// para instanciar el servicio de transferencias.
func NewTransferService(pool *pgxpool.Pool, aiClient *AIFraudClient) *TransferService {
	return &TransferService{
		pool:     pool,
		aiClient: aiClient,
	}
}

// ExecuteTransfer procesa una transferencia de fondos entre dos cuentas garantizando propiedades ACID.
// Implementa un ordenamiento determinista de bloqueos para prevenir deadlocks y utiliza
// bloqueos pesimistas (SELECT FOR UPDATE) para evitar condiciones de carrera.
func (s *TransferService) ExecuteTransfer(ctx context.Context, req model.TransferRequest, correlationID string) (*model.TransferResponse, error) {
	// 1. Validaciones básicas de negocio
	if req.Amount <= 0 {
		return nil, ErrInvalidAmount
	}
	if req.SourceAccountID == req.TargetAccountID {
		return nil, ErrSameAccount
	}

	// 2. Iniciar transacción en la base de datos
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al iniciar transaccion: %w", err)
	}
	defer tx.Rollback(ctx)

	// 3. Prevención de Deadlocks (Abrazo mortal)
	// Se ordenan los IDs de las cuentas lexicográficamente para asegurar que
	// siempre se bloqueen en el mismo orden transaccional, independientemente de
	// qué cuenta sea origen o destino en la petición concurrente.
	firstLockID, secondLockID := req.SourceAccountID, req.TargetAccountID
	if firstLockID > secondLockID {
		firstLockID, secondLockID = secondLockID, firstLockID
	}

	var dummyBalance float64
	queryLock := "SELECT balance FROM accounts WHERE account_id = $1 FOR UPDATE"

	// 4. Bloqueo pesimista de la primera cuenta
	if err := tx.QueryRow(ctx, queryLock, firstLockID).Scan(&dummyBalance); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %s", ErrAccountNotFound, firstLockID)
		}
		return nil, fmt.Errorf("error al bloquear cuenta %s: %w", firstLockID, err)
	}

	// 5. Bloqueo pesimista de la segunda cuenta
	if err := tx.QueryRow(ctx, queryLock, secondLockID).Scan(&dummyBalance); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %s", ErrAccountNotFound, secondLockID)
		}
		return nil, fmt.Errorf("error al bloquear cuenta %s: %w", secondLockID, err)
	}

	// 6. Verificación de fondos de la cuenta origen
	var sourceBalance float64
	err = tx.QueryRow(ctx, "SELECT balance FROM accounts WHERE account_id = $1", req.SourceAccountID).Scan(&sourceBalance)
	if err != nil {
		return nil, fmt.Errorf("error al verificar saldo de origen: %w", err)
	}

	if sourceBalance < req.Amount {
		return nil, ErrInsufficientFunds
	}

	// 7. Ejecución atómica de movimientos contables
	queryDebit := "UPDATE accounts SET balance = balance - $1, updated_at = NOW() WHERE account_id = $2"
	if _, err := tx.Exec(ctx, queryDebit, req.Amount, req.SourceAccountID); err != nil {
		return nil, fmt.Errorf("error debitando de cuenta origen: %w", err)
	}

	queryCredit := "UPDATE accounts SET balance = balance + $1, updated_at = NOW() WHERE account_id = $2"
	if _, err := tx.Exec(ctx, queryCredit, req.Amount, req.TargetAccountID); err != nil {
		return nil, fmt.Errorf("error acreditando a cuenta destino: %w", err)
	}

	// 8. Registro de la transacción
	// Este registro servirá como base para el patrón Transactional Outbox,
	// permitiendo la sincronización diferida con el sistema heredado.
	transactionID := uuid.NewString()
	queryInsertTx := `
		INSERT INTO transactions (transaction_id, correlation_id, source_account_id, destination_account_id, amount, status)
		VALUES ($1, $2, $3, $4, $5, 'COMPLETED')
	`
	_, err = tx.Exec(ctx, queryInsertTx, transactionID, correlationID, req.SourceAccountID, req.TargetAccountID, req.Amount)
	if err != nil {
		return nil, fmt.Errorf("error registrando transaccion contable: %w", err)
	}

	// 9. Consolidación de la transacción en disco
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("error en commit transaccional: %w", err)
	}

	// 10. Invocación asíncrona del modelo de IA
	// Se delega a una goroutine utilizando context.Background() debido a que el contexto HTTP (ctx)
	// original se cancelará en cuanto se envíe la respuesta al cliente.
	// Esto asegura que la latencia de la IA no degrade el tiempo de respuesta del API.
	go s.aiClient.AnalyzeTransaction(context.Background(), transactionID, req.SourceAccountID, req.TargetAccountID, req.Amount)

	return &model.TransferResponse{
		TransactionID: transactionID,
		CorrelationID: correlationID,
		Status:        "COMPLETED",
		Amount:        req.Amount,
	}, nil
}