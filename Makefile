.PHONY: build
build:
	docker compose build

.PHONY: up
up:
	docker compose up -d