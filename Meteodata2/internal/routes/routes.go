package routes

import (
	"time"

	"github.com/gin-gonic/gin"

	"meteodata2/internal/handlers"
	"meteodata2/internal/usecase"
)

// RegisterRoutes registers all application routes
func RegisterRoutes(r *gin.Engine, usecase usecase.MeteoDataUsecase) {
	// Initialize handlers
	regularHandler := handlers.NewRegularHandler(usecase)
	ledHandler := handlers.NewLEDHandler(usecase)
	tttrHandler := handlers.NewTTTRHandler(usecase)
	snowHandler := handlers.NewSNOWHandler(usecase)

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
			"message": "Server is running",
			"timestamp": time.Now().Unix(),
		})
	})

	// API routes
	api := r.Group("/api")
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
	}
}
