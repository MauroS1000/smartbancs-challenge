# Cargar variables de entorno si el archivo existe
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: up run down etl

# Levantar infraestructura
up:
	docker-compose up -d

# Ejecutar el servidor backend
run:
	cd backend && go run cmd/api/main.go

# Detener infraestructura
down:
	docker-compose down

# Ejecutar pipeline de datos
etl:
	python scripts/etl/transform.py