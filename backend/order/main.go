package main

import (
	"context"
	"encoding/json"
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

type OrderItem struct {
	gorm.Model
	OrderID   uint    `json:"orderId" gorm:"not null"`
	ProductID uint    `json:"productId" gorm:"not null"`
	Quantity  int     `json:"quantity" gorm:"not null"`
	Price     float64 `json:"price" gorm:"not null"`
}

type Order struct {
	gorm.Model
	UserID        uint        `json:"userId" gorm:"not null"`
	Items         []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
	Total         float64     `json:"total" gorm:"not null"`
	Status        string      `json:"status" gorm:"not null;default:'pending'"`
	PaymentMethod string      `json:"paymentMethod" gorm:"not null"`
	Address       string      `json:"address" gorm:"not null"`
}

type OrderService struct {
	db            *gorm.DB
	cartURL       string
	productsURL   string
	featureURL    string
}

func NewOrderService(db *gorm.DB, cartURL, productsURL, featureURL string) *OrderService {
	return &OrderService{
		db:            db,
		cartURL:       cartURL,
		productsURL:   productsURL,
		featureURL:    featureURL,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID uint, paymentMethod, address string) (*Order, error) {
	// Get cart
	cartURL := fmt.Sprintf("%s/api/cart/%d", s.cartURL, userID)
	resp, err := http.Get(cartURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch cart: %v", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("Error closing response body: %v", cerr)
		}
	}()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("cart not found")
	}

	var cart struct {
		Items []struct {
			ProductID uint    `json:"productId"`
			Quantity  int     `json:"quantity"`
			Price     float64 `json:"price"`
		} `json:"items"`
		Total float64 `json:"total"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&cart); err != nil {
		return nil, fmt.Errorf("failed to decode cart: %v", err)
	}

	if len(cart.Items) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	// Check if COD is enabled if payment method is COD
	if paymentMethod == "cod" {
		codEnabled, err := s.isFeatureEnabled(ctx, "enableCodPayment")
		if err != nil {
			return nil, fmt.Errorf("failed to check COD feature: %v", err)
		}
		if !codEnabled {
			return nil, fmt.Errorf("COD payment is not enabled")
		}
	}

	// Validate stock and create order
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	order := &Order{
		UserID:        userID,
		Total:         cart.Total,
		Status:        "pending",
		PaymentMethod: paymentMethod,
		Address:       address,
	}

	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create order items and update stock
	for _, item := range cart.Items {
		// Check stock
		productURL := fmt.Sprintf("%s/api/products/%d", s.productsURL, item.ProductID)
		resp, err := http.Get(productURL)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to fetch product: %v", err)
		}
		defer func() {
			if cerr := resp.Body.Close(); cerr != nil {
				log.Printf("Error closing response body: %v", cerr)
			}
		}()

		var product struct {
			Stock int `json:"stock"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to decode product: %v", err)
		}

		if product.Stock < item.Quantity {
			tx.Rollback()
			return nil, fmt.Errorf("insufficient stock for product %d", item.ProductID)
		}

		// Create order item
		orderItem := OrderItem{
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     item.Price,
		}
		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			return nil, err
		}

		// Update stock
		updateURL := fmt.Sprintf("%s/api/products/%d/stock", s.productsURL, item.ProductID)
		req, err := http.NewRequest("PUT", updateURL, nil)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create stock update request: %v", err)
		}
		q := req.URL.Query()
		q.Add("quantity", fmt.Sprintf("-%d", item.Quantity))
		req.URL.RawQuery = q.Encode()

		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update stock: %v", err)
		}
		if err := resp.Body.Close(); err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}

	// Clear cart
	clearURL := fmt.Sprintf("%s/api/cart/%d", s.cartURL, userID)
	req, err := http.NewRequest("DELETE", clearURL, nil)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create cart clear request: %v", err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to clear cart: %v", err)
	}
	if err := resp.Body.Close(); err != nil {
		log.Printf("Error closing response body: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return order, nil
}

func (s *OrderService) GetOrder(ctx context.Context, orderID uint) (*Order, error) {
	var order Order
	if err := s.db.Preload("Items").First(&order, orderID).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (s *OrderService) GetUserOrders(ctx context.Context, userID uint) ([]Order, error) {
	var orders []Order
	if err := s.db.Preload("Items").Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, orderID uint, status string) error {
	return s.db.Model(&Order{}).Where("id = ?", orderID).Update("status", status).Error
}

func (s *OrderService) isFeatureEnabled(ctx context.Context, featureName string) (bool, error) {
	url := fmt.Sprintf("%s/api/flags/%s", s.featureURL, featureName)
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("Error closing response body: %v", cerr)
		}
	}()

	var result struct {
		Enabled bool `json:"enabled"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, err
	}
	return result.Enabled, nil
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
	if err := db.AutoMigrate(&Order{}, &OrderItem{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize order service
	service := NewOrderService(
		db,
		os.Getenv("CART_SERVICE_URL"),
		os.Getenv("PRODUCTS_SERVICE_URL"),
		os.Getenv("FEATURE_TOGGLE_URL"),
	)

	// Initialize Gin router
	r := gin.Default()

	// API routes
	r.POST("/api/orders", func(c *gin.Context) {
		var input struct {
			UserID        uint   `json:"userId" binding:"required"`
			PaymentMethod string `json:"paymentMethod" binding:"required"`
			Address       string `json:"address" binding:"required"`
		}
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		order, err := service.CreateOrder(c.Request.Context(), input.UserID, input.PaymentMethod, input.Address)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, order)
	})

	r.GET("/api/orders/:id", func(c *gin.Context) {
		orderID := uint(parseUint(c.Param("id")))
		order, err := service.GetOrder(c.Request.Context(), orderID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		c.JSON(http.StatusOK, order)
	})

	r.GET("/api/orders/user/:userId", func(c *gin.Context) {
		userID := uint(parseUint(c.Param("userId")))
		orders, err := service.GetUserOrders(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
			return
		}
		c.JSON(http.StatusOK, orders)
	})

	r.PUT("/api/orders/:id/status", func(c *gin.Context) {
		orderID := uint(parseUint(c.Param("id")))
		var input struct {
			Status string `json:"status" binding:"required"`
		}
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		if err := service.UpdateOrderStatus(c.Request.Context(), orderID, input.Status); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order status"})
			return
		}
		c.Status(http.StatusOK)
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