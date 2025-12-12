package main

import (
	"Meteostation/handlers"
	"Meteostation/models"
	"Meteostation/pkg/geoservice"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	// Load environment variables
	godotenv.Load()
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	return db
}

func initGeoService() *geoservice.GeoHashClient {
	gs, err := geoservice.NewGeoHashClient(os.Getenv("GEOHASH_PRECISION"))
	if err != nil {
		log.Fatal("Failed to create GeoHashClient:", err)
	}
	return gs
}

func main() {
	db := initDB()

	if err := models.AutoMigrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	gs := initGeoService()

	meteoHandler := &handlers.MeteostationHandler{
		DB:         db,
		GeoService: gs,
	}

	router := gin.Default()
	api := router.Group("/api")
	{
		meteostations := api.Group("/meteostations")
		{
			meteostations.POST("", meteoHandler.Create)
			meteostations.GET("", meteoHandler.Index)
			meteostations.GET("/:index", meteoHandler.Get)
			meteostations.PUT("/:index", meteoHandler.Update)
			meteostations.DELETE("/:index", meteoHandler.Delete)
		}
		api.GET("/geohash/:geohash", meteoHandler.SearchByGeoHash)
	}

	router.Run(":8081") // listen and serve on localhost:8081
}
