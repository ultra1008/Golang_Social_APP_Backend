.PHONY: build test

build:
	go build -v ./cmd/highload-social-network

docker up:
	docker-compose -f deployment/docker-compose.yml up

docker down:
	docker-compose -f deployment/docker-compose.yml down