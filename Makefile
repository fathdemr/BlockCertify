# ─────────────────────────────────────────────────────────────────────────────
#  BlockCertify — Makefile
#  Usage: make <target>
# ─────────────────────────────────────────────────────────────────────────────

# Project settings
APP_NAME      := blockcertify
DOCKER_REPO   := fathdemr/blockcertify
GO_CMD        := ./cmd/server
BINARY        := ./bin/$(APP_NAME)-server
GO_FILES      := $(shell find . -name '*.go' -not -path './vendor/*')

# Docker image tags
BACKEND_TAG   := $(DOCKER_REPO)-backend
FRONTEND_TAG  := $(DOCKER_REPO)-frontend

# Colors for pretty output
CYAN  := \033[0;36m
GREEN := \033[0;32m
RESET := \033[0m

.PHONY: help \
        dev dev-backend dev-frontend \
        build build-backend build-frontend \
        test lint fmt tidy \
        docker-build docker-build-backend docker-build-frontend \
        docker-push docker-push-backend docker-push-frontend \
        up down logs logs-backend logs-frontend logs-db \
        db-shell db-reset \
        clean

# ─── Default: show help ───────────────────────────────────────────────────────
help:
	@echo ""
	@echo "$(CYAN)BlockCertify Makefile$(RESET)"
	@echo "─────────────────────────────────────────────────────────────"
	@echo "$(GREEN)Development$(RESET)"
	@echo "  make dev              Start all services locally (no Docker)"
	@echo "  make dev-backend      Run Go backend with hot-reload (air)"
	@echo "  make dev-frontend     Start Vite dev server"
	@echo ""
	@echo "$(GREEN)Build$(RESET)"
	@echo "  make build            Build Go backend binary"
	@echo "  make build-backend    Same as build"
	@echo "  make build-frontend   Build React production bundle"
	@echo ""
	@echo "$(GREEN)Code Quality$(RESET)"
	@echo "  make test             Run Go tests"
	@echo "  make lint             Run golangci-lint"
	@echo "  make fmt              Format Go source"
	@echo "  make tidy             go mod tidy"
	@echo ""
	@echo "$(GREEN)Docker$(RESET)"
	@echo "  make docker-build     Build all Docker images"
	@echo "  make docker-push      Push all images to registry (multi-arch)"
	@echo ""
	@echo "$(GREEN)Compose$(RESET)"
	@echo "  make up               Start full stack (db + backend + frontend)"
	@echo "  make down             Stop and remove containers"
	@echo "  make logs             Tail logs for all services"
	@echo "  make logs-backend     Tail backend logs"
	@echo "  make logs-frontend    Tail frontend logs"
	@echo "  make logs-db          Tail database logs"
	@echo ""
	@echo "$(GREEN)Database$(RESET)"
	@echo "  make db-shell         Open psql shell in running db container"
	@echo "  make db-reset         Drop and recreate the database volume"
	@echo ""
	@echo "$(GREEN)Misc$(RESET)"
	@echo "  make clean            Remove binary and node dist artifacts"
	@echo "─────────────────────────────────────────────────────────────"

# ─── Development (local, no Docker) ──────────────────────────────────────────
dev: dev-backend dev-frontend

dev-backend:
	@echo "$(CYAN)→ Starting Go backend...$(RESET)"
	@if command -v air > /dev/null 2>&1; then \
		air; \
	else \
		echo "  [air not found — running go run]"; \
		go run $(GO_CMD)/main.go; \
	fi

dev-frontend:
	@echo "$(CYAN)→ Starting Vite dev server...$(RESET)"
	cd frontend && npm run dev

# ─── Build ────────────────────────────────────────────────────────────────────
build: build-backend

build-backend:
	@echo "$(CYAN)→ Building Go backend binary...$(RESET)"
	@mkdir -p bin
	CGO_ENABLED=0 go build -ldflags="-w -s" -o $(BINARY) $(GO_CMD)
	@echo "$(GREEN)  Binary: $(BINARY)$(RESET)"

build-frontend:
	@echo "$(CYAN)→ Building React production bundle...$(RESET)"
	cd frontend && npm ci && npm run build
	@echo "$(GREEN)  Output: frontend/dist/$(RESET)"

# ─── Code Quality ─────────────────────────────────────────────────────────────
test:
	@echo "$(CYAN)→ Running Go tests...$(RESET)"
	go test ./... -v -race -timeout 120s

lint:
	@echo "$(CYAN)→ Running golangci-lint...$(RESET)"
	golangci-lint run ./...

fmt:
	@echo "$(CYAN)→ Formatting Go source...$(RESET)"
	gofmt -w $(GO_FILES)

tidy:
	@echo "$(CYAN)→ Running go mod tidy...$(RESET)"
	go mod tidy

# ─── Docker build (local, single-arch) ───────────────────────────────────────
docker-build: docker-build-backend docker-build-frontend

docker-build-backend:
	@echo "$(CYAN)→ Building backend Docker image...$(RESET)"
	docker build -f Dockerfile.backend -t $(BACKEND_TAG):latest .
	@echo "$(GREEN)  Image: $(BACKEND_TAG):latest$(RESET)"

docker-build-frontend:
	@echo "$(CYAN)→ Building frontend Docker image...$(RESET)"
	docker build -f Dockerfile.frontend -t $(FRONTEND_TAG):latest .
	@echo "$(GREEN)  Image: $(FRONTEND_TAG):latest$(RESET)"

# ─── Docker push (multi-arch, amd64 + arm64 for Fly/Render/M-chip) ───────────
docker-push: docker-push-backend docker-push-frontend

docker-push-backend:
	@echo "$(CYAN)→ Pushing backend image (linux/amd64,linux/arm64)...$(RESET)"
	docker buildx create --driver docker-container --name multiarch --use 2>/dev/null || \
		docker buildx use multiarch
	docker buildx build --platform linux/amd64,linux/arm64 \
		-f Dockerfile.backend \
		-t $(BACKEND_TAG):latest \
		--push .
	@echo "$(GREEN)  Pushed: $(BACKEND_TAG):latest$(RESET)"

docker-push-frontend:
	@echo "$(CYAN)→ Pushing frontend image (linux/amd64,linux/arm64)...$(RESET)"
	docker buildx create --driver docker-container --name multiarch --use 2>/dev/null || \
		docker buildx use multiarch
	docker buildx build --platform linux/amd64,linux/arm64 \
		-f Dockerfile.frontend \
		-t $(FRONTEND_TAG):latest \
		--push .
	@echo "$(GREEN)  Pushed: $(FRONTEND_TAG):latest$(RESET)"

# ─── Docker Compose ───────────────────────────────────────────────────────────
up:
	@echo "$(CYAN)→ Starting full stack (db + backend + frontend)...$(RESET)"
	docker compose up -d --build
	@echo "$(GREEN)  Frontend : http://localhost:$${FRONTEND_PORT:-80}$(RESET)"
	@echo "$(GREEN)  Backend  : http://localhost:$${PORT:-8080}$(RESET)"
	@echo "$(GREEN)  Postgres : localhost:$${DB_PORT:-5432}$(RESET)"

down:
	@echo "$(CYAN)→ Stopping containers...$(RESET)"
	docker compose down

logs:
	docker compose logs -f

logs-backend:
	docker compose logs -f backend

logs-frontend:
	docker compose logs -f frontend

logs-db:
	docker compose logs -f db

# ─── Database helpers ─────────────────────────────────────────────────────────
db-shell:
	@echo "$(CYAN)→ Opening psql shell (connected as postgres superuser)...$(RESET)"
	docker compose exec db psql -U postgres -d $${APP_DB_NAME:-blockcertify}

db-reset:
	@echo "$(CYAN)→ Resetting database volume (all data will be lost)...$(RESET)"
	docker compose down -v
	docker compose up -d db
	@echo "$(GREEN)  Database reset complete.$(RESET)"

# ─── Clean ────────────────────────────────────────────────────────────────────
clean:
	@echo "$(CYAN)→ Cleaning build artifacts...$(RESET)"
	rm -rf bin/ frontend/dist/
	@echo "$(GREEN)  Done.$(RESET)"