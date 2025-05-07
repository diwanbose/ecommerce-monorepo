.PHONY: install-prerequisites build test test-unit test-integration test-e2e run-local deploy-kind deploy-argocd clean lint lint-go lint-docs lint-api coverage coverage-go coverage-frontend

# Variables
GO_VERSION := 1.21.8
NODE_VERSION := 18
KIND_CLUSTER_NAME := ecommerce-cluster
KIND_CONFIG := k8s/kind-config.yaml

# Install all prerequisites
install-prerequisites:
	@echo "Installing prerequisites..."
	# Install Go
	@if ! command -v go &> /dev/null; then \
		echo "Installing Go..."; \
		wget https://go.dev/dl/go$(GO_VERSION).linux-amd64.tar.gz; \
		sudo rm -rf /usr/local/go; \
		sudo tar -C /usr/local -xzf go$(GO_VERSION).linux-amd64.tar.gz; \
		rm go$(GO_VERSION).linux-amd64.tar.gz; \
	fi
	# Install Node.js
	@if ! command -v node &> /dev/null; then \
		echo "Installing Node.js..."; \
		curl -fsSL https://deb.nodesource.com/setup_$(NODE_VERSION).x | sudo -E bash -; \
		sudo apt-get install -y nodejs; \
	fi
	# Install Docker
	@if ! command -v docker &> /dev/null; then \
		echo "Installing Docker..."; \
		sudo apt-get update; \
		sudo apt-get install -y docker.io; \
		sudo systemctl enable --now docker; \
		sudo usermod -aG docker $$USER; \
	fi
	# Install Docker Compose
	@if ! command -v docker-compose &> /dev/null; then \
		echo "Installing Docker Compose..."; \
		sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$$(uname -s)-$$(uname -m)" -o /usr/local/bin/docker-compose; \
		sudo chmod +x /usr/local/bin/docker-compose; \
	fi
	# Install kubectl
	@if ! command -v kubectl &> /dev/null; then \
		echo "Installing kubectl..."; \
		curl -LO "https://dl.k8s.io/release/$$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"; \
		sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl; \
		rm kubectl; \
	fi
	# Install Helm
	@if ! command -v helm &> /dev/null; then \
		echo "Installing Helm..."; \
		curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash; \
	fi
	# Install kind
	@if ! command -v kind &> /dev/null; then \
		echo "Installing kind..."; \
		curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-linux-amd64; \
		chmod +x ./kind; \
		sudo mv ./kind /usr/local/bin/kind; \
	fi
	# Install ArgoCD CLI
	@if ! command -v argocd &> /dev/null; then \
		echo "Installing ArgoCD CLI..."; \
		curl -sSL -o argocd-linux-amd64 https://github.com/argoproj/argo-cd/releases/latest/download/argocd-linux-amd64; \
		sudo install -m 555 argocd-linux-amd64 /usr/local/bin/argocd; \
		rm argocd-linux-amd64; \
	fi

# Build all components
build:
	@echo "Building all components..."
	# Build frontend
	cd frontend && npm install && npm run build
	# Build backend services
	cd backend/products && go build -o bin/products-service
	cd backend/cart && go build -o bin/cart-service
	cd backend/order && go build -o bin/order-service
	cd backend/feature-toggle && go build -o bin/feature-toggle-service

# Run all tests
test: test-unit test-integration test-e2e

# Run unit tests
test-unit:
	@echo "Running unit tests..."
	# Frontend unit tests
	cd frontend && npm run test:unit
	# Backend unit tests
	cd backend/products && go test ./... -v
	cd backend/cart && go test ./... -v
	cd backend/order && go test ./... -v
	cd backend/feature-toggle && go test ./... -v

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	cd backend/products && go test ./... -tags=integration
	cd backend/cart && go test ./... -tags=integration
	cd backend/order && go test ./... -tags=integration
	cd backend/feature-toggle && go test ./... -tags=integration

# Run E2E tests
test-e2e:
	@echo "Running E2E tests..."
	cd e2e && go test ./... -v

