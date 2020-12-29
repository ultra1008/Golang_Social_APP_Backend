.PHONY: build test

docker up:
	docker-compose -f deployment/docker-compose.yml up

docker down:
	docker-compose -f deployment/docker-compose.yml down