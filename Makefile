.PHONY: build run test migrate-up migrate-down docker-build docker-up clean

BINARY_NAME=devpulse

build:
	go build -o bin/$(BINARY_NAME) cmd/server/main.go

run:
	go run cmd/server/main.go

test:
	go test -v ./...

migrate-up:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/devpulse?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://postgres:postgres@localhost:5432/devpulse?sslmode=disable" down

docker-build:
	docker build -t devpulse:latest .

docker-up:
	docker compose up -d

clean:
	rm -rf bin/
	go clean
