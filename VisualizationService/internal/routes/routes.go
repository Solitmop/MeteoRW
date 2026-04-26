package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"visualizationservice/internal/handlers"
)

// SetupRoutes sets up all the routes for the visualization API
func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// Create handler instance
	visualizationHandler := handlers.NewVisualizationHandler()

	// Visualization routes
	api := router.Group("/api")
	{
		// Map routes
		mapRoutes := api.Group("/map")
		{
			mapRoutes.GET("/stations", visualizationHandler.GetStationMap)
			mapRoutes.GET("/regions", visualizationHandler.GetRegionMap)
			mapRoutes.GET("/lines", visualizationHandler.GetLineMap)
		}

		// Chart routes
		chartRoutes := api.Group("/charts")
		{
			chartRoutes.GET("/temperature", visualizationHandler.GetTemperatureChart)
			chartRoutes.GET("/humidity", visualizationHandler.GetHumidityChart)
			chartRoutes.GET("/pressure", visualizationHandler.GetPressureChart)
		}

		// Dashboard routes
		dashboardRoutes := api.Group("/dashboard")
		{
			dashboardRoutes.GET("/", visualizationHandler.GetDashboard)
		}
	}
}