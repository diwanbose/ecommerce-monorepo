package main

import (
	"context"
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

type FeatureFlag struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex"`
	Description string
	Enabled     bool
}

type FeatureToggleService struct {
	db *gorm.DB
}

func NewFeatureToggleService(db *gorm.DB) *FeatureToggleService {
	return &FeatureToggleService{db: db}
}

func (s *FeatureToggleService) IsEnabled(ctx context.Context, flagName string) (bool, error) {
	var flag FeatureFlag
	if err := s.db.WithContext(ctx).Where("name = ?", flagName).First(&flag).Error; err != nil {
		return false, err
	}
	return flag.Enabled, nil
}

func (s *FeatureToggleService) SetFlag(ctx context.Context, flagName string, enabled bool) error {
	return s.db.WithContext(ctx).Model(&FeatureFlag{}).Where("name = ?", flagName).Update("enabled", enabled).Error
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
	if err := db.AutoMigrate(&FeatureFlag{}); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Initialize feature toggle service
	service := NewFeatureToggleService(db)

	// Create default feature flags if they don't exist
	defaultFlags := []FeatureFlag{
		{Name: "enableCodPayment", Description: "Enable Cash on Delivery payment option", Enabled: false},
		{Name: "enableNewUI", Description: "Enable new UI design", Enabled: true},
	}

	for _, flag := range defaultFlags {
		var existingFlag FeatureFlag
		if err := db.Where("name = ?", flag.Name).First(&existingFlag).Error; err != nil {
			if err := db.Create(&flag).Error; err != nil {
				log.Printf("Failed to create default flag %s: %v", flag.Name, err)
			}
		}
	}

	// Initialize Gin router
	r := gin.Default()

	// API routes
	r.GET("/api/flags/:name", func(c *gin.Context) {
		flagName := c.Param("name")
		enabled, err := service.IsEnabled(c.Request.Context(), flagName)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Feature flag not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"enabled": enabled})
	})

	r.POST("/api/flags/:name", func(c *gin.Context) {
		flagName := c.Param("name")
		var input struct {
			Enabled bool `json:"enabled"`
		}
		if err := c.BindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}
		if err := service.SetFlag(c.Request.Context(), flagName, input.Enabled); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update feature flag"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Feature flag updated"})
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
	}
} 