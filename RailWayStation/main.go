// cmd/server/main.go
package main

import (
	"Railwaystation/internal/models"
	"Railwaystation/internal/routes"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
	_ "Railwaystation/docs"
)

// @title Railway Stations API
// @version 1.0
// @description REST API для управления железнодорожными станциями
// @host localhost:8080
// @BasePath /api
func main() {
	db := initDB()
	// Автомиграция моделей
	if err := autoMigrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Создание роутера
	router := gin.Default()

	// Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Middleware
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())

	// Настройка маршрутов
	routes.SetupRoutes(router, db)

	// Запуск сервера
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Server starting on :%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func initDB() *gorm.DB {
	// Load environment variables
	godotenv.Load()
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Successfully connected to database")

	return db
}

func autoMigrate(db *gorm.DB) error {
	modelsAr := []interface{}{
		&models.Region{}, &models.Area{}, &models.District{}, &models.Line{}, &models.RailStation{},
	}

	for _, model := range modelsAr {
		if err := db.AutoMigrate(model); err != nil {
			return err
		}
	}

	log.Println("Database migrated successfully")
	return nil
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
