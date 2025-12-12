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

// RegularHandler обработчик для регулярных измерений
type RegularHandler struct {
	*BaseHandler
}

// NewRegularHandler создает новый обработчик для Regular
func NewRegularHandler(db *gorm.DB) *RegularHandler {
	return &RegularHandler{
		BaseHandler: NewBaseHandler(db),
	}
}

// @Summary Создать регулярное измерение
// @Description Создает новое регулярное метеорологическое измерение
// @Tags regular
// @Accept json
// @Produce json
// @Param measurement body models.RegularRequest true "Данные измерения"
// @Success 201 {object} models.Regular
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /regular [post]
func (h *RegularHandler) Create(c *gin.Context) {
	var req models.RegularRequest
	if !h.validateRequest(c, &req) {
		return
	}

	measurement := req.ToModel()

	result := h.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&measurement)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create regular measurement",
			"details": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, measurement)
}

// @Summary Получить регулярные измерения
// @Description Возвращает список регулярных измерений с фильтрацией
// @Tags regular
// @Accept json
// @Produce json
// @Param index query string false "Индекс станции"
// @Param date_from query int false "Дата начала (Unix timestamp)"
// @Param date_to query int false "Дата окончания (Unix timestamp)"
// @Param limit query int false "Лимит записей" default(100)
// @Param offset query int false "Смещение" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /regular [get]
func (h *RegularHandler) Get(c *gin.Context) {
	var filter models.BaseFilter

	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var measurements []models.Regular
	query := h.buildBaseQuery(filter, &models.Regular{})
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

// @Summary Получить конкретное регулярное измерение
// @Description Возвращает регулярное измерение по индексу и дате
// @Tags regular
// @Accept json
// @Produce json
// @Param index path int true "Индекс станции"
// @Param date path int true "Дата измерения (Unix timestamp)"
// @Success 200 {object} models.Regular
// @Failure 404 {object} map[string]string
// @Router /regular/{index}/{date} [get]
func (h *RegularHandler) GetByID(c *gin.Context) {
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
	var measurement models.Regular

	if !h.checkMeasurementExists(c, &measurement, index, date) {
		return
	}

	c.JSON(http.StatusOK, measurement)
}

// @Summary Обновить регулярное измерение
// @Description Обновляет существующее регулярное измерение
// @Tags regular
// @Accept json
// @Produce json
// @Param index path int true "Индекс станции"
// @Param date path int true "Дата измерения (Unix timestamp)"
// @Param measurement body models.RegularRequest true "Обновленные данные"
// @Success 200 {object} models.Regular
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /regular/{index}/{date} [put]
func (h *RegularHandler) Update(c *gin.Context) {
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
	var existing models.Regular
	if !h.checkMeasurementExists(c, &existing, index, date) {
		return
	}

	var req models.RegularRequest
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

// @Summary Удалить регулярное измерение
// @Description Удаляет регулярное измерение по индексу и дате
// @Tags regular
// @Accept json
// @Produce json
// @Param index path int true "Индекс станции"
// @Param date path int true "Дата измерения (Unix timestamp)"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /regular/{index}/{date} [delete]
func (h *RegularHandler) Delete(c *gin.Context) {
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

	result := h.DB.Where("index = ? AND date = ?", index, date).Delete(&models.Regular{})

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
