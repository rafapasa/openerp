# swagger.ps1
Write-Host "Gerando Swagger..." -ForegroundColor Yellow

# Aponta para os diretórios do main.go e dos handlers
swag init --dir ./cmd/api/,./internal/handler/ -g main.go -o ./docs --parseDependency --parseInternal

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Swagger gerado com sucesso!" -ForegroundColor Green
} else {
    Write-Host "❌ Erro ao gerar Swagger" -ForegroundColor Red
}