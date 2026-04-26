package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"meteodata2/config"
	"meteodata2/internal/database"
	"meteodata2/internal/repository"
	"meteodata2/internal/routes"
	"meteodata2/internal/usecase"
)

func init() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
}

func main() {
	// Initialize configuration
	cfg := config.NewConfig()

	// Initialize database connection
	db, err := database.NewInfluxDBClient(cfg.InfluxDB)
	if err != nil {
		log.Fatal("Failed to connect to InfluxDB:", err)
	}
	defer db.Close()

	// Initialize repository
	meteoRepo := repository.NewMeteoDataRepository(
		db.GetClient(),
		cfg.InfluxDB.Bucket,
		cfg.InfluxDB.Org,
	)

	// Initialize usecase
	meteoUsecase := usecase.NewMeteoDataUsecase(meteoRepo)

	// Setup Gin router
	r := gin.Default()

	// Register routes
	routes.RegisterRoutes(r, meteoUsecase)

	// Start server
	log.Printf("Starting server on port %s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
