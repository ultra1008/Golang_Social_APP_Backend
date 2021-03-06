.PHONY: build test

build:
	go build -v ./cmd/highload-social-network

run:
	./build/highload-social-network

test:
	go test ./... -v -race -timeout=30s