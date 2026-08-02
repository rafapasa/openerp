# Makefile
.PHONY: swagger

swagger:
#	 swag init -g ./cmd/api/main.go -o ./docs --parseDependency --parseInternal
#	swag init --dir ./cmd/api/,./internal/handlers -g ./main.go -o ./docs --parseDependency --parseInternal
	swag init --dir ./cmd/api/,./internal/handler/ -g main.go -o ./docs --parseDependency --parseInternal

run:
	go run ./cmd/api/main.go

dev: swagger
	air

builder: 
	go build -o ./bin/openerp_api.exe .\cmd\api\main.go	


.PHONY: start-monitoring stop-monitoring

start-monitoring:
	docker-compose -f docker-compose.monitoring.yml up -d
	@echo "✅ Prometheus: http://localhost:9090"
	@echo "✅ Jaeger: http://localhost:16686"
	@echo "✅ Grafana: http://localhost:3000 (admin/admin)"

stop-monitoring:
	docker-compose -f docker-compose.monitoring.yml down -v

logs:
	docker-compose -f docker-compose.monitoring.yml logs -f