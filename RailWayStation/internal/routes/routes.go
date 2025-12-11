// internal/routes/routes.go
package routes

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "Railwaystation/internal/handlers"
    "Railwaystation/internal/service"
	"Railwaystation/internal/models"
    "github.com/aoliveti/geohash"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
    // Инициализация сервисов
    geoSvc := service.NewGeoHashService(geohash.City)
    
    // Инициализация обработчиков
    regionHandler := handlers.NewRegionHandler(db)
    areaHandler := handlers.NewAreaHandler(db)
    districtHandler := handlers.NewDistrictHandler(db)
    lineHandler := handlers.NewLineHandler(db)
    railStationHandler := handlers.NewRailStationHandler(db, geoSvc)
    
    // Группа API
    api := router.Group("/api")
    {
        // Регионы
        regions := api.Group("/regions")
        {
            regions.POST("", regionHandler.CreateRegion)
            regions.GET("", regionHandler.GetRegions)
            regions.GET(":id", regionHandler.GetRegion)
            regions.PUT(":id", regionHandler.UpdateRegion)
            regions.DELETE(":id", regionHandler.DeleteRegion)
        }
        
        // Области
        areas := api.Group("/areas")
        {
            areas.POST("", areaHandler.CreateArea)
            areas.GET("", areaHandler.GetAreas)
            areas.GET(":id", areaHandler.GetArea)
            areas.PUT(":id", areaHandler.UpdateArea)
            areas.DELETE(":id", areaHandler.DeleteArea)
            areas.GET("/region/:region_id", areaHandler.GetAreasByRegion)
        }
        
        // Районы
        districts := api.Group("/districts")
        {
            districts.POST("", districtHandler.CreateDistrict)
            districts.GET("", districtHandler.GetDistricts)
            districts.GET(":id", districtHandler.GetDistrict)
            districts.PUT(":id", districtHandler.UpdateDistrict)
            districts.DELETE(":id", districtHandler.DeleteDistrict)
            districts.GET("/area/:area_id", districtHandler.GetDistrictsByArea)
        }
        
        // Линии
        lines := api.Group("/lines")
        {
            lines.POST("", lineHandler.CreateLine)
            lines.GET("", lineHandler.GetLines)
            lines.GET(":id", lineHandler.GetLine)
            lines.PUT(":id", lineHandler.UpdateLine)
            lines.DELETE(":id", lineHandler.DeleteLine)
        }
        
        // ЖД станции
        stations := api.Group("/stations")
        {
            stations.POST("", railStationHandler.CreateRailStation)
            stations.GET("", railStationHandler.GetRailStations)
            stations.GET(":id", railStationHandler.GetRailStation)
            stations.PUT(":id", railStationHandler.UpdateRailStation)
            stations.DELETE(":id", railStationHandler.DeleteRailStation)
            stations.GET("/geohash/:hash", railStationHandler.SearchByGeoHash)
            stations.GET("/radius", railStationHandler.SearchByRadius)
            stations.GET("/district/:district_id", railStationHandler.GetStationsByDistrict)
            stations.GET("/line/:line_id", railStationHandler.GetStationsByLine)
            stations.POST("/batch-update-geohash", railStationHandler.BatchUpdateGeohash)
            stations.GET("/nearest-meteostation/:id", railStationHandler.GetNearestMeteostation)
        }
        
        // Health check
        api.GET("/health", func(c *gin.Context) {
            var count int64
            if err := db.Model(&models.Region{}).Count(&count).Error; err != nil {
                c.JSON(500, gin.H{"status": "unhealthy", "error": err.Error()})
                return
            }
            c.JSON(200, gin.H{
                "status":  "healthy",
                "message": "Railway service is running",
                "version": "1.0.0",
            })
        })
    }
}