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