/*
Package handler expone la API REST de la aplicación SmartBancs.
Actúa como la capa de adaptadores y transporte: decodifica las peticiones JSON entrantes,
delega la ejecución transaccional a la capa de servicio y traduce los errores de dominio
(ej. cuentas no encontradas, fondos insuficientes) a códigos de estado HTTP estandarizados.
*/
package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
	"smartbancs-backend/internal/model"
	"smartbancs-backend/internal/service"
)

// TransferHandler estructura el controlador HTTP inyectando la dependencia
// del servicio transaccional.
type TransferHandler struct {
	transferService *service.TransferService
}

// NewTransferHandler inicializa y retorna una nueva instancia del controlador de transferencias.
func NewTransferHandler(transferService *service.TransferService) *TransferHandler {
	return &TransferHandler{
		transferService: transferService,
	}
}

// HandleTransfer procesa la petición de transferencia financiera.
// Valida el método HTTP, gestiona el CorrelationID para observabilidad transaccional,
// decodifica el payload y mapea los resultados o errores del dominio hacia el cliente HTTP
// de forma estandarizada.
func (h *TransferHandler) HandleTransfer(w http.ResponseWriter, r *http.Request) {
	// 1. Validar verbo HTTP
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "método no permitido"})
		return
	}

	// 2. Extraer o generar el Correlation ID para observabilidad
	// Si el cliente no lo provee, se genera uno para asegurar el rastreo de inicio a fin.
	correlationID := r.Header.Get("X-Correlation-ID")
	if correlationID == "" {
		correlationID = uuid.NewString()
	}
	w.Header().Set("X-Correlation-ID", correlationID)

	// 3. Decodificar el payload JSON de la solicitud
	var req model.TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "payload JSON inválido"})
		return
	}

	// 4. Delegar la transacción a la capa de servicio
	res, err := h.transferService.ExecuteTransfer(r.Context(), req, correlationID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidAmount), errors.Is(err, service.ErrSameAccount):
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		case errors.Is(err, service.ErrAccountNotFound):
			writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		case errors.Is(err, service.ErrInsufficientFunds):
			writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		default:
			// Error no esperado, se registra en consola con su CorrelationID para auditoría
			log.Printf("ERROR transfer [%s]: %v", correlationID, err)
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return
	}

	// 5. Responder con éxito
	writeJSON(w, http.StatusOK, res)
}

// writeJSON es una función auxiliar para unificar la emisión de respuestas HTTP en formato JSON,
// asegurando que los encabezados y códigos de estado se configuren correctamente.
func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}