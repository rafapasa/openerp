-- Criar banco de dados se não existir
CREATE DATABASE IF NOT EXISTS openerp CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Criar usuário se não existir
CREATE USER IF NOT EXISTS 'openerp'@'%' IDENTIFIED BY 'openerp123';

-- Dar permissões
GRANT ALL PRIVILEGES ON openerp.* TO 'openerp'@'%';

-- Aplicar mudanças
FLUSH PRIVILEGES;

-- Criar tabelas (exemplo)
USE openerp;

-- Tabela de usuários (exemplo)
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;