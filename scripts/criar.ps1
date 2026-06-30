# Navegar para a raiz do projeto
cd D:\Projetos\OpenERP\backend

# Criar o script com a sintaxe correta
$content = @'
# ============================================================
# OpenERP - Development Commands
# ============================================================

# Usar $args para capturar os argumentos
$Command = $args[0]

if ($Command -eq $null -or $Command -eq "") {
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
    Write-Host "  docs     - Gerar documentação" -ForegroundColor White
    Write-Host "  env      - Mostrar variáveis de ambiente" -ForegroundColor White
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
        Write-Host "🔨 Building..." -ForegroundColor Yellow
        go build -o bin/api.exe cmd/api/main.go
        if ($LASTEXITCODE -eq 0) {
            Write-Host "✅ Build concluído!" -ForegroundColor Green
        } else {
            Write-Host "❌ Build falhou!" -ForegroundColor Red
        }
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
        if (Get-Command golangci-lint -ErrorAction SilentlyContinue) {
            golangci-lint run
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
            swag init -g cmd/api/main.go
        } else {
            Write-Host "⚠️  swag não instalado. Instale com:" -ForegroundColor Yellow
            Write-Host "   go install github.com/swaggo/swag/cmd/swag@latest" -ForegroundColor White
        }
    }

    "env" {
        Write-Host "📋 Environment Variables:" -ForegroundColor Yellow
        Get-Content .env
    }

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

    default {
        Write-Host "❌ Comando desconhecido: '$Command'" -ForegroundColor Red
        Write-Host ""
        Write-Host "Comandos disponíveis:" -ForegroundColor Cyan
        Write-Host "  build, run, test, coverage, tidy, fmt, lint, gen, docs, env, help" -ForegroundColor White
        Write-Host ""
        Write-Host "Exemplo: .\scripts\dev.ps1 run" -ForegroundColor Cyan
    }
}
'@

# Salvar com UTF-8 sem BOM
$Utf8NoBomEncoding = New-Object System.Text.UTF8Encoding $False
[System.IO.File]::WriteAllText("scripts\dev.ps1", $content, $Utf8NoBomEncoding)

Write-Host "✅ Script recriado com sucesso!" -ForegroundColor Green