package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string  `json:"name" gorm:"not null"`
	Description string  `json:"description"`
	Price       float64 `json:"price" gorm:"not null"`
	Image       string  `json:"image"`
	Stock       int     `json:"stock" gorm:"not null"`
}

type ProductService struct {
	db *gorm.DB
}

func NewProductService(db *gorm.DB) *ProductService {
	return &ProductService{db: db}
}

func (s *ProductService) GetProducts(ctx context.Context) ([]Product, error) {
	var products []Product
	if err := s.db.WithContext(ctx).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (s *ProductService) GetProduct(ctx context.Context, id uint) (*Product, error) {
	var product Product
	if err := s.db.WithContext(ctx).First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (s *ProductService) CreateProduct(ctx context.Context, product *Product) error {
	return s.db.WithContext(ctx).Create(product).Error
}

func (s *ProductService) UpdateProduct(ctx context.Context, id uint, product *Product) error {
	return s.db.WithContext(ctx).Model(&Product{}).Where("id = ?", id).Updates(product).Error
}

func (s *ProductService) DeleteProduct(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&Product{}, id).Error
}

func (s *ProductService) UpdateStock(ctx context.Context, id uint, quantity int) error {
	return s.db.WithContext(ctx).Model(&Product{}).
		Where("id = ?", id).
		Update("stock", gorm.Expr("stock + ?", quantity)).
		Error
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize database
	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate the schema
	if err := db.AutoMigrate(&Product{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize product service
	service := NewProductService(db)

	// Initialize Gin router
	r := gin.Default()

	// API routes
	r.GET("/api/products", func(c *gin.Context) {
		products, err := service.GetProducts(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
			return
		}
		c.JSON(http.StatusOK, products)
	})

	r.GET("/api/products/:id", func(c *gin.Context) {
		id := c.Param("id")
		var product Product
		if err := db.First(&product, id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusOK, product)
	})

	r.POST("/api/products", func(c *gin.Context) {
		var product Product
		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		if err := service.CreateProduct(c.Request.Context(), &product); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
			return
		}
		c.JSON(http.StatusCreated, product)
	})

	r.PUT("/api/products/:id", func(c *gin.Context) {
		id := c.Param("id")
		var product Product
		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		if err := service.UpdateProduct(c.Request.Context(), uint(parseUint(id)), &product); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
			return
		}
		c.JSON(http.StatusOK, product)
	})

	r.DELETE("/api/products/:id", func(c *gin.Context) {
		id := c.Param("id")
		if err := service.DeleteProduct(c.Request.Context(), uint(parseUint(id))); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
			return
		}
		c.Status(http.StatusNoContent)
	})

	// Start server
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server: ", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
		// Let the program exit naturally after defer cancel()
	}
}

func parseUint(s string) uint64 {
	var result uint64
	_, err := fmt.Sscanf(s, "%d", &result)
	if err != nil {
		return 0
	}
	return result
}
