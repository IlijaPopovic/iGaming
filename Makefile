.PHONY: build run migrate-up migrate-down test docker-up docker-down docker-rebuild swagger-regenerate

build:
	go build -o bin/main ./cmd

run: build
	./bin/main

migrate-up:
	docker-compose up migrate

test:
	go test ./...

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-rebuild:
	docker-compose down && docker-compose up -d --build

swagger-regenerate:
	rm -rf docs/ && swag init -g cmd/main.go --output docs --dir . --exclude internal/migrations,vendor,testdata