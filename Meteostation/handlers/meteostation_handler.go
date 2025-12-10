package handlers

import (
	"Meteostation/models"
	"Meteostation/pkg/geoservice"
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

// CreateMeteostation создает новый продукт
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

	// Автоматически генерируем geohash, если его нет, но есть координаты
	if meteostation.Geohash == "" {
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

// GetMeteostations возвращает все продукты
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

// getMeteostationByIndex внутренний метод для получения продукта по Index
// Возвращает продукт и флаг существования (без отправки HTTP ответа)
func (h *MeteostationHandler) getByIndex(c *gin.Context) (*models.Meteostation, error) {
	index := c.Param("index")
	// TODO: Add validation

	var meteostation models.Meteostation
	result := h.DB.First(&meteostation, index)
	return &meteostation, result.Error
}

// GetMeteostation возвращает продукт по Index
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
	//index := c.Param("index")
	// TODO: Add validation
	/*
		    if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meteostation Index"})
				return
			}
	*/

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

// UpdateMeteostation обновляет продукт по Index
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
	// TODO: Add validation

	if err := c.ShouldBindJSON(&meteostation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	/*index, err := strconv.ParseUint(c.Param("index"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meteostation Index"})
		return
	}
	updatedMeteostation.Index = uint(index)
	*/
	result := h.DB.Save(&meteostation)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, meteostation)

}

// DeleteMeteostation удаляет продукт по Index
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
