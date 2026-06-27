-- create_openerp_user.sql
USE etools_openerp;

-- Remover usuário se existir
DROP USER IF EXISTS 'openerp'@'localhost';

-- Criar usuário
CREATE USER 'openerp'@'localhost' IDENTIFIED BY '1234';

-- Conceder permissões
GRANT ALL PRIVILEGES ON etools_openerp.* TO 'openerp'@'localhost';

-- Aplicar
FLUSH PRIVILEGES;

-- Mostrar resultado
SELECT 'Usuário openerp criado com sucesso!' AS Status;
SHOW GRANTS FOR 'openerp'@'localhost';