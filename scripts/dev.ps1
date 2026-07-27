<#
.SYNOPSIS
    Script de desenvolvimento para o projeto OpenERP.
.DESCRIPTION
    Centraliza os comandos comuns de desenvolvimento como build, run, test, lint, etc.
.PARAMETER Command
    O comando a ser executado.
.EXAMPLE
    .\scripts\dev.ps1 run
    Executa o projeto.
.EXAMPLE
    .\scripts\dev.ps1 test -v
    Executa os testes em modo verbose.
.EXAMPLE
    .\scripts\dev.ps1
    Mostra a ajuda.
#>
param(
    [Parameter(Mandatory=$false, Position=0)]
    [string]$Command,

    [Parameter(ValueFromRemainingArguments=$true)]
    [string[]]$RemainingArgs
)

# Define a codificação de saída para UTF-8 para evitar problemas com caracteres especiais
$OutputEncoding = [System.Text.Encoding]::UTF8

function Show-Help {
    Write-Host ""
    Write-Host "Comandos disponíveis:" -ForegroundColor Cyan
    Write-Host "  build    - Compilar o projeto" -ForegroundColor White
    Write-Host "  run      - Executar o projeto" -ForegroundColor White
    Write-Host "  test     - Rodar testes" -ForegroundColor White
    Write-Host "  coverage - Rodar testes com cobertura" -ForegroundColor White
    Write-Host "  tidy     - Limpar dependências" -ForegroundColor White
    Write-Host "  fmt      - Formatar código" -ForegroundColor White
    Write-Host "  lint     - Rodar linter" -ForegroundColor White
    Write-Host "  gen      - Gerar models" -ForegroundColor White
    Write-Host "  docs     - Gerar documentação Swagger" -ForegroundColor White
    Write-Host "  air      - Iniciar hot reload com Air (inclui geração de docs)" -ForegroundColor White
    Write-Host "  db:migrate - Rodar migrações do banco de dados" -ForegroundColor White
    Write-Host "  env      - Mostrar variáveis de ambiente" -ForegroundColor White
    Write-Host "  help     - Mostrar esta ajuda" -ForegroundColor White
    Write-Host ""
    Write-Host "Exemplo: .\scripts\dev.ps1 run" -ForegroundColor Cyan
    Write-Host "         .\scripts\dev.ps1 test -v" -ForegroundColor Cyan
}

if ([string]::IsNullOrWhiteSpace($Command) -or $Command -eq "help") {
    Show-Help
    exit
}

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "OpenERP - Development Commands" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Comando: $Command" -ForegroundColor Yellow -NoNewline
if ($RemainingArgs.Count -gt 0) {
    Write-Host " $($RemainingArgs)" -ForegroundColor DarkGray
} else {
    Write-Host ""
}
Write-Host ""

switch ($Command.Trim()) {
    "build" {
        Write-Host "🔨 Building..." -ForegroundColor Yellow
        go build -o bin/api.exe cmd/api/main.go @RemainingArgs
        if ($LASTEXITCODE -eq 0) {
            Write-Host "✅ Build concluído!" -ForegroundColor Green
        } else {
            Write-Host "❌ Build falhou!" -ForegroundColor Red
        }
    }

    "run" {
        Write-Host "🚀 Running..." -ForegroundColor Yellow
        go run cmd/api/main.go @RemainingArgs
    }

    "test" {
        Write-Host "🧪 Testing..." -ForegroundColor Yellow
        go test ./... @RemainingArgs
    }

    "coverage" {
        Write-Host "📊 Coverage..." -ForegroundColor Yellow
        go test -coverprofile=coverage.out ./... @RemainingArgs
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
        if (Get-Command golangci-lint -ErrorAction SilentlyContinue) {
            golangci-lint run @RemainingArgs
        } else {
            Write-Host "⚠️  golangci-lint não instalado. Instale com:" -ForegroundColor Yellow
            Write-Host "   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" -ForegroundColor White
        }
    }

    "gen" {
        Write-Host "📝 Generating models..." -ForegroundColor Yellow
        go run cmd/gen/main.go
    }

    "docs" {
        Write-Host "📚 Generating docs..." -ForegroundColor Yellow
        if (Get-Command swag -ErrorAction SilentlyContinue) {
            swag init --dir ./cmd/api/,./internal/handler/ -g main.go -o ./docs --parseDependency --parseInternal
        } else {
            Write-Host "⚠️  swag não instalado. Instale com:" -ForegroundColor Yellow
            Write-Host "   go install github.com/swaggo/swag/cmd/swag@latest" -ForegroundColor White
        }
    }

    "env" {
        Write-Host "📋 Environment Variables:" -ForegroundColor Yellow
        Get-Content .env
    }

    "air" {
        Write-Host "🚀 Starting hot reload with Air..." -ForegroundColor Yellow
        # Rodar docs primeiro, similar ao target 'dev' do makefile
        Write-Host "📚 Generating docs (pre-air)..." -ForegroundColor DarkCyan
        if (Get-Command swag -ErrorAction SilentlyContinue) {
            swag init --dir ./cmd/api/,./internal/handler/ -g main.go -o ./docs --parseDependency --parseInternal
        } else {
            Write-Host "⚠️  swag não instalado. Instale com:" -ForegroundColor Yellow
            Write-Host "   go install github.com/swaggo/swag/cmd/swag@latest" -ForegroundColor White
            exit 1 # Sair se o swag não estiver instalado e os docs não puderem ser gerados
        }
        Write-Host "✅ Docs generated." -ForegroundColor Green

        if (Get-Command air -ErrorAction SilentlyContinue) {
            air @RemainingArgs
        } else {
            Write-Host "⚠️  Air não instalado. Instale com:" -ForegroundColor Yellow
            Write-Host "   go install github.com/cosmtrek/air@latest" -ForegroundColor White
            exit 1 # Sair se o air não estiver instalado
        }
    }

    "db:migrate" {
        Write-Host "🔄 Running database migrations..." -ForegroundColor Yellow
        go run cmd/migrate/main.go @RemainingArgs
    }

    default {
        Write-Host "❌ Comando desconhecido: '$Command'" -ForegroundColor Red
        Show-Help
    }
}