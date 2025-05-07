# Cart Service

The Cart Service manages shopping cart functionality, including adding/removing items, updating quantities, and cart persistence.

## Features

- Cart CRUD operations
- Item management
- Price calculations
- Cart persistence
- Cart expiration
- Cart merging (guest to logged-in user)

## API Endpoints

### Cart

- `GET /api/cart/:userId` - Get user's cart
- `POST /api/cart/:userId` - Create new cart
- `PUT /api/cart/:userId` - Update cart
- `DELETE /api/cart/:userId` - Delete cart

### Cart Items

- `POST /api/cart/:userId/items` - Add item to cart
- `PUT /api/cart/:userId/items/:itemId` - Update cart item
- `DELETE /api/cart/:userId/items/:itemId` - Remove item from cart

## Environment Variables

```env
DB_HOST=localhost
DB_PORT=5432
DB_NAME=ecommerce
DB_USER=postgres
DB_PASSWORD=postgres
REDIS_HOST=localhost
REDIS_PORT=6379
PRODUCTS_SERVICE_URL=http://products:8080
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
docker build -t ecommerce-cart .

# Run container
docker run -p 8080:8080 ecommerce-cart
```

## Dependencies

- Go 1.21+
- PostgreSQL
- Redis
- Gin Web Framework
- GORM
- Docker (optional)
