# OpenERP - Backend 
 
## Tecnologias 
- Go 1.21 
- MySQL 8.0 
- Redis 
- JWT Authentication 
 
## Estrutura 
``` 
cmd/         # Pontos de entrada 
internal/    # Codigo interno 
pkg/         # Pacotes reutilizaveis 
migrations/  # Migracoes SQL 
api/         # Especificacoes da API 
``` 

Subir tudo (aplicação + monitoramento)
bash

docker-compose up -d

Ver logs
bash

docker-compose logs -f app

Ver todos os serviços
bash

docker-compose ps

Parar tudo
bash

docker-compose down

Parar e remover volumes (limpeza total)
bash

docker-compose down -v

Reconstruir a aplicação
bash

docker-compose up -d --build app

🌐 Acessos
Serviço	URL	Usuário/Senha
API	http://localhost:8080	-
Swagger	http://localhost:8080/swagger/index.html	-
Metrics	http://localhost:8080/metrics	-
Health	http://localhost:8080/health	-
Prometheus	http://localhost:9090	-
Jaeger UI	http://localhost:16686	-
Grafana	http://localhost:3000	admin/admin
MySQL	localhost:3306	root/root
Redis	localhost:6379	-
Node Exporter	http://localhost:9100

Acessar o MySQL
bash

docker exec -it openerp-mysql mysql -u root -proot

Acessar com o usuário openerp
bash

docker exec -it openerp-mysql mysql -u openerp -popenerp123 openerp

Verificar versão
bash

docker exec -it openerp-mysql mysql -u root -proot -e "SELECT VERSION();"

Backup do banco
bash

docker exec openerp-mysql mysqldump -u root -proot openerp > backup.sql

Restaurar backup
bash

cat backup.sql | docker exec -i openerp-mysql mysql -u root -proot openerp