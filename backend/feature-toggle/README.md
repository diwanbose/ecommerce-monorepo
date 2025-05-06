# Feature Toggle Service

The Feature Toggle Service manages feature flags and feature availability across the e-commerce platform, enabling controlled feature rollouts and A/B testing.

## Features

- Feature flag management
- Dynamic feature configuration
- Feature targeting
- A/B testing support
- Feature analytics
- Environment-specific flags

## API Endpoints

### Feature Flags

- `GET /api/flags` - List all feature flags
- `GET /api/flags/:name` - Get feature flag status
- `POST /api/flags` - Create new feature flag
- `PUT /api/flags/:name` - Update feature flag
- `DELETE /api/flags/:name` - Delete feature flag

### Feature Analytics

- `GET /api/flags/:name/analytics` - Get feature usage analytics
- `GET /api/flags/analytics` - Get all feature analytics

## Environment Variables

```env
DB_HOST=localhost
DB_PORT=5432
DB_NAME=ecommerce
DB_USER=postgres
DB_PASSWORD=postgres
REDIS_HOST=localhost
REDIS_PORT=6379
```

## Development

```bash
# Run locally
go run main.go

# Run tests
go test ./...

# Run with coverage
go test ./... -coverprofile=coverage.out
```

## Docker

```bash
# Build image
docker build -t ecommerce-feature-toggle .

# Run container
docker run -p 8080:8080 ecommerce-feature-toggle
```

## Dependencies

- Go 1.21+
- PostgreSQL
- Redis
- Gin Web Framework
- GORM
- Docker (optional)

## Default Feature Flags

- `enableCodPayment`: Enable/disable Cash on Delivery payment
- `enableNewCheckout`: Enable/disable new checkout flow
- `enableProductReviews`: Enable/disable product reviews
- `enableWishlist`: Enable/disable wishlist functionality 