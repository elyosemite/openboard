.DEFAULT_GOAL := help
.PHONY: help check backend-check backend-fmt backend-vet backend-lint backend-complexity backend-test \
        frontend-install frontend-check sonar up down

SONAR_HOST_URL ?= http://host.docker.internal:9000

help: ## List targets
	@grep -E '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "} {printf "  %-20s %s\n", $$1, $$2}'

check: backend-check frontend-check ## Run every quality gate

backend-check: backend-fmt backend-vet backend-lint backend-complexity backend-test ## All backend gates

backend-fmt: ## Fail if any Go file is not gofmt-formatted
	@out=$$(cd backend && gofmt -l .); if [ -n "$$out" ]; then echo "$$out"; exit 1; fi

backend-vet: ## go vet
	cd backend && go vet ./...

backend-lint: ## golangci-lint
	cd backend && golangci-lint run ./...

backend-complexity: ## gocyclo <= 10, gocognit <= 15
	cd backend && gocyclo -over 10 . && gocognit -over 15 .

backend-test: ## Tests with coverage profile (backend/coverage.out)
	cd backend && go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out | tail -1

frontend-install: ## Install frontend dependencies
	cd frontend && pnpm install --frozen-lockfile

frontend-check: ## Prettier, ESLint, tsc and tests with coverage
	cd frontend && pnpm format:check && pnpm lint && pnpm typecheck && pnpm test:coverage

sonar: backend-test ## Analyze on the local SonarQube (needs OPENBOARD_SONAR_TOKEN)
	cd frontend && pnpm test:coverage
	docker run --rm -e SONAR_HOST_URL=$(SONAR_HOST_URL) -e SONAR_TOKEN=$$OPENBOARD_SONAR_TOKEN \
		-v "$(CURDIR):/usr/src" sonarsource/sonar-scanner-cli

up: ## Start the local stack (postgres, liquibase, backend, frontend)
	docker compose up -d --build

down: ## Stop the local stack and drop its volumes
	docker compose down -v