# Run locally using Docker Compose
run-local:
	@echo "Starting local development environment..."
	docker compose up --build

# Create kind cluster
deploy-kind:
	@echo "Creating kind cluster..."
	kind create cluster --name $(KIND_CLUSTER_NAME) --config $(KIND_CONFIG)
	# Install ArgoCD
	kubectl create namespace argocd
	kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml
	# Wait for ArgoCD to be ready
	kubectl wait --for=condition=available --timeout=300s deployment/argocd-server -n argocd

# Deploy using ArgoCD
deploy-argocd:
	@echo "Deploying using ArgoCD..."
	kubectl apply -f k8s/argocd/applications.yaml
	# Wait for applications to be synced
	kubectl wait --for=condition=available --timeout=300s deployment/frontend -n ecommerce
	kubectl wait --for=condition=available --timeout=300s deployment/products-service -n ecommerce
	kubectl wait --for=condition=available --timeout=300s deployment/cart-service -n ecommerce
	kubectl wait --for=condition=available --timeout=300s deployment/order-service -n ecommerce
	kubectl wait --for=condition=available --timeout=300s deployment/feature-toggle-service -n ecommerce

# Clean up
clean:
	@echo "Cleaning up..."
	kind delete cluster --name $(KIND_CLUSTER_NAME)
	docker compose down -v
	rm -rf frontend/node_modules
	rm -rf frontend/build
	rm -rf frontend/coverage
	rm -rf backend/*/bin
	rm -rf coverage
	rm -rf .cache
	find . -type d -name "dist" -exec rm -rf {} +
	find . -type f -name "*.out" -delete
	find . -type f -name "*.test" -delete
	find . -type f -name "*.prof" -delete
	find . -type f -name "*.trace" -delete
	find . -type f -name "*.log" -delete
	find . -type f -name ".DS_Store" -delete
	find . -type f -name "*.swp" -delete
	find . -type f -name "*.swo" -delete
	find . -type f -name "*~" -delete

# Linting
lint: lint-go lint-docs lint-api

# Go linting using golangci-lint
lint-go:
	@echo "Running Go linters..."
	cd backend/products && golangci-lint run --no-config ./...
	cd backend/cart && golangci-lint run --no-config ./...
	cd backend/order && golangci-lint run --no-config ./...
	cd backend/feature-toggle && golangci-lint run --no-config ./...

# Documentation linting using markdownlint
lint-docs:
	@echo "Running documentation linters..."
	markdownlint '**/*.md'

# API linting using spectral
lint-api:
	@echo "Running API linters..."
	spectral lint api/**/*.yaml api/**/*.json

# Test coverage
coverage: coverage-go coverage-frontend

# Go test coverage
coverage-go:
	@echo "Running Go test coverage..."
	@mkdir -p coverage
	cd backend/products && go test ./... -coverprofile=../../coverage/products.out -covermode=atomic
	cd backend/cart && go test ./... -coverprofile=../../coverage/cart.out -covermode=atomic
	cd backend/order && go test ./... -coverprofile=../../coverage/order.out -covermode=atomic
	cd backend/feature-toggle && go test ./... -coverprofile=../../coverage/feature-toggle.out -covermode=atomic
	cd coverage && go tool cover -func=products.out -o products.txt
	cd coverage && go tool cover -func=cart.out -o cart.txt
	cd coverage && go tool cover -func=order.out -o order.txt
	cd coverage && go tool cover -func=feature-toggle.out -o feature-toggle.txt
	cd coverage && go tool cover -html=products.out -o products.html
	cd coverage && go tool cover -html=cart.out -o cart.html
	cd coverage && go tool cover -html=order.out -o order.html
	cd coverage && go tool cover -html=feature-toggle.out -o feature-toggle.html

# Frontend test coverage
coverage-frontend:
	@echo "Running Frontend test coverage..."
	cd frontend && npm run test:coverage

# Default target
all: lint test build

# Install dependencies
install-deps:
	@echo "Installing dependencies..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2
	npm install -g markdownlint-cli
	npm install -g @stoplight/spectral-cli
