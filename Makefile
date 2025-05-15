.PHONY: build run migrate-up migrate-down test docker-up docker-down docker-rebuild

build:
	go build -o bin/main ./cmd/server

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
