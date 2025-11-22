APP_NAME=avito-shop-service
MAIN_FILE_PATH=./cmd/avito-shop-service
GO_CMD=go
DOCKER_COMPOSE_CMD=docker compose

.PHONY: all build run up down logs

all: build

# -----------------------------
# Local commands
# -----------------------------

build:
	$(GO_CMD) build -o bin/$(APP_NAME).exe $(MAIN_FILE_PATH)

run:
	$(GO_CMD) run $(MAIN_FILE_PATH)

# -----------------------------
# Docker
# -----------------------------

up:
	$(DOCKER_COMPOSE_CMD) up --build -d

down:
	$(DOCKER_COMPOSE_CMD) down

logs:
	$(DOCKER_COMPOSE_CMD) logs -f app
