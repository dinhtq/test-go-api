.PHONY: build run swag-init test lint

all: swag-init build

swag-init:
	go run github.com/swaggo/swag/cmd/swag init

build: swag-init
	go build -o todo-api main.go

run: swag-init
	go run main.go

test:
	go test -v ./...

lint:
	golangci-lint run
