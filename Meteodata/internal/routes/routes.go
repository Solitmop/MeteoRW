package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"meteodata/internal/handlers"
)

// SetupRoutes настраивает маршруты API
func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// Инициализируем обработчики
	regularHandler := handlers.NewRegularHandler(db)
	ledHandler := handlers.NewLEDHandler(db)
	tttrHandler := handlers.NewTTTRHandler(db)
	snowHandler := handlers.NewSNOWHandler(db)
	statsHandler := handlers.NewStatisticsHandler(db)

	api := router.Group("/api")
	{
		// Regular measurements
		regular := api.Group("/regular")
		{
			regular.POST("", regularHandler.Create)
			regular.GET("", regularHandler.Get)
			regular.GET("/:index/:date", regularHandler.GetByID)
			regular.PUT("/:index/:date", regularHandler.Update)
			regular.DELETE("/:index/:date", regularHandler.Delete)
		}

		// LED measurements
		led := api.Group("/led")
		{
			led.POST("", ledHandler.Create)
			led.GET("", ledHandler.Get)
			led.GET("/:index/:date", ledHandler.GetByID)
			led.PUT("/:index/:date", ledHandler.Update)
			led.DELETE("/:index/:date", ledHandler.Delete)
		}

		// TTTR measurements
		tttr := api.Group("/tttr")
		{
			tttr.POST("", tttrHandler.Create)
			tttr.GET("", tttrHandler.Get)
			tttr.GET("/:index/:date", tttrHandler.GetByID)
			tttr.DELETE("/:index/:date", tttrHandler.Delete)
			tttr.GET("/monthly-temperature", tttrHandler.GetMonthlyTemperature)
		}

		// SNOW measurements
		snow := api.Group("/snow")
		{
			snow.POST("", snowHandler.Create)
			snow.GET("", snowHandler.Get)
			snow.GET("/:index/:date", snowHandler.GetByID)
			snow.DELETE("/:index/:date", snowHandler.Delete)
		}

		// Statistics
		api.GET("/statistics", statsHandler.GetStatistics)

		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":    "OK",
				"timestamp": time.Now().Unix(),
			})
		})
	}
}
