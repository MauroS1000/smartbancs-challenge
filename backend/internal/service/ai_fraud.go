/*
Package service contiene la lógica de negocio central y la integración con sistemas externos.
Este archivo provee AIFraudClient, un mock funcional del sistema de prevención de fraude
basado en modelos de Machine Learning, cumpliendo con el requerimiento de evaluación de IA (Sección 3.3).
*/
package service

import (
	"context"
	"log"
	"math/rand"
	"time"
)

// AIFraudClient representa el cliente de integración con el servicio externo
// de Inteligencia Artificial para la detección de anomalías y fraudes.
type AIFraudClient struct{}

// NewAIFraudClient inicializa y retorna una nueva instancia del mock del cliente de IA.
func NewAIFraudClient() *AIFraudClient {
	return &AIFraudClient{}
}

// AnalyzeTransaction simula el consumo de un modelo predictivo de Machine Learning.
// Introduce una latencia artificial (500ms - 1500ms) para emular el tiempo de inferencia y red.
// CRÍTICO: Este método debe ser invocado obligatoriamente de forma asíncrona (mediante Goroutines)
// para garantizar que la evaluación de riesgo no bloquee la ejecución de la transferencia principal
// ni degrade el tiempo de respuesta exigido (SLA < 2 segundos).
func (ai *AIFraudClient) AnalyzeTransaction(ctx context.Context, transactionID, source, target string, amount float64) {
	// Simulación de latencia de inferencia (500ms - 1500ms)
	time.Sleep(time.Duration(rand.Intn(1000)+500) * time.Millisecond)

	score := rand.Float64() // Genera un valor entre 0.0 y 1.0
	riskLevel := "LOW"
	if score > 0.85 {
		riskLevel = "HIGH"
	} else if score > 0.5 {
		riskLevel = "MEDIUM"
	}

	// Registro de observabilidad del modelo de IA, emitiendo métricas estructuradas
	// para facilitar el monitoreo de las decisiones algorítmicas en producción.
	log.Printf("[IA Observabilidad] TxID: %s | Riesgo: %s (Score: %.2f) | Origen: %s -> Destino: %s",
		transactionID, riskLevel, score, source, target)
}