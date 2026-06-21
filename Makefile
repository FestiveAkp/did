.PHONY: dev build run test

dev:
	go run .

build:
	go build -o bin/did .

run: build
	./bin/did

test:
	go test ./...
