/*
Package handler contiene la capa de adaptadores y transporte HTTP.
Este archivo provee los componentes de telemetría e instrumentación necesarios para cumplir
con los requerimientos de observabilidad (Sección 3.4). Intercepta el tráfico entrante
para generar métricas operativas precisas que facilitan el diagnóstico de incidentes.
*/
package handler

import (
	"log"
	"net/http"
	"time"
)

// statusWriter es un decorador para http.ResponseWriter que permite interceptar
// y capturar el código de estado HTTP (ej. 200, 422, 500) emitido por el controlador.
// Esto es necesario porque la interfaz original no expone un método para leer el estado
// una vez que ha sido escrito.
type statusWriter struct {
	http.ResponseWriter
	status int
}

// WriteHeader sobrescribe el método original para registrar el código de estado
// en memoria antes de enviarlo al cliente.
func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// ObservabilityMiddleware envuelve los controladores HTTP para registrar métricas
// de volumen transaccional, códigos de error y latencia con precisión de microsegundos.
// Extrae el CorrelationID de las cabeceras para permitir el rastreo (tracing)
// de la transacción a lo largo de todos los componentes del sistema.
func ObservabilityMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next(sw, r)

		duration := time.Since(start)
		correlationID := sw.Header().Get("X-Correlation-ID")

		// Registro estructurado de operaciones críticas para el monitoreo del rendimiento
		// y detección temprana de degradación del servicio o cuellos de botella.
		log.Printf("[Métrica HTTP] Método: %s | Ruta: %s | Estado: %d | Latencia: %v | CorrelationID: %s",
			r.Method, r.URL.Path, sw.status, duration, correlationID)
	}
}