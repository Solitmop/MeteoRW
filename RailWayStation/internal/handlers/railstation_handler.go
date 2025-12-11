// internal/handlers/railstation_handler.go
package handlers

import (
	"Railwaystation/internal/models"
	"Railwaystation/internal/service"
	"Railwaystation/pkg/validator"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RailStationHandler struct {
	DB         *gorm.DB
	GeoHashSvc *service.GeoHashService
}

func NewRailStationHandler(db *gorm.DB, geoSvc *service.GeoHashService) *RailStationHandler {
	return &RailStationHandler{
		DB:         db,
		GeoHashSvc: geoSvc,
	}
}

// CreateRailStation создает новую жд станцию
// @Summary Create a new railway station
// @Description Create a new railway station with automatic geohash generation
// @Tags stations
// @Accept json
// @Produce json
// @Param station body models.RailStationCreateRequest true "Railway station object"
// @Success 201 {object} models.RailStation
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /stations [post]
func (h *RailStationHandler) CreateRailStation(c *gin.Context) {
	var input models.RailStationCreateRequest

	if err := validator.BindAndValidate(c, &input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": validator.FormatValidationError(err),
		})
		return
	}

	// Проверяем существование district
	var district models.District
	if err := h.DB.First(&district, input.DistrictID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid district",
				"message": "District with specified ID does not exist",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": err.Error(),
		})
		return
	}

	// Проверяем существование line
	var line models.Line
	if err := h.DB.Where("id = ?", input.LineID).First(&line).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid line",
				"message": "Line with specified ID does not exist",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": err.Error(),
		})
		return
	}

	// Генерируем geohash
	hash, err := h.GeoHashSvc.Generate(input.Lat, input.Lon)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid coordinates",
			"message": err.Error(),
		})
		return
	}

	station := models.RailStation{
		ID:         input.ID,
		Name:       input.Name,
		Lat:        input.Lat,
		Lon:        input.Lon,
		DistrictID: input.DistrictID,
		LineID:     input.LineID,
		Hash:       hash,
	}

	result := h.DB.Create(&station)
	if result.Error != nil {
		if isDuplicateKeyError(result.Error) {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Rail station already exists",
				"message": "A station with similar data already exists",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": result.Error.Error(),
		})
		return
	}

	// Загружаем связанные данные
	h.DB.Preload("District").Preload("Line").First(&station, station.ID)

	c.JSON(http.StatusCreated, station)
}

// GetRailStation возвращает станцию по ID
// @Summary Get railway station by ID
// @Description Get a single railway station by its ID
// @Tags stations
// @Produce json
// @Param id path int true "Station ID"
// @Success 200 {object} models.RailStation
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /stations/{id} [get]
func (h *RailStationHandler) GetRailStation(c *gin.Context) {
	id := c.Param("id")

	stationID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID",
			"message": "ID must be a valid number",
		})
		return
	}

	var station models.RailStation
	result := h.DB.Preload("District").Preload("Line").First(&station, stationID)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Rail station not found",
				"message": "Station with specified ID does not exist",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, station)
}

// GetRailStations возвращает все станции с пагинацией
// @Summary Get all railway stations
// @Description Get a list of all railway stations with pagination and filtering
// @Tags stations
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(20)
// @Param district_id query int false "Filter by district ID"
// @Param line_id query string false "Filter by line ID"
// @Param name query string false "Filter by station name"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /stations [get]
func (h *RailStationHandler) GetRailStations(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	var stations []models.RailStation
	var total int64

	// Подготовка запроса с предзагрузкой связанных данных
	query := h.DB.Preload("District", func(db *gorm.DB) *gorm.DB {
		return db.Preload("Area", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Region")
		})
	}).Preload("Line")

	// Фильтрация по району
	if districtID := c.Query("district_id"); districtID != "" {
		query = query.Where("district_id = ?", districtID)
	}

	// Фильтрация по линии
	if lineID := c.Query("line_id"); lineID != "" {
		query = query.Where("line_id = ?", lineID)
	}

	// Поиск по имени
	if name := c.Query("name"); name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}

	// Подсчет общего количества
	query.Model(&models.RailStation{}).Count(&total)

	// Получение данных с пагинацией
	result := query.Offset(offset).Limit(limit).Find(&stations)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": result.Error.Error(),
		})
		return
	}

	response := gin.H{
		"data": stations,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": (int(total) + limit - 1) / limit,
		},
	}

	c.JSON(http.StatusOK, response)
}

