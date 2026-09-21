# SmartBancs App - MVP Backend

Este repositorio contiene el Producto Mínimo Viable (MVP) del backend de SmartBancs. La solución expone un microservicio transaccional diseñado para soportar alta concurrencia, garantizando propiedades ACID e integrando la evaluación asíncrona de modelos de Inteligencia Artificial sin degradar el tiempo de respuesta.

## Prerrequisitos

Para ejecutar este proyecto, el entorno de despliegue debe contar con las siguientes herramientas instaladas:
* **Go** (v1.27 o superior)
* **Docker** y **Docker Compose**
* **Python** (v3.10 o superior, para la ejecución del pipeline ETL)

## Instalación y Configuración

**1. Configurar variables de entorno:**
Copie el archivo de entorno de ejemplo para establecer las credenciales locales de la base de datos.
```bash
cp .env.example .env
```

**2. Provisionar la infraestructura:**
Levante PostgreSQL y Redis de forma automatizada mediante el script de infraestructura como código.
```bash
docker-compose up -d
```

**3. Ejecutar transformación de datos (ETL):**
Ejecute el script que limpia y estructura los datos transaccionales heredados para su consumo por el modelo de Inteligencia Artificial.
```bash
python scripts/etl/transform.py
```

## Ejecución del Servidor

Inicie el microservicio inyectando las variables de entorno definidas.

Para terminales Bash/Zsh:
```bash
set -a; source .env; set +a; go run backend/cmd/api/main.go
```

Para terminales Fish:
```bash
env (cat .env) go run backend/cmd/api/main.go
```

## Pruebas de Validación

El API expone un endpoint para procesar transferencias. Ejecute el siguiente comando para validar una transacción exitosa. En los registros del servidor se evidenciará el rastreo mediante un identificador, la métrica de latencia y la ejecución asíncrona del servicio de Inteligencia Artificial.

```bash
curl -i -X POST http://localhost:8080/api/v1/transfer \
  -H "Content-Type: application/json" \
  -d '{"source_account_id": "acc-001", "target_account_id": "acc-002", "amount": 25.50}'
```

**Validación de Reglas de Negocio:**
* **Fondos insuficientes:** Envíe un monto superior al saldo para validar el rechazo de la base de datos transaccional mediante un código HTTP 422 Unprocessable Entity.
* **Monto inválido:** Envíe un valor negativo para recibir un rechazo en la capa de negocio mediante un código HTTP 400 Bad Request.

## Operaciones de Apagado

Para detener la ejecución del servidor presione `Ctrl + C`. Para detener y destruir los contenedores de infraestructura ejecute:
```bash
docker-compose down
```
## Evidencias

* **Video Demostrativo:** [https://youtu.be/qIKL8xEzvUs]