# Documento Técnico de Arquitectura: SmartBancs App

## 1. Diseño de Arquitectura y Decisiones Técnicas

El diseño de la plataforma SmartBancs se fundamenta en una arquitectura de microservicios orientada a eventos para soportar picos transaccionales de hasta 10 000 peticiones por segundo, garantizando tiempos de respuesta inferiores a 2 segundos.

**Decisiones Tecnológicas:**
*   **Backend (Go):** Se seleccionó Go por su modelo de concurrencia nativa (*goroutines*). Esto permite manejar miles de conexiones simultáneas con un consumo mínimo de memoria y ejecutar la integración con el modelo de IA de forma asíncrona y no bloqueante.
*   **Base de Datos Relacional (PostgreSQL):** Utilizada como repositorio transaccional principal debido a su robusto cumplimiento de las propiedades ACID. Se implementó un control de concurrencia mediante bloqueos a nivel de fila (`SELECT ... FOR UPDATE`) para prevenir condiciones de carrera durante la actualización de saldos.
*   **Caché en Memoria (Redis):** Integrado como capa de almacenamiento temporal rápida. Previene la sobrecarga de la base de datos primaria y actúa como mecanismo base para implementar limitadores de tasa (*Rate Limiting*) en futuras interacciones con sistemas de terceros.

## 2. Estrategia de Sincronización con el Core Legado (Bancs)

El sistema central "Bancs" es un componente transaccional heredado que sufre degradación de rendimiento ante un alto volumen de consultas directas. Para integrar la nueva aplicación sin saturar este sistema, se define la siguiente estrategia de sincronización asíncrona:

*   **Patrón Transaccional Outbox:** Las transacciones procesadas exitosamente en SmartBancs se registran en una tabla `outbox` dentro de PostgreSQL en la misma transacción local (ACID).
*   **Gestor de Colas (Message Broker):** Un proceso *Relay* lee continuamente la tabla `outbox` y publica los eventos transaccionales en un intermediario de mensajes (como Apache Kafka o RabbitMQ).
*   **Consumidor con Limitador de Tasa (Rate-limited Consumer):** Un microservicio dedicado consume los mensajes de la cola y los envía al core "Bancs" de manera secuencial o en lotes (Micro-batching) aplicando políticas estrictas de limitación de tasa. Esto garantiza que "Bancs" reciba las actualizaciones de saldo a un ritmo predecible y controlable, protegiendo su disponibilidad.

## 3. Manejo y Ciclo de Vida del Modelo de Inteligencia Artificial (MLOps)

El servicio de recomendaciones financieras y prevención de fraude operará bajo un ciclo de vida estructurado para asegurar su precisión y disponibilidad en producción:

*   **Ejecución No Bloqueante:** La inferencia del modelo se ejecuta en un proceso en segundo plano (asíncrono) desacoplado del flujo crítico transaccional. Esto asegura que la latencia de la IA no afecte el tiempo de procesamiento de la transferencia principal.
*   **Alimentación con Nuevos Datos:** Los datos transaccionales crudos son extraídos periódicamente y procesados mediante canalizaciones ETL/ELT (como la implementada en Python) para limpiar valores nulos y estandarizar formatos. Estos datos estructurados se almacenan en un *Data Lake* para el reentrenamiento programado del modelo.
*   **Monitoreo de Data Drift:** Se implementará un análisis estadístico continuo (por ejemplo, métricas de distancia de Wasserstein o pruebas KS) comparando la distribución de los datos de entrenamiento con los datos entrantes en tiempo real. Si se detecta una desviación significativa (Data Drift), el sistema emitirá una alerta automática al equipo de ciencia de datos para iniciar el proceso de recalibración.
*   **Gestión de Consumo de Recursos:** El servicio de inferencia de IA se desplegará en clústeres elásticos (Kubernetes) con escalamiento automático de pods (HPA). El escalamiento estará condicionado por métricas de saturación de CPU/GPU y longitud de la cola de peticiones, optimizando los costos de infraestructura computacional.