.PHONY: build run migrate-up migrate-down test docker-up docker-down docker-rebuild

build:
	go build -o bin/main ./cmd/server

run: build
	./bin/main

migrate-up:
	docker run --rm \
		-v $(PWD)/migrations:/migrations \
		--network="host" \
		mattes/goose:latest \
		mysql "root:password@tcp(localhost:3306)/igaming" up

migrate-down:
	docker run --rm \
		-v $(PWD)/migrations:/migrations \
		--network="host" \
		mattes/goose:latest \
		mysql "root:password@tcp(localhost:3306)/igaming" down

test:
	go test ./...

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-rebuild:
	docker-compose down && docker-compose up -d --build
