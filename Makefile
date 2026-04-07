.PHONY: build test lint run tidy

build:
	go build -o bin/server ./cmd/server

test:
	go test -race -coverprofile=coverage.out ./...

lint:
	go vet ./...

run:
	go run ./cmd/server

tidy:
	go mod tidy

coverage:
	go tool cover -html=coverage.out
