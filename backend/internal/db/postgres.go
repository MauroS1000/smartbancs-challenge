/*
Package db proporciona la capa de gestión de conexiones e infraestructura de persistencia.
Implementa el manejo eficiente de recursos mediante connection pooling (agrupación de conexiones)
para soportar el procesamiento de alta concurrencia transaccional sin agotar los descriptores
de archivo ni la memoria del motor de base de datos.
*/
package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NewPostgresPool inicializa y afina un pool de conexiones optimizado para PostgreSQL.
// Define límites máximos y mínimos de concurrencia, junto con políticas de expiración
// de conexiones inactivas para prevenir fugas de recursos (connection leaks).
// Implementa una verificación estricta de conectividad (Ping) con timeout para asegurar
// que la infraestructura subyacente está operativa bajo un enfoque de falla rápida (Fail-Fast).
func NewPostgresPool(ctx context.Context, connString string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("error al parsear config de base de datos: %w", err)
	}

	// Configuración ajustada para mantener un flujo constante bajo picos de concurrencia
	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 5 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("error al inicializar pool de pgx: %w", err)
	}

	// Validación de red y credenciales antes de ceder el control a la aplicación
	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("no se pudo conectar a postgres: %w", err)
	}

	return pool, nil
}