package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"railwaystation/internal/handlers"
	"railwaystation/internal/service"
)

// SetupRoutes sets up all the routes for the railway station API
func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// Create service instance
	stationService := service.NewRailwayStationService(db)
	stationHandler := handlers.NewRailwayStationHandler(stationService)

	// Railway stations routes
	api := router.Group("/api")
	{
		// Stations routes
		stations := api.Group("/stations")
		{
			stations.GET("/", stationHandler.GetAllStations)
			stations.GET("/:id", stationHandler.GetStationByID)
			stations.POST("/", stationHandler.CreateStation)
			stations.PUT("/:id", stationHandler.UpdateStation)
			stations.DELETE("/:id", stationHandler.DeleteStation)
			stations.GET("/search", stationHandler.SearchStations)
		}

		// Region routes
		regions := api.Group("/regions")
		{
			regions.GET("/", stationHandler.GetAllRegions)
			regions.GET("/:id", stationHandler.GetRegionByID)
			regions.POST("/", stationHandler.CreateRegion)
			regions.PUT("/:id", stationHandler.UpdateRegion)
			regions.DELETE("/:id", stationHandler.DeleteRegion)
		}

		// Area routes
		areas := api.Group("/areas")
		{
			areas.GET("/", stationHandler.GetAllAreas)
			areas.GET("/:id", stationHandler.GetAreaByID)
			areas.POST("/", stationHandler.CreateArea)
			areas.PUT("/:id", stationHandler.UpdateArea)
			areas.DELETE("/:id", stationHandler.DeleteArea)
		}

		// District routes
		districts := api.Group("/districts")
		{
			districts.GET("/", stationHandler.GetAllDistricts)
			districts.GET("/:id", stationHandler.GetDistrictByID)
			districts.POST("/", stationHandler.CreateDistrict)
			districts.PUT("/:id", stationHandler.UpdateDistrict)
			districts.DELETE("/:id", stationHandler.DeleteDistrict)
		}

		// Line routes
		lines := api.Group("/lines")
		{
			lines.GET("/", stationHandler.GetAllLines)
			lines.GET("/:id", stationHandler.GetLineByID)
			lines.POST("/", stationHandler.CreateLine)
			lines.PUT("/:id", stationHandler.UpdateLine)
			lines.DELETE("/:id", stationHandler.DeleteLine)
		}
	}
}