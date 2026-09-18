/*
Package model define las entidades principales del dominio y los objetos de transferencia
de datos (DTOs) de la aplicación SmartBancs. Provee las estructuras necesarias para el mapeo
objeto-relacional y la serialización/deserialización de payloads JSON en la capa HTTP.
*/
package model

import "time"

// Account representa una cuenta bancaria dentro del sistema.
// Mantiene el estado actual del saldo, el cual es modificado bajo estrictos
// controles de concurrencia y bloqueos pesimistas.
type Account struct {
	AccountID string    `json:"account_id"`
	OwnerName string    `json:"owner_name"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Transaction representa un registro inmutable de un movimiento financiero.
// Además de funcionar como el libro mayor (ledger) de la base de datos, actúa
// como la tabla base para el patrón Transactional Outbox, permitiendo la sincronización
// diferida con el sistema central heredado (Bancs).
type Transaction struct {
	TransactionID   string    `json:"transaction_id"`
	CorrelationID   string    `json:"correlation_id"`
	SourceAccountID string    `json:"source_account_id"`
	TargetAccountID string    `json:"target_account_id"`
	Amount          float64   `json:"amount"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
}

// TransferRequest define la estructura del payload JSON esperado del cliente
// para iniciar una solicitud de transferencia de fondos.
type TransferRequest struct {
	SourceAccountID string  `json:"source_account_id"`
	TargetAccountID string  `json:"target_account_id"`
	Amount          float64 `json:"amount"`
}

// TransferResponse define la estructura del payload JSON devuelto al cliente
// una vez que la transferencia ha sido procesada y consolidada en la base de datos.
type TransferResponse struct {
	TransactionID string  `json:"transaction_id"`
	CorrelationID string  `json:"correlation_id"`
	Status        string  `json:"status"`
	Amount        float64 `json:"amount"`
}