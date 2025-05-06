# Order Service

The Order Service handles order processing, payment integration, and order status management in the e-commerce platform.

## Features

- Order creation and management
- Payment processing
- Order status tracking
- Stock validation
- Cart integration
- Feature flag integration for payment methods

## API Endpoints

### Orders

- `POST /api/orders` - Create new order
- `GET /api/orders/:id` - Get order details
- `GET /api/orders/user/:userId` - Get user's orders
- `PUT /api/orders/:id/status` - Update order status

## Environment Variables

```env
DB_HOST=localhost
DB_PORT=5432
DB_NAME=ecommerce
DB_USER=postgres
DB_PASSWORD=postgres
CART_SERVICE_URL=http://cart:8080
PRODUCTS_SERVICE_URL=http://products:8080
FEATURE_TOGGLE_URL=http://feature-toggle:8080
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
docker build -t ecommerce-order .

# Run container
docker run -p 8080:8080 ecommerce-order
```

## Dependencies

- Go 1.21+
- PostgreSQL
- Gin Web Framework
- GORM
- Docker (optional)

## Integration Points

- Cart Service: Fetches cart items and clears cart after order creation
- Products Service: Validates stock and updates inventory
- Feature Toggle Service: Checks payment method availability 