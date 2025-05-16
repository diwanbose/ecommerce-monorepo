# E-commerce Monorepo

A modern, cloud-native e-commerce platform built with microservices architecture and best DevOps practices.

## Features

- Modern microservices architecture
- Cloud-native design
- Feature toggles for safe deployments
- Comprehensive test coverage
- CI/CD pipeline with GitHub Actions
- Kubernetes deployment with ArgoCD
- Real-time inventory management
- Secure payment processing
- Order tracking and management

## Getting Started

### Prerequisites

- Go 1.21.8+
- Node.js 18+
- Docker
- kubectl
- Helm
- kind
- ArgoCD CLI

### Quick Start

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

Feature flags are managed through the feature-toggle service. To enable/disable features:

1. Update the feature flag in the configuration:

   ```yaml
   features:
     new_checkout: true
     dark_mode: false
   ```

2. The changes will be picked up automatically by the services.

3. Monitor the feature flag status in the admin dashboard.

### Code Quality Checks

- **Linting**: Enforces code style and best practices
- **Static Analysis**: Checks for potential bugs and security issues
- **Code Coverage**: Ensures adequate test coverage
- **Dependency Scanning**: Checks for vulnerable dependencies

### Testing

- **Unit Tests**: Tests individual components
- **Integration Tests**: Tests service interactions
- **E2E Tests**: Tests complete user flows
- **Performance Tests**: Tests system under load
- **Security Tests**: Tests for vulnerabilities

### Build and Deployment

1. Creates a kind cluster for local development
2. Installs ArgoCD for GitOps
3. Deploys all microservices:
   - Frontend
   - Products Service
   - Cart Service
   - Order Service
   - Feature Toggle Service
4. Sets up monitoring and logging

### Automated Checks

- Dependency updates and security patches
- Code quality and test coverage
- Docker image scanning
- Kubernetes manifest validation

## Development

To start development:

```bash
# Clone the repository
git clone https://github.com/yourusername/ecommerce-monorepo.git
cd ecommerce-monorepo

# Install dependencies
make install-deps
```

## Architecture

```plaintext
├── frontend/          # React frontend
├── backend/
│   ├── products/     # Product catalog service
│   ├── cart/         # Shopping cart service
│   ├── order/        # Order management service
│   └── feature-toggle/ # Feature flag service
├── k8s/              # Kubernetes manifests
└── scripts/          # Development scripts
```

## Test Suite

The project includes:

- Unit tests for all components
- Integration tests for service interactions
- End-to-end tests for critical flows
- Performance tests for key endpoints

```bash
# Run all tests
make test

# Run specific test suites
make test-unit
make test-integration
make test-e2e
```

## Environment Variables

The following environment variables are required:

```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=ecommerce

# API
API_PORT=8080
API_ENV=development

# Feature Flags
FEATURE_TOGGLE_URL=http://localhost:8081
```

## Troubleshooting

### Common Issues

1. **Database Connection Issues**
   - Check if PostgreSQL is running: `docker ps | grep postgres`
   - Verify database credentials in environment variables
   - Check database logs: `docker logs ecommerce-db`

2. **API Service Issues**
   - Check service logs: `docker logs ecommerce-api`
   - Verify all required environment variables are set
   - Check service health: `curl http://localhost:8080/health`

3. **Frontend Issues**
   - Clear browser cache
   - Check browser console for errors
   - Verify API endpoints are accessible

### Getting Help

- Check the [issues page](https://github.com/yourusername/ecommerce-monorepo/issues)
- Join our [Discord community](https://discord.gg/ecommerce-monorepo)
- Contact the development team at [dev@example.com](mailto:dev@example.com)

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Security

Please report security issues to [security@example.com](mailto:security@example.com).
See our [Security Policy](SECURITY.md) for more details.

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for a list of changes in each version.

## License

MIT

## Support

- Documentation: [docs.example.com](https://docs.example.com)
- Community: [Discord](https://discord.gg/ecommerce-monorepo)
- Email: [support@example.com](mailto:support@example.com)

## Questions?

Feel free to open an issue or contact the development team at [dev@example.com](mailto:dev@example.com)
