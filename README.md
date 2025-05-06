# E-commerce Monorepo

This monorepo contains a complete e-commerce application with microservices architecture, featuring a React frontend and Go backend services.

## Architecture

- Frontend: React.js application
- Backend Services:
  - Products Service: Product catalog and inventory management
  - Cart Service: Shopping cart management
  - Order Service: Order processing and management
  - Feature Toggle Service: Feature flag management with GitOps integration

## Prerequisites

- Go 1.21 or later
- Node.js 18 or later
- Docker
- Docker Compose
- Make
- kubectl
- Helm 3
- kind
- argocd CLI

## Local Development Setup

1. Install prerequisites:
```bash
make install-prerequisites
```

2. Build all components:
```bash
make build
```

3. Run locally using Docker Compose:
```bash
make run-local
```

4. Run tests:
```bash
make test
```

## Feature Flags

The application uses a custom feature toggle implementation that is GitOps-driven. Feature flags are stored in Kubernetes ConfigMaps and can be managed through Git.

### Managing Feature Flags

1. Feature flags are defined in `k8s/helm/feature-toggle/values.yaml`
2. Changes to feature flags should be made through Git commits
3. ArgoCD will automatically sync the changes to the cluster

### Example: Hiding COD Payment Option

To hide the COD payment option:

1. Update the feature flag in `k8s/helm/feature-toggle/values.yaml`:
```yaml
featureFlags:
  enableCodPayment: false
```

2. Commit and push the changes
3. ArgoCD will automatically apply the changes

## CI/CD

The project uses GitHub Actions for CI/CD. The pipeline includes the following stages:

### Code Quality Checks
- **Linting**:
  - Go code linting using `golangci-lint`
  - Documentation linting using `markdownlint`
  - API specification linting using `spectral`
  - Frontend code linting using ESLint

### Testing
- **Unit Tests**:
  - Frontend unit tests using Jest
  - Backend unit tests for all Go services
- **Integration Tests**:
  - Backend service integration tests
- **End-to-End Tests**:
  - Complete application flow testing
- **Test Coverage**:
  - Go test coverage reports for all services
  - Frontend test coverage reports

### Build and Deployment
1. Creates a kind cluster for testing
2. Builds all components:
   - Frontend React application
   - Go backend services
3. Runs all test suites
4. Deploys to the kind cluster using ArgoCD:
   - Frontend application
   - All backend services
   - Feature toggle service
   - Required infrastructure components

### Automated Checks
- Dependency updates and security scanning
- Code formatting verification
- Documentation validation
- API specification validation

The CI pipeline can be run locally using the following commands:
```bash
make lint          # Run all linters
make test         # Run all tests
make coverage     # Generate coverage reports
make build        # Build all components
```

## Directory Structure

```
.
├── frontend/              # React frontend application
├── backend/              # Go backend services
│   ├── products/        # Products service
│   ├── cart/           # Cart service
│   ├── order/          # Order service
│   └── feature-toggle/ # Feature toggle service
├── k8s/                 # Kubernetes manifests
│   ├── helm/           # Helm charts
│   └── argocd/         # ArgoCD manifests
├── .github/            # GitHub Actions workflows
└── Makefile           # Build and deployment automation
```

## Testing

The project includes:
- Unit tests for all components
- Integration tests for backend services
- E2E tests for the complete application

Run tests using:
```bash
make test-unit      # Run unit tests
make test-integration  # Run integration tests
make test-e2e       # Run E2E tests
```

## License

MIT
