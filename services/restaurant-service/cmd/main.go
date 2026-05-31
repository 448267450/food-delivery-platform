package main

import (
	"fmt"
	"log"

	"github.com/448267450/food-delivery-platform/services/restaurant-service/internal/handler"
	"github.com/448267450/food-delivery-platform/services/restaurant-service/internal/model"
	"github.com/448267450/food-delivery-platform/services/restaurant-service/internal/repository"
	"github.com/448267450/food-delivery-platform/services/restaurant-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// Defaults
	viper.SetDefault("server.port", "8082")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.sslmode", "disable")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("failed to read config: %v", err)
		}
	}

	// Connect to database
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		viper.GetString("database.host"),
		viper.GetString("database.port"),
		viper.GetString("database.user"),
		viper.GetString("database.password"),
		viper.GetString("database.dbname"),
		viper.GetString("database.sslmode"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Auto migrate
	if err := db.AutoMigrate(&model.Restaurant{}, &model.MenuItem{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	// Wire up dependencies: Repository -> Service -> Handler
	restaurantRepo := repository.NewRestaurantRepository(db)
	menuItemRepo := repository.NewMenuItemRepository(db)
	restaurantSvc := service.NewRestaurantService(restaurantRepo, menuItemRepo)
	restaurantHandler := handler.NewRestaurantHandler(restaurantSvc)

	// Setup router
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "restaurant-service"})
	})

	restaurantHandler.RegisterRoutes(r)

	port := viper.GetString("server.port")
	log.Printf("restaurant-service running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
