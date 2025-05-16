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
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

type CartItem struct {
	ProductID uint    `json:"productId"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type Cart struct {
	UserID uint       `json:"userId"`
	Items  []CartItem `json:"items"`
	Total  float64    `json:"total"`
}

type CartService struct {
	redisClient *redis.Client
	productsURL string
}

func NewCartService(redisClient *redis.Client, productsURL string) *CartService {
	return &CartService{
		redisClient: redisClient,
		productsURL: productsURL,
	}
}

func (s *CartService) GetCart(ctx context.Context, userID uint) (*Cart, error) {
	key := fmt.Sprintf("cart:%d", userID)
	data, err := s.redisClient.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return &Cart{UserID: userID, Items: []CartItem{}, Total: 0}, nil
	}
	if err != nil {
		return nil, err
	}

	var cart Cart
	if err := json.Unmarshal(data, &cart); err != nil {
		return nil, err
	}
	return &cart, nil
}

func (s *CartService) AddToCart(ctx context.Context, userID uint, item CartItem) error {
	// Get current cart
	cart, err := s.GetCart(ctx, userID)
	if err != nil {
		return err
	}

	// Check if product exists and get current price
	productURL := fmt.Sprintf("%s/api/products/%d", s.productsURL, item.ProductID)
	resp, err := http.Get(productURL)
	if err != nil {
		return fmt.Errorf("failed to fetch product: %v", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("Error closing response body: %v", cerr)
		}
	}()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("product not found")
	}

	var product struct {
		Price float64 `json:"price"`
		Stock int     `json:"stock"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return fmt.Errorf("failed to decode product: %v", err)
	}

	// Validate stock
	if product.Stock < item.Quantity {
		return fmt.Errorf("insufficient stock")
	}

	// Update cart
	found := false
	for i, existingItem := range cart.Items {
		if existingItem.ProductID == item.ProductID {
			cart.Items[i].Quantity += item.Quantity
			cart.Items[i].Price = product.Price
			found = true
			break
		}
	}

	if !found {
		item.Price = product.Price
		cart.Items = append(cart.Items, item)
	}

	// Calculate total
	cart.Total = 0
	for _, item := range cart.Items {
		cart.Total += item.Price * float64(item.Quantity)
	}

	// Save to Redis
	data, err := json.Marshal(cart)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("cart:%d", userID)
	return s.redisClient.Set(ctx, key, data, 24*time.Hour).Err()
}

func (s *CartService) UpdateCartItem(ctx context.Context, userID uint, productID uint, quantity int) error {
	cart, err := s.GetCart(ctx, userID)
	if err != nil {
		return err
	}

	// Check if product exists and get current price
	productURL := fmt.Sprintf("%s/api/products/%d", s.productsURL, productID)
	resp, err := http.Get(productURL)
	if err != nil {
		return fmt.Errorf("failed to fetch product: %v", err)
	}
	defer func() {
		if cerr := resp.Body.Close(); cerr != nil {
			log.Printf("Error closing response body: %v", cerr)
		}
	}()

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("product not found")
	}

	var product struct {
		Price float64 `json:"price"`
		Stock int     `json:"stock"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&product); err != nil {
		return fmt.Errorf("failed to decode product: %v", err)
	}

	// Validate stock
	if product.Stock < quantity {
		return fmt.Errorf("insufficient stock")
	}

	// Update cart
	found := false
	for i, item := range cart.Items {
		if item.ProductID == productID {
			if quantity <= 0 {
				cart.Items = append(cart.Items[:i], cart.Items[i+1:]...)
			} else {
				cart.Items[i].Quantity = quantity
				cart.Items[i].Price = product.Price
			}
			found = true
			break
		}
	}

	if !found && quantity > 0 {
		cart.Items = append(cart.Items, CartItem{
			ProductID: productID,
			Quantity:  quantity,
			Price:     product.Price,
		})
	}

	// Calculate total
	cart.Total = 0
	for _, item := range cart.Items {
		cart.Total += item.Price * float64(item.Quantity)
	}

	// Save to Redis
	data, err := json.Marshal(cart)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("cart:%d", userID)
	return s.redisClient.Set(ctx, key, data, 24*time.Hour).Err()
}

func (s *CartService) ClearCart(ctx context.Context, userID uint) error {
	key := fmt.Sprintf("cart:%d", userID)
	return s.redisClient.Del(ctx, key).Err()
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	// Initialize cart service
	service := NewCartService(redisClient, os.Getenv("PRODUCTS_SERVICE_URL"))

	// Initialize Gin router
	r := gin.Default()

	// API routes
	r.GET("/api/cart/:userId", func(c *gin.Context) {
		userID := uint(parseUint(c.Param("userId")))
		cart, err := service.GetCart(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart"})
			return
		}
		c.JSON(http.StatusOK, cart)
	})

	r.POST("/api/cart/:userId/items", func(c *gin.Context) {
		userID := uint(parseUint(c.Param("userId")))
		var item CartItem
		if err := c.BindJSON(&item); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		if err := service.AddToCart(c.Request.Context(), userID, item); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		cart, _ := service.GetCart(c.Request.Context(), userID)
		c.JSON(http.StatusOK, cart)
	})

	r.PUT("/api/cart/:userId/items/:productId", func(c *gin.Context) {
		userID := uint(parseUint(c.Param("userId")))
		productID := uint(parseUint(c.Param("productId")))
		var input struct {
			Quantity int `json:"quantity"`
		}
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		if err := service.UpdateCartItem(c.Request.Context(), userID, productID, input.Quantity); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		cart, _ := service.GetCart(c.Request.Context(), userID)
		c.JSON(http.StatusOK, cart)
	})

	r.DELETE("/api/cart/:userId", func(c *gin.Context) {
		userID := uint(parseUint(c.Param("userId")))
		if err := service.ClearCart(c.Request.Context(), userID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear cart"})
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
