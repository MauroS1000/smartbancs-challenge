/*
Package config centraliza la carga y validación de las variables de entorno de la aplicación.
Se encarga de construir las cadenas de conexión (DSN) para la infraestructura de persistencia
y de asegurar que los parámetros críticos y secretos estén presentes antes de permitir
el arranque del servidor.
*/
package config

import (
	"fmt"
	"log"
	"os"
)

// Config almacena los parámetros de configuración globales del sistema.
type Config struct {
	Port        string
	DatabaseURL string
	RedisURL    string
}

// Load lee el entorno de ejecución, aplica valores por defecto para facilitar
// el despliegue en entornos de desarrollo local y aborta la ejecución (Fatal)
// si faltan variables obligatorias, garantizando un enfoque Fail-Fast.
func Load() *Config {
	port := getEnv("PORT", "8080")
	dbHost := getEnv("POSTGRES_HOST", "localhost")
	dbPort := getEnv("POSTGRES_PORT", "5432")
	dbUser := getEnv("POSTGRES_USER", "bancs_user")
	dbName := getEnv("POSTGRES_DB", "smartbancs")
	dbPass := getRequiredEnv("POSTGRES_PASSWORD")

	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPass, dbHost, dbPort, dbName)

	redisURL := fmt.Sprintf("%s:%s", redisHost, redisPort)

	return &Config{
		Port:        port,
		DatabaseURL: dbURL,
		RedisURL:    redisURL,
	}
}

// getEnv recupera el valor de una variable de entorno o retorna el valor
// de respaldo (fallback) especificado si la variable no está definida o está vacía.
func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}

// getRequiredEnv recupera una variable de entorno obligatoria.
// Si la variable no existe o está vacía, finaliza el proceso inmediatamente
// para prevenir inconsistencias en la conexión con la infraestructura.
func getRequiredEnv(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		log.Fatalf("FATAL: variable de entorno requerida no encontrada: %s", key)
	}
	return val
}