.PHONY: build run migrate-up migrate-down test docker-up docker-down docker-rebuild swagger open start-project

build:
	go build -o bin/main ./cmd

run: build
	./bin/main

migrate-up:
	docker-compose run --rm migrate up

migrate-down:
	docker-compose run --rm migrate down

test:
	go test -v ./...

docker-up:
	docker-compose up -d app db

docker-down:
	docker-compose down

docker-rebuild:
	docker-compose down && docker-compose up -d --build

open:
	@echo "Opening Swagger UI..."
	@xdg-open http://localhost:8080/swagger/ 2>/dev/null || \
	 open http://localhost:8080/swagger/ 2>/dev/null || \
	 echo "Could not open browser - please visit http://localhost:8080/swagger/"

# before swagger-regenerate
# export PATH=$PATH:$(go env GOPATH)/bin
swagger-regenerate:
	rm -rf docs/ && swag init -g cmd/main.go --output docs --dir . --exclude internal/migrations,vendor,testdata

start-project: docker-rebuild open
	@echo "Project is running! Access:"
	@echo "- API: http://localhost:8080"
	@echo "- Swagger: http://localhost:8080/swagger/"