// UpdateRailStation обновляет станцию
// @Summary Update railway station
// @Description Update an existing railway station by ID
// @Tags stations
// @Accept json
// @Produce json
// @Param id path int true "Station ID"
// @Param station body models.RailStationUpdateRequest true "Station update object"
// @Success 200 {object} models.RailStation
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /stations/{id} [put]
func (h *RailStationHandler) UpdateRailStation(c *gin.Context) {
	id := c.Param("id")

	stationID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID",
			"message": "ID must be a valid number",
		})
		return
	}

	// Проверяем существование станции
	var station models.RailStation
	if err := h.DB.First(&station, stationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Rail station not found",
				"message": "Station with specified ID does not exist",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": err.Error(),
		})
		return
	}

	var input models.RailStationUpdateRequest
	if err := validator.BindAndValidate(c, &input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": validator.FormatValidationError(err),
		})
		return
	}

	// Проверяем district если передан
	if input.DistrictID != nil {
		var district models.District
		if err := h.DB.First(&district, *input.DistrictID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Invalid district",
					"message": "District with specified ID does not exist",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Database error",
				"details": err.Error(),
			})
			return
		}
	}

	// Проверяем line если передан
	if input.LineID != nil {
		var line models.Line
		if err := h.DB.First(&line, *input.LineID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Invalid line",
					"message": "Line with specified ID does not exist",
				})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "Database error",
				"details": err.Error(),
			})
			return
		}
	}

	updates := make(map[string]interface{})

	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Lat != nil {
		updates["lat"] = *input.Lat
	}
	if input.Lon != nil {
		updates["lon"] = *input.Lon
	}
	if input.DistrictID != nil {
		updates["district_id"] = *input.DistrictID
	}
	if input.LineID != nil {
		updates["line_id"] = *input.LineID
	}

	// Если изменились координаты, пересчитываем geohash
	if input.Lat != nil || input.Lon != nil {
		lat := station.Lat
		lon := station.Lon

		if input.Lat != nil {
			lat = *input.Lat
		}
		if input.Lon != nil {
			lon = *input.Lon
		}

		hash, err := h.GeoHashSvc.Generate(lat, lon)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid coordinates",
				"message": err.Error(),
			})
			return
		}
		updates["hash"] = hash
	}

	result := h.DB.Model(&station).Updates(updates)
	if result.Error != nil {
		if isDuplicateKeyError(result.Error) {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Update conflict",
				"message": "Update would create duplicate data",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": result.Error.Error(),
		})
		return
	}

	// Загружаем обновленные данные с связями
	h.DB.Preload("District").Preload("Line").First(&station, stationID)

	c.JSON(http.StatusOK, station)
}

// DeleteRailStation удаляет станцию
// @Summary Delete railway station
// @Description Delete a railway station by ID
// @Tags stations
// @Produce json
// @Param id path int true "Station ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /stations/{id} [delete]
func (h *RailStationHandler) DeleteRailStation(c *gin.Context) {
	id := c.Param("id")

	stationID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID",
			"message": "ID must be a valid number",
		})
		return
	}

	var station models.RailStation
	if err := h.DB.First(&station, stationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Rail station not found",
				"message": "Station with specified ID does not exist",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": err.Error(),
		})
		return
	}

	result := h.DB.Delete(&station)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Rail station deleted successfully",
		"id":      stationID,
	})
}

// SearchByGeoHash ищет станции по geohash и соседним хешам
// @Summary Search stations by geohash
// @Description Search railway stations by geohash with optional neighbor inclusion
// @Tags stations
// @Produce json
// @Param hash path string true "Geohash string"
// @Param neighbors query boolean false "Include neighboring geohashes" default(false)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /stations/geohash/{hash} [get]
func (h *RailStationHandler) SearchByGeoHash(c *gin.Context) {
	hash := c.Param("hash")
	includeNeighbors := c.DefaultQuery("neighbors", "false") == "true"

	var hashesToSearch []string
	hashesToSearch = append(hashesToSearch, hash)

	// Добавляем соседние хеши если нужно
	if includeNeighbors {
		neighbors, err := h.GeoHashSvc.GenerateNeighbors(hash)
		if err == nil {
			hashesToSearch = append(hashesToSearch, neighbors...)
		}
	}

	var stations []models.RailStation
	result := h.DB.Preload("District").Preload("Line").
		Where("hash IN ?", hashesToSearch).
		Find(&stations)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"hash":     hash,
		"stations": stations,
		"count":    len(stations),
	})
}

