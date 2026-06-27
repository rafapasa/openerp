#!/bin/bash
# Script para desenvolvimento

echo "========================================"
echo "OpenERP - Development Commands"
echo "========================================"
echo ""

case "$1" in
    "build")
        echo "🔨 Building..."
        go build -o bin/api cmd/api/main.go
        ;;
    "run")
        echo "🚀 Running..."
        go run cmd/api/main.go
        ;;
    "test")
        echo "🧪 Testing..."
        go test -v ./...
        ;;
    "coverage")
        echo "📊 Coverage..."
        go test -coverprofile=coverage.out ./...
        go tool cover -html=coverage.out
        ;;
    "tidy")
        echo "🧹 Tidying..."
        go mod tidy
        go mod verify
        ;;
    "fmt")
        echo "🎨 Formatting..."
        go fmt ./...
        ;;
    "lint")
        echo "🔍 Linting..."
        golangci-lint run
        ;;
    "gen")
        echo "📝 Generating models..."
        go run cmd/gen/main.go
        ;;
    "docs")
        echo "📚 Generating docs..."
        swag init -g cmd/api/main.go
        ;;
    *)
        echo "Comandos disponíveis:"
        echo "  build    - Compilar o projeto"
        echo "  run      - Executar o projeto"
        echo "  test     - Rodar testes"
        echo "  coverage - Rodar testes com cobertura"
        echo "  tidy     - Limpar dependências"
        echo "  fmt      - Formatar código"
        echo "  lint     - Rodar linter"
        echo "  gen      - Gerar models"
        echo "  docs     - Gerar documentação"
        ;;
esac