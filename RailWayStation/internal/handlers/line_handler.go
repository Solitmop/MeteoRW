// internal/handlers/line_handler.go
package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"Railwaystation/internal/models"
	"Railwaystation/pkg/validator"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LineHandler struct {
	DB *gorm.DB
}

func NewLineHandler(db *gorm.DB) *LineHandler {
	return &LineHandler{DB: db}
}

// CreateLine создает новую линию
// @Summary Create a new line
// @Description Create a new railway line in the database
// @Tags lines
// @Accept json
// @Produce json
// @Param line body models.LineCreateRequest true "Line object"
// @Success 201 {object} models.Line
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /lines [post]
func (h *LineHandler) CreateLine(c *gin.Context) {
	var input models.LineCreateRequest

	if err := validator.BindAndValidate(c, &input); err != nil {
		c.JSON(http.StatusBadRequest, validator.FormatError(err))
		return
	}

	line := models.Line{
		ID:   strings.TrimSpace(input.ID),
		Name: strings.TrimSpace(input.Name),
	}

	result := h.DB.Create(&line)
	if result.Error != nil {
		if validator.IsDuplicateKeyError(result.Error) {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Line already exists",
				"message": "A line with this ID or name already exists",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, line)
}

// GetLine возвращает линию по ID
// @Summary Get line by ID
// @Description Get a single railway line by its ID
// @Tags lines
// @Produce json
// @Param id path string true "Line ID"
// @Success 200 {object} models.Line
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /lines/{id} [get]
func (h *LineHandler) GetLine(c *gin.Context) {
	id := c.Param("id")

	var line models.Line
	if err := h.DB.Preload("RailStations").First(&line, id).Error; err != nil {
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

	c.JSON(http.StatusOK, line)
}

// GetLines возвращает все линии с пагинацией
// @Summary Get all lines
// @Description Get a list of all railway lines with pagination
// @Tags lines
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(20)
// @Param search query string false "Search by name or ID"
// @Param sort_by query string false "Sort field" default(id)
// @Param order query string false "Sort order" Enums(asc, desc) default(asc)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /lines [get]
func (h *LineHandler) GetLines(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	var lines []models.Line
	var total int64

	query := h.DB.Model(&models.Line{})

	// Поиск по имени или ID
	if search := c.Query("search"); search != "" {
		query = query.Where("name ILIKE ? OR id ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Сортировка
	sortBy := c.DefaultQuery("sort_by", "id")
	order := c.DefaultQuery("order", "asc")
	if order != "asc" && order != "desc" {
		order = "asc"
	}
	query = query.Order(sortBy + " " + order)

	// Подсчет и получение
	query.Count(&total)
	result := query.Offset(offset).Limit(limit).Find(&lines)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": result.Error.Error(),
		})
		return
	}

	response := gin.H{
		"data": lines,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": (int(total) + limit - 1) / limit,
		},
	}

	c.JSON(http.StatusOK, response)
}

// UpdateLine обновляет линию
// @Summary Update line
// @Description Update an existing railway line by ID
// @Tags lines
// @Accept json
// @Produce json
// @Param id path string true "Line ID"
// @Param line body models.LineUpdateRequest true "Line update object"
// @Success 200 {object} models.Line
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /lines/{id} [put]
func (h *LineHandler) UpdateLine(c *gin.Context) {
	id := c.Param("id")

	var line models.Line
	if err := h.DB.First(&line, id).Error; err != nil {
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

	var input models.LineUpdateRequest
	if err := validator.BindAndValidate(c, &input); err != nil {
		c.JSON(http.StatusBadRequest, validator.FormatError(err))
		return
	}

	updates := make(map[string]interface{})
	if input.Name != nil {
		updates["name"] = strings.TrimSpace(*input.Name)
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "No data to update",
			"message": "Provide at least one field to update",
		})
		return
	}

	result := h.DB.Model(&line).Updates(updates)
	if result.Error != nil {
		if validator.IsDuplicateKeyError(result.Error) {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Update conflict",
				"message": "A line with this name already exists",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": result.Error.Error(),
		})
		return
	}

	h.DB.First(&line, id)
	c.JSON(http.StatusOK, line)
}

// DeleteLine удаляет линию
// @Summary Delete line
// @Description Delete a railway line by ID
// @Tags lines
// @Produce json
// @Param id path string true "Line ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /lines/{id} [delete]
func (h *LineHandler) DeleteLine(c *gin.Context) {
	id := c.Param("id")

	var line models.Line
	if err := h.DB.First(&line, id).Error; err != nil {
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

	// Проверяем, есть ли связанные станции
	var stationsCount int64
	h.DB.Model(&models.RailStation{}).Where("line_id = ?", id).Count(&stationsCount)

	if stationsCount > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error":          "Cannot delete line",
			"message":        "Line has associated rail stations. Delete them first.",
			"stations_count": stationsCount,
		})
		return
	}

	result := h.DB.Delete(&line)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Line deleted successfully",
		"id":      id,
	})
}