// SearchByRadius ищет станции в радиусе от точки
// @Summary Search stations by radius
// @Description Search railway stations within a radius from a given point
// @Tags stations
// @Produce json
// @Param lat query number true "Latitude of the center point"
// @Param lon query number true "Longitude of the center point"
// @Param radius query number false "Search radius in kilometers" default(1)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /stations/radius [get]
func (h *RailStationHandler) SearchByRadius(c *gin.Context) {
	latStr := c.Query("lat")
	lonStr := c.Query("lon")
	radiusStr := c.DefaultQuery("radius", "1")

	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid latitude"})
		return
	}

	lon, err := strconv.ParseFloat(lonStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid longitude"})
		return
	}

	radius, err := strconv.ParseFloat(radiusStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid radius"})
		return
	}

	// Генерируем geohash для центральной точки
	hash, err := h.GeoHashSvc.Generate(lat, lon)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Получаем соседей
	neighbors, err := h.GeoHashSvc.GenerateNeighbors(hash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Ищем станции в этом и соседних геохешах
	var allStations []models.RailStation
	hashesToSearch := append(neighbors, hash)

	result := h.DB.Preload("District").Preload("Line").
		Where("hash IN ?", hashesToSearch).
		Find(&allStations)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": result.Error.Error(),
		})
		return
	}

	// Фильтруем по реальному расстоянию (в километрах)
	type StationWithDistance struct {
		models.RailStation
		Distance float64 `json:"distance_km"`
	}

	var filteredStations []StationWithDistance
	for _, station := range allStations {
		distance := calculateDistance(lat, lon, station.Lat, station.Lon)
		if distance <= radius {
			filteredStations = append(filteredStations, StationWithDistance{
				RailStation: station,
				Distance:    distance,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"center":  gin.H{"lat": lat, "lon": lon},
		"radius":  radius,
		"results": filteredStations,
		"count":   len(filteredStations),
	})
}

// GetStationsByDistrict возвращает станции по району
// @Summary Get stations by district
// @Description Get all railway stations belonging to a specific district
// @Tags stations
// @Produce json
// @Param district_id path int true "District ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /stations/district/{district_id} [get]
func (h *RailStationHandler) GetStationsByDistrict(c *gin.Context) {
	districtIDStr := c.Param("district_id")

	districtID, err := validator.ValidateID(districtIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid district ID",
			"message": err.Error(),
		})
		return
	}

	// Проверяем существование района
	var district models.District
	if err := h.DB.Preload("Area.Region").First(&district, districtID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "District not found",
				"message": "District with specified ID does not exist",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": err.Error(),
		})
		return
	}

	var stations []models.RailStation
	result := h.DB.Where("district_id = ?", districtID).
		Preload("District").Preload("Line").
		Find(&stations)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"district": district,
		"stations": stations,
		"count":    len(stations),
	})
}

// GetStationsByLine возвращает станции по линии
// @Summary Get stations by line
// @Description Get all railway stations belonging to a specific line
// @Tags stations
// @Produce json
// @Param line_id path string true "Line ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /stations/line/{line_id} [get]
func (h *RailStationHandler) GetStationsByLine(c *gin.Context) {
	lineID := c.Param("line_id")

	// Проверяем существование линии
	var line models.Line
	if err := h.DB.First(&line, lineID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Line not found",
				"message": "Line with specified ID does not exist",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": err.Error(),
		})
		return
	}

	var stations []models.RailStation
	result := h.DB.Where("line_id = ?", lineID).
		Preload("District").Preload("Line").
		Find(&stations)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"line":     line,
		"stations": stations,
		"count":    len(stations),
	})
}

