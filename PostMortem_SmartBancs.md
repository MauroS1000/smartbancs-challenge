# Informe Post Mortem y Gestión de Incidentes TI

**Incidente:** Caída del sistema transaccional por bloqueos mutuos (Deadlocks) y alta latencia.
**Estado:** Resuelto en Producto Mínimo Viable (MVP).

## 1. Resumen del Incidente

Durante pruebas de estrés simuladas, el sistema experimentó una degradación severa del servicio. El tiempo de respuesta de la API superó los umbrales aceptables (SLA > 2s) y un porcentaje significativo de las transacciones fallaron debido a bloqueos a nivel de base de datos. 

## 2. Análisis de Causa Raíz

Se identificaron dos causas principales que originaron el incidente concurrente:

1.  **Condiciones de Carrera y Deadlocks:** Las transferencias concurrentes intentaban actualizar los saldos de las mismas cuentas en órdenes distintos. Al no existir un mecanismo de bloqueo jerárquico o un orden determinista en la base de datos, múltiples transacciones se bloqueaban entre sí esperando la liberación de recursos (Circular Wait), resultando en un *Deadlock* que abortaba la transacción.
2.  **Integración Síncrona de IA:** El llamado al modelo de Inteligencia Artificial para la evaluación de fraude se estaba realizando de manera síncrona dentro del flujo principal de la petición HTTP. El tiempo de inferencia del modelo retenía la conexión abierta, agotando el pool de conexiones y elevando la latencia general de la aplicación.

## 3. Resolución Implementada

Para mitigar la caída del sistema, se aplicaron las siguientes soluciones a nivel de código y base de datos:

*   **Bloqueo de Filas Ordenado (Row-Level Locking):** Se implementó el uso de `SELECT ... FOR UPDATE` en PostgreSQL dentro de una transacción serializable. Para evitar los *Deadlocks*, el sistema ahora ordena lexicográficamente los identificadores de las cuentas (origen y destino) antes de solicitar los bloqueos. Esto asegura que todas las transacciones concurrentes adquieran los bloqueos siempre en el mismo orden.
*   **Desacoplamiento Asíncrono (Goroutines):** Se extrajo la validación del modelo de IA del hilo principal de la petición. Utilizando el modelo de concurrencia de Go (*goroutines*), la evaluación de riesgo ahora se ejecuta en segundo plano (simulada con un retraso artificial), permitiendo que la respuesta HTTP se libere inmediatamente después de asegurar la transacción ACID.

## 4. Acciones Preventivas y de Escalamiento

Para asegurar la estabilidad en el entorno productivo a largo plazo, se definen las siguientes acciones:

*   **Instrumentación y Observabilidad:** Se ha implementado un `CorrelationID` por cada petición para trazar el ciclo de vida de la transacción a través de los logs y sistemas de monitoreo (como Prometheus y Grafana).
*   **Timeouts Estrictos:** Implementación de tiempos de espera máximos (timeouts) tanto en el contexto HTTP como en el *connection pool* de la base de datos para evitar que peticiones atascadas consuman todos los recursos del servidor.
*   **Alertas Automatizadas:** Configuración de alertas en el sistema de monitoreo que notifiquen al equipo de operaciones si el porcentaje de errores HTTP 5xx supera el 1% o si el percentil 95 de latencia (P95) sobrepasa los 1500 milisegundos.