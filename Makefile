.PHONY: image

image:
	docker build --pull --rm -f "message_processor\Dockerfile" -t "message_processor:latest" .

compose:
    docker compose -f "message_processor\docker-compose.yaml" up -d --build