// BatchUpdateGeohash обновляет geohash для всех станций
// @Summary Batch update geohash
// @Description Update geohash for all railway stations (useful for data migration)
// @Tags stations
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /stations/batch-update-geohash [post]
func (h *RailStationHandler) BatchUpdateGeohash(c *gin.Context) {
	var stations []models.RailStation
	result := h.DB.Find(&stations)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": result.Error.Error(),
		})
		return
	}

	updatedCount := 0
	failedCount := 0

	for _, station := range stations {
		geohashStr, err := h.GeoHashSvc.Generate(station.Lat, station.Lon)
		if err == nil {
			h.DB.Model(&station).Update("hash", geohashStr)
			updatedCount++
		} else {
			failedCount++
			log.Printf("Failed to generate geohash for station %d: %v", station.ID, err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Batch geohash update completed",
		"total":   len(stations),
		"updated": updatedCount,
		"failed":  failedCount,
	})
}

// GetNearestMeteostationHandler возвращает ближайшую метеостанцию для железнодорожной станции
func (h *RailStationHandler) GetNearestMeteostation(c *gin.Context) {
	stationID := c.Param("id")

	var railStation models.RailStation

	if err := h.DB.First(&railStation, stationID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Railway station not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": err.Error(),
		})
		return
	}

	geohashStr := railStation.Hash

	// Вызываем внешний сервис для получения метеостанций
	meteostations, err := h.fetchMeteostationsByGeohash(geohashStr[:3]) // Берем 4 символа
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch meteostations",
			"details": err.Error(),
		})
		return
	}

	if len(meteostations) == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message":     "No meteostations found near the railway station",
			"coordinates": [2]float64{railStation.Lat, railStation.Lon},
		})
		return
	}

	// Находим ближайшую метеостанцию
	nearestMeteo := h.findNearestMeteostation(
		railStation.Lat,
		railStation.Lon,
		meteostations,
	)

	c.JSON(http.StatusOK, gin.H{
		"id":           railStation.ID,
		"name":         railStation.Name,
		"latitude":     railStation.Lat,
		"longitude":    railStation.Lon,
		"meteostation": nearestMeteo.Index,
		"distance":     nearestMeteo.Distance,
	})
}

// Вспомогательная функция для проверки дублирования ключей
func isDuplicateKeyError(err error) bool {
	errorStr := err.Error()
	return strings.Contains(errorStr, "23505") ||
		strings.Contains(errorStr, "duplicate key") ||
		strings.Contains(errorStr, "unique constraint")
}

// Вспомогательная функция для расчета расстояния (Haversine formula)
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	// Реализация формулы Haversine
	const R = 6371 // Earth radius in km

	dLat := (lat2 - lat1) * (3.14159265358979323846 / 180.0)
	dLon := (lon2 - lon1) * (3.14159265358979323846 / 180.0)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*(3.14159265358979323846/180.0))*
			math.Cos(lat2*(3.14159265358979323846/180.0))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

type Meteostation struct {
	Index     uint    `json:"index"`
	Latitude  float64 `json:"latitude"`  // широта
	Longitude float64 `json:"longitude"` // долгота
	Distance float64 `json:"distance_km"`
}

// GeoHashResponse структура ответа от внешнего сервиса
type GeoHashResponse struct {
	Count    int            `json:"count"`
	Geohash  string         `json:"geohash"`
	Stations []Meteostation `json:"stations"`
}

// fetchMeteostationsByGeohash получает метеостанции из внешнего сервиса
func (h *RailStationHandler) fetchMeteostationsByGeohash(geohash string) ([]Meteostation, error) {
	// Формируем URL для запроса
	url := fmt.Sprintf("http://localhost:8081/api/geohash/%s?neighbors=true", geohash)

	// Создаем HTTP клиент с таймаутом
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// Выполняем GET запрос
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("external service returned status: %d", resp.StatusCode)
	}

	// Декодируем ответ
	var geoHashResp GeoHashResponse
	if err := json.NewDecoder(resp.Body).Decode(&geoHashResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return geoHashResp.Stations, nil
}

// findNearestMeteostation находит ближайшую метеостанцию
func (h *RailStationHandler) findNearestMeteostation(lat, lon float64,
	stations []Meteostation) *Meteostation {

	if len(stations) == 0 {
		return nil
	}

	// Инициализируем первую станцию как ближайшую
	nearest := stations[0]
	nearest.Distance = calculateDistance(lat, lon, stations[0].Latitude, stations[0].Longitude)

	// Проходим по остальным станциям
	for _, station := range stations[1:] {
		distance := calculateDistance(lat, lon, station.Latitude, station.Longitude)

		if distance < nearest.Distance {
			nearest = station
			nearest.Distance = distance
		}
	}

	return &nearest
}
