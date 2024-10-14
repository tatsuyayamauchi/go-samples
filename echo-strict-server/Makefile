CURRENT_DIR := $(shell pwd)

.PHONY: $(shell egrep -o '^(\._)?[a-z_-]+:' $(MAKEFILE_LIST) | sed 's/://')

up: ## Up all containers
	docker compose up -d
	open http://localhost:3030/

down: ## Shutdown all containers
	docker compose down

build: ## Build all containers
	docker-compose build --no-cache

setup: ## Setup tools
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest

codegen: ## Code generate
	@rm -rf $(CURRENT_DIR)/spec/*
	oapi-codegen --config oapi-codegen.config.yaml ./swagger.yaml

help: ## Show options
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
