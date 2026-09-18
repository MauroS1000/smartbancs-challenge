/*
Archivo: init.sql
Propósito: Script de inicialización DDL/DML para el componente de persistencia (PostgreSQL).
Establece el esquema relacional central de SmartBancs, garantizando la integridad referencial
y aplicando restricciones a nivel de motor de base de datos (Defense in Depth) para soportar
escenarios de alta concurrencia transaccional.
*/

/*
Tabla: accounts
Almacena el estado actual de los fondos de los clientes.
Se implementa la restricción CHECK (balance >= 0) como un mecanismo de seguridad
estricto a nivel de base de datos. Esto garantiza que, incluso frente a un fallo
inesperado en la lógica de control de concurrencia de la aplicación, es matemáticamente
imposible que se registre un saldo negativo (sobregiro no autorizado).
*/
CREATE TABLE accounts (
    account_id VARCHAR(36) PRIMARY KEY,
    owner_name VARCHAR(100) NOT NULL,
    balance NUMERIC(15, 2) NOT NULL CHECK (balance >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

/*
Tabla: transactions
Actúa como un libro mayor (Ledger) de solo anexado que registra todos los movimientos financieros.
Representa el núcleo del patrón arquitectónico "Transactional Outbox": cada transferencia
exitosa se registra aquí en la misma transacción atómica que actualiza los saldos de 'accounts'.
Un proceso asíncrono leerá esta tabla posteriormente mediante micro-batching para propagar
los cambios al sistema central heredado (Bancs) sin afectar la latencia del API.
El campo correlation_id es indispensable para el rastreo y la observabilidad distribuida.
*/
CREATE TABLE transactions (
    transaction_id VARCHAR(36) PRIMARY KEY,
    source_account_id VARCHAR(36) NOT NULL REFERENCES accounts(account_id),
    destination_account_id VARCHAR(36) NOT NULL REFERENCES accounts(account_id),
    amount NUMERIC(15, 2) NOT NULL CHECK (amount > 0),
    status VARCHAR(20) NOT NULL, -- Estados: PENDING, COMPLETED, FAILED
    correlation_id VARCHAR(64) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Datos iniciales preconfigurados (Seed Data) para facilitar la ejecución,
-- validación y las pruebas de carga del Producto Mínimo Viable (MVP).
INSERT INTO accounts (account_id, owner_name, balance) VALUES
('acc-001', 'Juan Perez', 1500.00),
('acc-002', 'Maria Lopez', 350.00);