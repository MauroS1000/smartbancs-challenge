// Módulo principal del backend de SmartBancs.
// Gestiona las dependencias del proyecto garantizando construcciones reproducibles
// y un control estricto sobre las librerías de terceros utilizadas en el MVP.
module smartbancs-backend

// Versión de Go requerida para la compilación del proyecto.
go 1.27.1

require (
	// Utilizado para la generación de identificadores únicos universales, 
	// específicamente para el TransactionID y el manejo del CorrelationID 
	// en el middleware de observabilidad.
	github.com/google/uuid v1.6.0 // indirect

	// Dependencias internas del driver pgx para la lectura de archivos .pgpass y configuración.
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect

	// Driver y toolkit PostgreSQL de alto rendimiento (pgx).
	// Seleccionado por su soporte avanzado y nativo para Connection Pooling (puddle)
	// en lugar del paquete estándar 'database/sql'. Es fundamental para soportar 
	// las exigencias de alta concurrencia transaccional sin degradación de servicio.
	github.com/jackc/pgx/v5 v5.11.0 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect

	// Librerías de soporte subyacentes requeridas para primitivas de concurrencia 
	// y manipulación de codificaciones de texto por el driver de base de datos.
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/text v0.29.0 // indirect
)