.PHONY: image

image:
	docker build --pull --rm -f "Dockerfile" -t "message_processor:latest" .

compose-up:
    docker compose -f "docker-compose.yaml" up -d --build .

