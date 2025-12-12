package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"meteodata/internal/models"
	"meteodata/pkg/validators"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BaseHandler содержит общие методы для всех обработчиков
type BaseHandler struct {
	DB *gorm.DB
}

// NewBaseHandler создает новый базовый обработчик
func NewBaseHandler(db *gorm.DB) *BaseHandler {
	return &BaseHandler{DB: db}
}

// buildBaseQuery строит базовый запрос с фильтрацией
func (h *BaseHandler) buildBaseQuery(filter models.BaseFilter, model interface{}) *gorm.DB {
	query := h.DB.Model(model)

	// Фильтр по индексам
	if filter.Index != "" {
		indices := h.parseIndexList(filter.Index)
		if len(indices) > 0 {
			query = query.Where("index IN ?", indices)
		}
	}

	// Фильтр по датам
	if filter.DateFrom > 0 {
		query = query.Where("date >= ?", time.Unix(filter.DateFrom, 0))
	}

	if filter.DateTo > 0 {
		query = query.Where("date <= ?", time.Unix(filter.DateTo, 0))
	}

	// Пагинация
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}

	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}

	// Сортировка
	query = query.Order("date DESC")

	return query
}

// parseIndexList парсит список индексов
func (h *BaseHandler) parseIndexList(indexStr string) []int {
	if indexStr == "" {
		return nil
	}

	parts := strings.Split(indexStr, ",")
	indices := make([]int, 0, len(parts))

	for _, part := range parts {
		if idx, err := strconv.Atoi(strings.TrimSpace(part)); err == nil {
			indices = append(indices, idx)
		}
	}

	return indices
}

// validateRequest валидирует запрос
func (h *BaseHandler) validateRequest(c *gin.Context, request interface{}) bool {
	if err := c.ShouldBindJSON(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}

	if err := validators.Validate(request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}

	return true
}

// checkMeasurementExists проверяет существование измерения
func (h *BaseHandler) checkMeasurementExists(c *gin.Context, model interface{}, index int, date time.Time) bool {
	result := h.DB.Where("index = ? AND date = ?", index, date).First(model)
	if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "Measurement not found"})
		return false
	}
	return true
}
