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

## Evidencias

* **Video Demostrativo:** [https://youtu.be/qIKL8xEzvUs]