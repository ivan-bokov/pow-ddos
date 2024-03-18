.SILENT:
.DEFAULT_GOAL := help

help: ## Show this help
	@echo "Usage:\n  make <target>\n"
	@echo "Targets:"
	@grep -h -E '^[a-zA-Z_-].+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'


.PHONY: build
build: ## Build all binaries
	go build -o bin/word-of-wisdom-cli cmd/client/main.go
	go build -o bin/word-of-wisdom-server cmd/server/main.go


.PHONY: client
client: ## Run client
	go mod download
	go mod tidy
	go run cmd/client/main.go

.PHONY: server
server: ## Run server
	go mod download
	go mod tidy
	go run cmd/server/main.go

.PHONY: test
test: ## Run tests
	go clean --testcache
	go test ./...

.PHONY: start
start: ## Start server and client by docker-compose
	docker-compose build server client
	docker-compose up --abort-on-container-exit --force-recreate server client