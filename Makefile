APP_NAME         := blockcertify
BACKEND_IMG      := fathdemr/$(APP_NAME)-backend
FRONTEND_IMG     := fathdemr/$(APP_NAME)-frontend
PLATFORMS        := linux/amd64,linux/arm64

# ─── Local Development ────────────────────────────────────────────────────────
.PHONY: run
run:                            ## Run backend locally (requires .env)
	go run ./cmd/server

.PHONY: run-frontend
run-frontend:                   ## Run frontend dev server
	cd frontend && npm run dev

# ─── Build ────────────────────────────────────────────────────────────────────
.PHONY: build
build:                          ## Build Go binary
	CGO_ENABLED=0 go build -o $(APP_NAME) ./cmd/server

.PHONY: clean
clean:                          ## Remove built binary
	rm -f $(APP_NAME)

# ─── Docker Build ─────────────────────────────────────────────────────────────
.PHONY: docker-build-backend
docker-build-backend:           ## Build backend Docker image (current platform)
	docker build -f backend.Dockerfile -t $(BACKEND_IMG):latest .

.PHONY: docker-build-frontend
docker-build-frontend:          ## Build frontend Docker image (current platform)
	docker build -f frontend.Dockerfile -t $(FRONTEND_IMG):latest .

.PHONY: docker-build
docker-build: docker-build-backend docker-build-frontend  ## Build all Docker images

# ─── Docker Push (multi-platform) ────────────────────────────────────────────
.PHONY: docker-push-backend
docker-push-backend:            ## Build & push backend to Docker Hub (multi-platform)
	@echo "==> Building & pushing $(BACKEND_IMG):latest for $(PLATFORMS)..."
	docker buildx build --platform $(PLATFORMS) \
		-f backend.Dockerfile \
		-t $(BACKEND_IMG):latest \
		--push .

.PHONY: docker-push-frontend
docker-push-frontend:           ## Build & push frontend to Docker Hub (multi-platform)
	@echo "==> Building & pushing $(FRONTEND_IMG):latest for $(PLATFORMS)..."
	docker buildx build --platform $(PLATFORMS) \
		-f frontend.Dockerfile \
		-t $(FRONTEND_IMG):latest \
		--push .

.PHONY: docker-push
docker-push: docker-push-backend docker-push-frontend  ## Push all images to Docker Hub

# ─── Docker Compose (local) ──────────────────────────────────────────────────
.PHONY: docker-up
docker-up:                      ## Start all services locally
	docker compose up -d --build

.PHONY: docker-down
docker-down:                    ## Stop all services
	docker compose down

.PHONY: docker-logs
docker-logs:                    ## Tail logs from all services
	docker compose logs -f

# ─── Server Deploy ────────────────────────────────────────────────────────────
# Usage: make deploy SERVER=user@185.252.234.84
.PHONY: deploy
deploy: docker-push             ## Push images & deploy to remote server via SSH
	@echo "==> Deploying to $(SERVER)..."
	ssh $(SERVER) "\
		cd ~/blockcertify && \
		docker compose pull && \
		docker compose up -d && \
		docker image prune -f"
	@echo "==> Deploy complete!"

# ─── Help ─────────────────────────────────────────────────────────────────────
.PHONY: help
help:                           ## Show available commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-25s\033[0m %s\n", $$1, $$2}'
