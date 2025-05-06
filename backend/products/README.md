# Products Service

The Products Service is responsible for managing product catalog, inventory, and stock management in the e-commerce platform.

## Features

- Product CRUD operations
- Inventory management
- Stock tracking
- Product search and filtering
- Category management

## API Endpoints

### Products

- `GET /api/products` - List all products
- `GET /api/products/:id` - Get product details
- `POST /api/products` - Create new product
- `PUT /api/products/:id` - Update product
- `DELETE /api/products/:id` - Delete product

### Inventory

- `GET /api/products/:id/stock` - Get product stock
- `PUT /api/products/:id/stock` - Update product stock
- `GET /api/products/stock/low` - Get low stock products

## Environment Variables

```env
DB_HOST=localhost
DB_PORT=5432
DB_NAME=ecommerce
DB_USER=postgres
DB_PASSWORD=postgres
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
docker build -t ecommerce-products .

# Run container
docker run -p 8080:8080 ecommerce-products
```

## Dependencies

- Go 1.21+
- PostgreSQL
- Gin Web Framework
- GORM
- Docker (optional) 