package handlers

import (
	"Meteostation/models"
	"Meteostation/pkg/geoservice"
	"Meteostation/validators"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MeteostationHandler struct {
	DB         *gorm.DB
	GeoService *geoservice.GeoHashClient
}

// CreateMeteostation создает новую станцию
// @Summary Create a new meteostation
// @Description Create a new meteostation in the database
// @Tags meteostations
// @Accept json
// @Produce json
// @Param meteostation body models.Meteostation true "Meteostation object"
// @Success 201 {object} models.Meteostation
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /meteostations [post]
func (h *MeteostationHandler) Create(c *gin.Context) {
	var meteostation models.Meteostation
	// TODO: Add validation
	if err := c.ShouldBindJSON(&meteostation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Валидация
	if err := validators.Validate(meteostation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Автоматически генерируем geohash, если его нет, но есть координаты
	if meteostation.Geohash == "" {
		if meteostation.Longitude > 180 && meteostation.Longitude < 360 {
			meteostation.Longitude = meteostation.Longitude - 360
		}
		if meteostation.Latitude > 90 && meteostation.Longitude < 180 {
			meteostation.Longitude = meteostation.Longitude - 180
		}
		geohashStr, err := h.GeoService.GetGeohash(meteostation.Latitude, meteostation.Longitude)
		if err != nil {
			// Ошибка возможна только при неверных координатах (выход за границы)
			log.Printf("Failed to generate geohash: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid coordinates provided for geohash generation",
				"details": err.Error(),
			})
			return
		}
		meteostation.Geohash = geohashStr
		log.Printf("Geohash '%s' generated for coordinates (%f, %f)",
			geohashStr, meteostation.Latitude, meteostation.Longitude)
	}

	result := h.DB.Create(&meteostation)

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Meteostation already exists"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, meteostation)
}

// GetMeteostations возвращает все станции
// @Summary Get all meteostations
// @Description Get a list of all meteostations
// @Tags meteostations
// @Produce json
// @Success 200 {array} models.Meteostation
// @Failure 500 {object} map[string]string
// @Router /meteostations [get]
func (h *MeteostationHandler) Index(c *gin.Context) {
	var meteostations []models.Meteostation
	result := h.DB.Find(&meteostations)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, meteostations)
}

// getMeteostationByIndex внутренний метод для получения станции по Index
// Возвращает станцию и флаг существования (без отправки HTTP ответа)
func (h *MeteostationHandler) getByIndex(c *gin.Context) (*models.Meteostation, error) {
	index := c.Param("index")

	var meteostation models.Meteostation
	result := h.DB.First(&meteostation, index)
	return &meteostation, result.Error
}

// GetMeteostation возвращает станцию по Index
// @Summary Get meteostation by Index
// @Description Get a single meteostation by its Index
// @Tags meteostations
// @Produce json
// @Param index path int true "Meteostation Index"
// @Success 200 {object} models.Meteostation
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /meteostations/{index} [get]
func (h *MeteostationHandler) Get(c *gin.Context) {
	meteostation, err := h.getByIndex(c)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Meteostation not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, meteostation)
}

// UpdateMeteostation обновляет станцию по Index
// @Summary Update meteostation
// @Description Update an existing meteostation by Index
// @Tags meteostations
// @Accept json
// @Produce json
// @Param index path int true "Meteostation Index"
// @Param meteostation body models.Meteostation true "Meteostation object"
// @Success 200 {object} models.Meteostation
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /meteostations/{index} [put]
func (h *MeteostationHandler) Update(c *gin.Context) {
	meteostation, err := h.getByIndex(c)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Meteostation not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&meteostation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validators.Validate(meteostation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := h.DB.Save(&meteostation)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, meteostation)

}

// DeleteMeteostation удаляет станцию по Index
// @Summary Delete meteostation
// @Description Delete a meteostation by Index
// @Tags meteostations
// @Produce json
// @Param index path int true "Meteostation Index"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /meteostations/{index} [delete]
func (h *MeteostationHandler) Delete(c *gin.Context) {
	_, err := h.getByIndex(c)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Meteostation not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	index, err := strconv.ParseUint(c.Param("index"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meteostation Index"})
		return
	}
	result := h.DB.Delete(&models.Meteostation{}, index)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, index)
}

// SearchByGeoHash ищет станции по geohash и соседним хешам
// @Summary Search stations by geohash
// @Description Search railway stations by geohash with optional neighbor inclusion
// @Tags stations
// @Produce json
// @Param geohash path string true "Geohash string"
// @Param neighbors query boolean false "Include neighboring geohashes" default(false)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /stations/geohash/{geohash} [get]
func (h *MeteostationHandler) SearchByGeoHash(c *gin.Context) {
	limit := 100
	geohash := c.Param("geohash")
	includeNeighbors := c.DefaultQuery("neighbors", "false") == "true"

	var hashesToSearch []string
	hashesToSearch = append(hashesToSearch, geohash)

	// Добавляем соседние хеши если нужно
	if includeNeighbors {
		neighbors, err := h.GeoService.GenerateNeighbors(geohash)
		if err == nil {
			hashesToSearch = append(hashesToSearch, neighbors...)
		}
	}

	var stations []models.Meteostation

	// Создаем условия для поиска по префиксам
    var conditions []string
    var args []interface{}
	for _, searchHash := range hashesToSearch {
        conditions = append(conditions, "geohash LIKE ?")
        args = append(args, searchHash+"%")
    }
    
    // Собираем SQL запрос
    query := h.DB
    
    if len(conditions) > 0 {
        // Используем OR для всех префиксов
        query = query.Where(strings.Join(conditions, " OR "), args...)
    } else {
        // Поиск только по основному префиксу
        query = query.Where("geohash LIKE ?", geohash+"%")
    }
	
	// Выполняем запрос с лимитом
    result := query.Limit(limit).Find(&stations)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"geohash":  geohash,
		"stations": stations,
		"count":    len(stations),
	})
}