
---

### 2. Criar Arquivo de Comandos: `docs/GO_COMMANDS.md`

```markdown
# Guia Rápido - Comandos Go

## Gerenciamento de Dependências

```bash
# Inicializar módulo
go mod init github.com/openerp/backend

# Adicionar dependência
go get github.com/gin-gonic/gin

# Adicionar dependência com versão específica
go get gorm.io/gorm@v1.25.5

# Remover dependências não utilizadas
go mod tidy

# Verificar dependências
go mod verify

# Listar todas as dependências
go list -m all

# Visualizar grafo de dependências
go mod graph

# Executar arquivo específico
go run cmd/api/main.go

# Executar com variáveis de ambiente
DB_HOST=localhost go run cmd/api/main.go

# Build do executável
go build -o bin/api cmd/api/main.go

# Build para Windows
 

# Build para Linux
GOOS=linux GOARCH=amd64 go build -o bin/api cmd/api/main.go

# Build com otimização (produção)
go build -ldflags="-s -w" -o bin/api cmd/api/main.go

# Instalar Air (hot reload)
go install github.com/air-verse/air@latest

# Rodar com Air
air

# Instalar Swagger (documentação API)
go install github.com/swaggo/swag/cmd/swag@latest

# Gerar documentação Swagger
swag init -g cmd/api/main.go

# Formatar código
go fmt ./...

# Verificar código (lint)
go vet ./...

# Instalar golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Rodar lint
golangci-lint run

Comando	Descrição
go mod tidy	Limpar dependências
go run .	Executar o projeto atual
go test ./...	Testar tudo
go build	Compilar
go fmt ./...	Formatar tudo
go vet ./...	Verificar código