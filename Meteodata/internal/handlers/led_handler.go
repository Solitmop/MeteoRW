package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"meteodata/internal/models"
)

// LEDHandler обработчик для LED измерений
type LEDHandler struct {
	*BaseHandler
}

// NewLEDHandler создает новый обработчик для LED
func NewLEDHandler(db *gorm.DB) *LEDHandler {
	return &LEDHandler{
		BaseHandler: NewBaseHandler(db),
	}
}

// @Summary Создать LED измерение
// @Description Создает новое LED измерение
// @Tags led
// @Accept json
// @Produce json
// @Param measurement body models.LEDRequest true "Данные измерения"
// @Success 201 {object} models.LED
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /led [post]
func (h *LEDHandler) Create(c *gin.Context) {
	var req models.LEDRequest
	if !h.validateRequest(c, &req) {
		return
	}

	measurement := req.ToModel()

	result := h.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&measurement)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create LED measurement",
			"details": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, measurement)
}

// @Summary Получить LED измерения
// @Description Возвращает список LED измерений с фильтрацией
// @Tags led
// @Accept json
// @Produce json
// @Param index query string false "Индекс станции"
// @Param date_from query int false "Дата начала (Unix timestamp)"
// @Param date_to query int false "Дата окончания (Unix timestamp)"
// @Param limit query int false "Лимит записей" default(100)
// @Param offset query int false "Смещение" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /led [get]
func (h *LEDHandler) Get(c *gin.Context) {
	var filter models.BaseFilter

	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var measurements []models.LED
	query := h.buildBaseQuery(filter, &models.LED{})
	result := query.Find(&measurements)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":   measurements,
		"count":  len(measurements),
		"limit":  filter.Limit,
		"offset": filter.Offset,
	})
}

// @Summary Получить конкретное LED измерение
// @Description Возвращает LED измерение по индексу и дате
// @Tags led
// @Accept json
// @Produce json
// @Param index path int true "Индекс станции"
// @Param date path int true "Дата измерения (Unix timestamp)"
// @Success 200 {object} models.LED
// @Failure 404 {object} map[string]string
// @Router /led/{index}/{date} [get]
func (h *LEDHandler) GetByID(c *gin.Context) {
	indexStr := c.Param("index")
	dateStr := c.Param("date")

	index, err := strconv.Atoi(indexStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid index"})
		return
	}

	dateUnix, err := strconv.ParseInt(dateStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, expected Unix timestamp"})
		return
	}

	date := time.Unix(dateUnix, 0)
	var measurement models.LED

	if !h.checkMeasurementExists(c, &measurement, index, date) {
		return
	}

	c.JSON(http.StatusOK, measurement)
}

// @Summary Обновить LED измерение
// @Description Обновляет существующее LED измерение
// @Tags led
// @Accept json
// @Produce json
// @Param index path int true "Индекс станции"
// @Param date path int true "Дата измерения (Unix timestamp)"
// @Param measurement body models.LEDRequest true "Обновленные данные"
// @Success 200 {object} models.LED
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /led/{index}/{date} [put]
func (h *LEDHandler) Update(c *gin.Context) {
	indexStr := c.Param("index")
	dateStr := c.Param("date")

	index, err := strconv.Atoi(indexStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid index"})
		return
	}

	dateUnix, err := strconv.ParseInt(dateStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, expected Unix timestamp"})
		return
	}

	date := time.Unix(dateUnix, 0)

	// Проверяем существование
	var existing models.LED
	if !h.checkMeasurementExists(c, &existing, index, date) {
		return
	}

	var req models.LEDRequest
	if !h.validateRequest(c, &req) {
		return
	}

	// Обновляем
	updated := req.ToModel()
	updated.Index = index
	updated.Date = time.Unix(req.Date, 0)

	result := h.DB.Save(&updated)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// @Summary Удалить LED измерение
// @Description Удаляет LED измерение по индексу и дате
// @Tags led
// @Accept json
// @Produce json
// @Param index path int true "Индекс станции"
// @Param date path int true "Дата измерения (Unix timestamp)"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /led/{index}/{date} [delete]
func (h *LEDHandler) Delete(c *gin.Context) {
	indexStr := c.Param("index")
	dateStr := c.Param("date")

	index, err := strconv.Atoi(indexStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid index"})
		return
	}

	dateUnix, err := strconv.ParseInt(dateStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format, expected Unix timestamp"})
		return
	}

	date := time.Unix(dateUnix, 0)

	result := h.DB.Where("index = ? AND date = ?", index, date).Delete(&models.LED{})

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Measurement not found"})
		return
	}

	c.Status(http.StatusNoContent)
}
