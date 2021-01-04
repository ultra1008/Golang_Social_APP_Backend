.PHONY: build test docker

build:
	go build -v -o ./build/highload-social-network ./cmd/highload-social-network

run:
	./build/highload-social-network

docker up:
	docker-compose -f deployment/docker-compose.yml up

docker down:
	docker-compose -f deployment/docker-compose.yml down