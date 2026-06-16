# TaskFlow developer tasks. Run `make help` to see all commands.

.PHONY: help run build test cover fmt vet tidy docker-build docker-run clean

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-14s\033[0m %s\n", $$1, $$2}'

run: ## Run the server locally
	go run .

build: ## Build the binary into ./taskflow
	CGO_ENABLED=0 go build -ldflags="-s -w" -o taskflow .

test: ## Run all tests
	go test ./...

cover: ## Run tests with a coverage summary
	go test -cover ./...

fmt: ## Format all code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

tidy: ## Tidy go.mod/go.sum
	go mod tidy

docker-build: ## Build the Docker image
	docker build -t taskflow:latest .

docker-run: ## Run the Docker image on port 8080
	docker run --rm -p 8080:8080 -e TASKFLOW_JWT_SECRET=change-me taskflow:latest

clean: ## Remove build artifacts and local DB
	rm -f taskflow cover.out *.db *.db-journal
