# ============================================================
# OpenERP - Development Commands
# ============================================================

# Usar $args para capturar os argumentos
$Command = $args[0]

if ($Command -eq $null -or $Command -eq "") {
    Write-Host ""
    Write-Host "Comandos disponÃ­veis:" -ForegroundColor Cyan
    Write-Host "  build    - Compilar o projeto" -ForegroundColor White
    Write-Host "  run      - Executar o projeto" -ForegroundColor White
    Write-Host "  test     - Rodar testes" -ForegroundColor White
    Write-Host "  coverage - Rodar testes com cobertura" -ForegroundColor White
    Write-Host "  tidy     - Limpar dependÃªncias" -ForegroundColor White
    Write-Host "  fmt      - Formatar cÃ³digo" -ForegroundColor White
    Write-Host "  lint     - Rodar linter" -ForegroundColor White
    Write-Host "  gen      - Gerar models" -ForegroundColor White
    Write-Host "  docs     - Gerar documentaÃ§Ã£o" -ForegroundColor White
    Write-Host "  env      - Mostrar variÃ¡veis de ambiente" -ForegroundColor White
    Write-Host "  help     - Mostrar esta ajuda" -ForegroundColor White
    Write-Host ""
    Write-Host "Exemplo: .\scripts\dev.ps1 run" -ForegroundColor Cyan
    exit
}

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "OpenERP - Development Commands" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Comando: $Command" -ForegroundColor Yellow
Write-Host ""

switch ($Command.Trim()) {
    "build" {
        Write-Host "ðŸ”¨ Building..." -ForegroundColor Yellow
        go build -o bin/api.exe cmd/api/main.go
        if ($LASTEXITCODE -eq 0) {
            Write-Host "âœ… Build concluÃ­do!" -ForegroundColor Green
        } else {
            Write-Host "âŒ Build falhou!" -ForegroundColor Red
        }
    }

    "run" {
        Write-Host "ðŸš€ Running..." -ForegroundColor Yellow
        go run cmd/api/main.go
    }

    "test" {
        Write-Host "ðŸ§ª Testing..." -ForegroundColor Yellow
        go test -v ./...
    }

    "coverage" {
        Write-Host "ðŸ“Š Coverage..." -ForegroundColor Yellow
        go test -coverprofile=coverage.out ./...
        go tool cover -html=coverage.out
    }

    "tidy" {
        Write-Host "ðŸ§¹ Tidying..." -ForegroundColor Yellow
        go mod tidy
        go mod verify
        Write-Host "âœ… DependÃªncias limpas!" -ForegroundColor Green
    }

    "fmt" {
        Write-Host "ðŸŽ¨ Formatting..." -ForegroundColor Yellow
        go fmt ./...
        Write-Host "âœ… CÃ³digo formatado!" -ForegroundColor Green
    }

    "lint" {
        Write-Host "ðŸ” Linting..." -ForegroundColor Yellow
        if (Get-Command golangci-lint -ErrorAction SilentlyContinue) {
            golangci-lint run
        } else {
            Write-Host "âš ï¸  golangci-lint nÃ£o instalado. Instale com:" -ForegroundColor Yellow
            Write-Host "   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" -ForegroundColor White
        }
    }

    "gen" {
        Write-Host "ðŸ“ Generating models..." -ForegroundColor Yellow
        go run cmd/gen/main.go
    }

    "docs" {
        Write-Host "ðŸ“š Generating docs..." -ForegroundColor Yellow
        if (Get-Command swag -ErrorAction SilentlyContinue) {
            swag init --dir ./cmd/api/,./internal/handler/ -g main.go -o ./docs --parseDependency --parseInternal
        } else {
            Write-Host "âš ï¸  swag nÃ£o instalado. Instale com:" -ForegroundColor Yellow
            Write-Host "   go install github.com/swaggo/swag/cmd/swag@latest" -ForegroundColor White
        }
    }

    "env" {
        Write-Host "ðŸ“‹ Environment Variables:" -ForegroundColor Yellow
        Get-Content .env
    }

    "help" {
        Write-Host "Comandos disponÃ­veis:" -ForegroundColor Cyan
        Write-Host ""
        Write-Host "  build    - Compilar o projeto" -ForegroundColor White
        Write-Host "  run      - Executar o projeto" -ForegroundColor White
        Write-Host "  test     - Rodar testes" -ForegroundColor White
        Write-Host "  coverage - Rodar testes com cobertura" -ForegroundColor White
        Write-Host "  tidy     - Limpar dependÃªncias" -ForegroundColor White
        Write-Host "  fmt      - Formatar cÃ³digo" -ForegroundColor White
        Write-Host "  lint     - Rodar linter" -ForegroundColor White
        Write-Host "  gen      - Gerar models" -ForegroundColor White
        Write-Host "  docs     - Gerar documentaÃ§Ã£o" -ForegroundColor White
        Write-Host "  env      - Mostrar variÃ¡veis de ambiente" -ForegroundColor White
        Write-Host "  help     - Mostrar esta ajuda" -ForegroundColor White
        Write-Host ""
        Write-Host "Exemplo: .\scripts\dev.ps1 run" -ForegroundColor Cyan
    }

    default {
        Write-Host "âŒ Comando desconhecido: '$Command'" -ForegroundColor Red
        Write-Host ""
        Write-Host "Comandos disponÃ­veis:" -ForegroundColor Cyan
        Write-Host "  build, run, test, coverage, tidy, fmt, lint, gen, docs, env, help" -ForegroundColor White
        Write-Host ""
        Write-Host "Exemplo: .\scripts\dev.ps1 run" -ForegroundColor Cyan
    }
}