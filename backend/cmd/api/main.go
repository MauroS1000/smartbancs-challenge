/*
Package main es el punto de entrada de la aplicación SmartBancs.
Orquesta la inicialización de la infraestructura subyacente (conexiones a PostgreSQL),
configura la inyección de dependencias entre las capas (Service -> Handler),
e inicializa el servidor HTTP con protección de timeouts.
Implementa un apagado ordenado (Graceful Shutdown) para asegurar que las transacciones
en vuelo no se interrumpan abruptamente al recibir señales del sistema operativo.
*/
package main

import (
        "context"
        "errors"
        "log"
        "net/http"
        "os"
        "os/signal"
        "syscall"
        "time"

        "smartbancs-backend/internal/config"
        "smartbancs-backend/internal/db"
        "smartbancs-backend/internal/handler"
        "smartbancs-backend/internal/service"
)

func main() {
        // 1. Cargar variables de entorno y configuración base
        cfg := config.Load()

        // 2. Inicializar pool de conexiones a PostgreSQL
        // Se define un timeout estricto para evitar bloqueos en el arranque si la base de datos no responde.
        ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()

        pool, err := db.NewPostgresPool(ctx, cfg.DatabaseURL)
        if err != nil {
                log.Fatalf("Fallo crítico al conectar con la base de datos: %v", err)
        }
        defer pool.Close()
        log.Println("Pool de conexiones a PostgreSQL establecido exitosamente")

        // 3. Inicializar capas (Inyección de dependencias)
        // Se instancian los componentes asegurando el aislamiento de responsabilidades.
        aiClient := service.NewAIFraudClient()
        transferService := service.NewTransferService(pool, aiClient)
        transferHandler := handler.NewTransferHandler(transferService)

        // 4. Configurar enrutador HTTP nativo
        mux := http.NewServeMux()

        // Se aplica el middleware de observabilidad al endpoint transaccional para cumplir
        // con el registro automatizado de métricas de latencia y estado (Sección 3.4 del reto).
        mux.HandleFunc("/api/v1/transfer", handler.ObservabilityMiddleware(transferHandler.HandleTransfer))

        // Endpoint de validación de estado para orquestadores (ej. Kubernetes, Docker Compose).
        mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
                w.WriteHeader(http.StatusOK)
                _, _ = w.Write([]byte(`{"status":"UP"}`))
        })

        // 5. Configurar servidor HTTP con timeouts de protección
        // Previene ataques de denegación de servicio (Slowloris) y agotamiento de descriptores de archivo.
        server := &http.Server{
                Addr:         ":" + cfg.Port,
                Handler:      mux,
                ReadTimeout:  5 * time.Second,
                WriteTimeout: 10 * time.Second,
                IdleTimeout:  15 * time.Second,
        }

        // 6. Arrancar servidor en una goroutine secundaria
        // Permite que el hilo principal continúe hacia el bloqueo del canal de señales de apagado.
        go func() {
                log.Printf("Servidor bancario escuchando en el puerto :%s", cfg.Port)
                if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
                        log.Fatalf("Error en el servidor HTTP: %v", err)
                }
        }()

        // 7. Esperar señales de apagado del sistema operativo (SIGINT, SIGTERM)
        shutdownSignal := make(chan os.Signal, 1)
        signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGTERM)
        <-shutdownSignal

        log.Println("Señal de terminación recibida. Iniciando Graceful Shutdown...")

        // 8. Contexto con límite de 5 segundos para drenar peticiones activas
        // Garantiza que las transacciones en procesamiento finalicen correctamente en disco antes de cerrar el proceso.
        shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer shutdownCancel()

        if err := server.Shutdown(shutdownCtx); err != nil {
                log.Printf("Fallo durante el apagado del servidor: %v", err)
        }

        log.Println("Servidor apagado de forma segura.")
}