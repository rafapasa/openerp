# scripts/dev.ps1

param(
    [string]$Command = ""
)

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "OpenERP - Development Commands" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

switch ($Command) {
    "build" {
        Write-Host "🔨 Building..." -ForegroundColor Yellow
        go build -o bin/api.exe cmd/api/main.go
        Write-Host "✅ Build concluído!" -ForegroundColor Green
    }
    "run" {
        Write-Host "🚀 Running..." -ForegroundColor Yellow
        go run cmd/api/main.go
    }
    "test" {
        Write-Host "🧪 Testing..." -ForegroundColor Yellow
        go test -v ./...
    }
    "coverage" {
        Write-Host "📊 Coverage..." -ForegroundColor Yellow
        go test -coverprofile=coverage.out ./...
        go tool cover -html=coverage.out
    }
    "tidy" {
        Write-Host "🧹 Tidying..." -ForegroundColor Yellow
        go mod tidy
        go mod verify
        Write-Host "✅ Dependências limpas!" -ForegroundColor Green
    }
    "fmt" {
        Write-Host "🎨 Formatting..." -ForegroundColor Yellow
        go fmt ./...
        Write-Host "✅ Código formatado!" -ForegroundColor Green
    }
    "lint" {
        Write-Host "🔍 Linting..." -ForegroundColor Yellow
        golangci-lint run
    }
    "gen" {
        Write-Host "📝 Generating models..." -ForegroundColor Yellow
        go run cmd/gen/main.go
    }
    "docs" {
        Write-Host "📚 Generating docs..." -ForegroundColor Yellow
        swag init -g cmd/api/main.go
    }
    "env" {
        Write-Host "📋 Environment Variables:" -ForegroundColor Yellow
        Get-Content .env
    }
    "help" {
        Write-Host "Comandos disponíveis:" -ForegroundColor Cyan
        Write-Host "  build    - Compilar o projeto" -ForegroundColor White
        Write-Host "  run      - Executar o projeto" -ForegroundColor White
        Write-Host "  test     - Rodar testes" -ForegroundColor White
        Write-Host "  coverage - Rodar testes com cobertura" -ForegroundColor White
        Write-Host "  tidy     - Limpar dependências" -ForegroundColor White
        Write-Host "  fmt      - Formatar código" -ForegroundColor White
        Write-Host "  lint     - Rodar linter" -ForegroundColor White
        Write-Host "  gen      - Gerar models" -ForegroundColor White
        Write-Host "  docs     - Gerar documentação" -ForegroundColor White
        Write-Host "  env      - Mostrar variáveis de ambiente" -ForegroundColor White
        Write-Host "  help     - Mostrar esta ajuda" -ForegroundColor White
    }
    default {
        Write-Host "❌ Comando desconhecido: $Command" -ForegroundColor Red
        Write-Host ""
        Write-Host "Comandos disponíveis:" -ForegroundColor Cyan
        Write-Host "  build, run, test, coverage, tidy, fmt, lint, gen, docs, env, help" -ForegroundColor White
    }
}