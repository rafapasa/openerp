# ============================================================
# OpenERP - Development Commands
# ============================================================

param(
    [string]$Command = ""
)

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "OpenERP - Development Commands" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

switch ($Command) {
    # ============================================================
    # BUILD - Compilar o projeto
    # ============================================================
    "build" {
        Write-Host "🔨 Building..." -ForegroundColor Yellow
        go build -o bin/api.exe cmd/api/main.go
        if ($LASTEXITCODE -eq 0) {
            Write-Host "✅ Build concluído!" -ForegroundColor Green
        } else {
            Write-Host "❌ Build falhou!" -ForegroundColor Red
        }
    }

    # ============================================================
    # RUN - Executar o projeto
    # ============================================================
    "run" {
        Write-Host "🚀 Running..." -ForegroundColor Yellow
        go run cmd/api/main.go
    }

    # ============================================================
    # TEST - Rodar testes
    # ============================================================
    "test" {
        Write-Host "🧪 Testing..." -ForegroundColor Yellow
        go test -v ./...
    }

    # ============================================================
    # COVERAGE - Testes com cobertura
    # ============================================================
    "coverage" {
        Write-Host "📊 Coverage..." -ForegroundColor Yellow
        go test -coverprofile=coverage.out ./...
        go tool cover -html=coverage.out
    }

    # ============================================================
    # TIDY - Limpar dependências
    # ============================================================
    "tidy" {
        Write-Host "🧹 Tidying..." -ForegroundColor Yellow
        go mod tidy
        go mod verify
        Write-Host "✅ Dependências limpas!" -ForegroundColor Green
    }

    # ============================================================
    # FMT - Formatar código
    # ============================================================
    "fmt" {
        Write-Host "🎨 Formatting..." -ForegroundColor Yellow
        go fmt ./...
        Write-Host "✅ Código formatado!" -ForegroundColor Green
    }

    # ============================================================
    # LINT - Rodar linter
    # ============================================================
    "lint" {
        Write-Host "🔍 Linting..." -ForegroundColor Yellow
        if (Get-Command golangci-lint -ErrorAction SilentlyContinue) {
            golangci-lint run
        } else {
            Write-Host "⚠️  golangci-lint não instalado. Instale com:" -ForegroundColor Yellow
            Write-Host "   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" -ForegroundColor White
        }
    }

    # ============================================================
    # GEN - Gerar models
    # ============================================================
    "gen" {
        Write-Host "📝 Generating models..." -ForegroundColor Yellow
        go run cmd/gen/main.go
    }

    # ============================================================
    # DOCS - Gerar documentação Swagger
    # ============================================================
    "docs" {
        Write-Host "📚 Generating docs..." -ForegroundColor Yellow
        if (Get-Command swag -ErrorAction SilentlyContinue) {
            swag init -g cmd/api/main.go
        } else {
            Write-Host "⚠️  swag não instalado. Instale com:" -ForegroundColor Yellow
            Write-Host "   go install github.com/swaggo/swag/cmd/swag@latest" -ForegroundColor White
        }
    }

    # ============================================================
    # ENV - Mostrar variáveis de ambiente
    # ============================================================
    "env" {
        Write-Host "📋 Environment Variables:" -ForegroundColor Yellow
        Get-Content .env
    }

    # ============================================================
    # HELP - Mostrar ajuda
    # ============================================================
    "help" {
        Write-Host "Comandos disponíveis:" -ForegroundColor Cyan
        Write-Host ""
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
        Write-Host ""
        Write-Host "Exemplo: .\scripts\dev.ps1 run" -ForegroundColor Cyan
    }

    # ============================================================
    # COMANDO DESCONHECIDO
    # ============================================================
    default {
        Write-Host "❌ Comando desconhecido: $Command" -ForegroundColor Red
        Write-Host ""
        Write-Host "Comandos disponíveis:" -ForegroundColor Cyan
        Write-Host "  build, run, test, coverage, tidy, fmt, lint, gen, docs, env, help" -ForegroundColor White
        Write-Host ""
        Write-Host "Exemplo: .\scripts\dev.ps1 run" -ForegroundColor Cyan
    }
}