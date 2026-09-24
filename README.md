
# SmartBancs App - MVP Backend

## Instrucciones de Ejecución

Este proyecto incluye un `Makefile` en la raíz del repositorio para estandarizar y simplificar el despliegue del entorno en cualquier sistema operativo.

### Prerrequisitos
* Docker y Docker Compose
* Go 1.20 o superior
* Python 3.x
* Make (GNU Make)

### Pasos operativos

**1. Desplegar la infraestructura (PostgreSQL y Redis)**
```bash
make up
```

**2. Iniciar el microservicio backend**
```bash
make run
```
*El servidor inyectará automáticamente las variables de entorno locales y escuchará en el puerto 8080.*

**3. Ejecutar el pipeline de transformación de datos (ETL)**
```bash
make etl
```
*El script procesará los registros heredados y generará el archivo `ai_ready_transactions.json` en el directorio de datos.*

**4. Detener y limpiar el entorno**
```bash
make down
```

## Pruebas de la API

Para verificar el funcionamiento del microservicio, ejecute los siguientes comandos en una nueva terminal mientras el servidor se encuentra en ejecución:

**Escenario 1: Transacción exitosa**
```bash
curl -X POST http://localhost:8080/api/v1/transferencias \
  -H "Content-Type: application/json" \
  -d '{
    "cuenta_origen": "1000001234",
    "cuenta_destino": "2000009876",
    "monto": 250.50,
    "moneda": "USD"
  }'
```
*Resultado esperado: HTTP 200 OK. La transferencia se procesa y se dispara la evaluación de IA de forma asíncrona.*

**Escenario 2: Error por fondos insuficientes**
```bash
curl -X POST http://localhost:8080/api/v1/transferencias \
  -H "Content-Type: application/json" \
  -d '{
    "cuenta_origen": "1000001234",
    "cuenta_destino": "2000009876",
    "monto": 999999.00,
    "moneda": "USD"
  }'
```
*Resultado esperado: HTTP 400 Bad Request o HTTP 422 Unprocessable Entity. El sistema rechaza la transacción protegiendo la consistencia de la base de datos.*

**Escenario 3: Error por formato inválido (cuenta faltante)**
```bash
curl -X POST http://localhost:8080/api/v1/transferencias \
  -H "Content-Type: application/json" \
  -d '{
    "cuenta_destino": "2000009876",
    "monto": 100.00,
    "moneda": "USD"
  }'
```
*Resultado esperado: HTTP 400 Bad Request. El middleware de validación rechaza la petición antes de interactuar con la base de datos.*

## Evidencias

* **Video Demostrativo:** [[https://youtu.be/qIKL8xEzvUs](https://youtu.be/qIKL8xEzvUs)]

```
