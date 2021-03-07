.PHONY: build test

build:
	go build -v -o /build/hsn ./cmd/highload-social-network

run:
	./build/highload-social-network

test:
	go test ./... -v -race -timeout=30s