package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"meteodata/internal/models"
	"meteodata/pkg/validators"
)

// TTTRHandler обработчик для TTTR измерений
type TTTRHandler struct {
	*BaseHandler
}

// NewTTTRHandler создает новый обработчик для TTTR
func NewTTTRHandler(db *gorm.DB) *TTTRHandler {
	return &TTTRHandler{
		BaseHandler: NewBaseHandler(db),
	}
}

// @Summary Создать TTTR измерение
// @Description Создает новое TTTR измерение
// @Tags tttr
// @Accept json
// @Produce json
// @Param measurement body models.TTTRRequest true "Данные измерения"
// @Success 201 {object} models.TTTR
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tttr [post]
func (h *TTTRHandler) Create(c *gin.Context) {
	var req models.TTTRRequest
	if !h.validateRequest(c, &req) {
		return
	}

	measurement := req.ToModel()

	result := h.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&measurement)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create TTTR measurement",
			"details": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, measurement)
}

// @Summary Получить TTTR измерения
// @Description Возвращает список TTTR измерений с фильтрацией
// @Tags tttr
// @Accept json
// @Produce json
// @Param index query string false "Индекс станции"
// @Param date_from query int false "Дата начала (Unix timestamp)"
// @Param date_to query int false "Дата окончания (Unix timestamp)"
// @Param limit query int false "Лимит записей" default(100)
// @Param offset query int false "Смещение" default(0)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Router /tttr [get]
func (h *TTTRHandler) Get(c *gin.Context) {
	var filter models.BaseFilter

	if err := c.ShouldBindQuery(&filter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var measurements []models.TTTR
	query := h.buildBaseQuery(filter, &models.TTTR{})
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

// @Summary Получить конкретное TTTR измерение
// @Description Возвращает TTTR измерение по индексу и дате
// @Tags tttr
// @Accept json
// @Produce json
// @Param index path int true "Индекс станции"
// @Param date path int true "Дата измерения (Unix timestamp)"
// @Success 200 {object} models.TTTR
// @Failure 404 {object} map[string]string
// @Router /tttr/{index}/{date} [get]
func (h *TTTRHandler) GetByID(c *gin.Context) {
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
	var measurement models.TTTR

	if !h.checkMeasurementExists(c, &measurement, index, date) {
		return
	}

	c.JSON(http.StatusOK, measurement)
}

// @Summary Удалить TTTR измерение
// @Description Удаляет TTTR измерение по индексу и дате
// @Tags tttr
// @Accept json
// @Produce json
// @Param index path int true "Индекс станции"
// @Param date path int true "Дата измерения (Unix timestamp)"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /tttr/{index}/{date} [delete]
func (h *TTTRHandler) Delete(c *gin.Context) {
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

	result := h.DB.Where("index = ? AND date = ?", index, date).Delete(&models.TTTR{})

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

// MonthlyTemperatureRequest структура запроса для средней температуры за месяц
type MonthlyTemperatureRequest struct {
	Month     int `form:"month" validate:"required,min=1,max=12" json:"month"`
	StartYear int `form:"start_year" validate:"required,min=1900,max=2100" json:"start_year"`
	EndYear   int `form:"end_year" validate:"required,min=1900,max=2100,gtefield=StartYear" json:"end_year"`
	Index     string `form:"index"` // Опционально: индекс станции (можно несколько через запятую)
}

// MonthlyTemperatureResponse структура ответа со средней температурой
type MonthlyTemperatureResponse struct {
	Month              int     `json:"month"`
	MonthName          string  `json:"month_name"`
	StartYear          int     `json:"start_year"`
	EndYear            int     `json:"end_year"`
	Indexes            []int   `json:"indexes,omitempty"`
	AverageTemperature float64 `json:"average_temperature"`
	MinTemperature     float64 `json:"min_temperature"`
	MaxTemperature     float64 `json:"max_temperature"`
	DaysCount          int64   `json:"days_count"`
	YearsCount         int     `json:"years_count"`
	YearsInRange       []int   `json:"years_in_range"`
}

// @Summary Получить среднюю температуру за месяц
// @Description Рассчитывает среднюю температуру за указанный месяц в диапазоне лет
// @Tags tttr
// @Accept json
// @Produce json
// @Param month query int true "Месяц (1-12)"
// @Param start_year query int true "Год начала диапазона"
// @Param end_year query int true "Год окончания диапазона"
// @Param index query string false "Индекс станции (можно несколько через запятую)"
// @Success 200 {object} MonthlyTemperatureResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tttr/monthly-temperature [get]
func (h *TTTRHandler) GetMonthlyTemperature(c *gin.Context) {
    var req MonthlyTemperatureRequest
    
    // Получаем параметры из запроса
    if err := c.ShouldBindQuery(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Валидация
    if err := validators.Validate(req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Получаем среднюю температуру
    result, err := h.calculateMonthlyTemperature(req.Month, req.StartYear, req.EndYear, req.Index)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Failed to calculate monthly temperature",
            "details": err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, result)
}

// calculateMonthlyTemperature вычисляет среднюю температуру за месяц
func (h *TTTRHandler) calculateMonthlyTemperature(month, startYear, endYear int, indexStr string) (*MonthlyTemperatureResponse, error) {
    // Формируем диапазон дат
    startDate := time.Date(startYear, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
    endDate := time.Date(endYear, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
    endDate = endDate.AddDate(0, 1, -1) // Последний день месяца
    
    var result struct {
        AverageTemp float64 `gorm:"column:avg_temp"`
        MinTemp     float64 `gorm:"column:min_temp"`
        MaxTemp     float64 `gorm:"column:max_temp"`
        DaysCount   int64   `gorm:"column:days_count"`
    }
    
    // Строим базовый запрос
    query := h.DB.Model(&models.TTTR{}).
        Select(`
            AVG(t_avg) as avg_temp,
            MIN(t_min) as min_temp,
            MAX(t_max) as max_temp,
            COUNT(*) as days_count
        `).
        Where("date >= ? AND date <= ?", startDate, endDate).
        Where("EXTRACT(MONTH FROM date) = ?", month).
        Where("EXTRACT(YEAR FROM date) BETWEEN ? AND ?", startYear, endYear)
    
    // Добавляем фильтр по индексу, если указан
    parsedIndexes := []int{}
    if indexStr != "" {
        indexes := h.parseIndexList(indexStr)
        parsedIndexes = indexes
        
        if len(indexes) > 0 {
            query = query.Where("index IN ?", indexes)
        }
    }
    
    if err := query.Scan(&result).Error; err != nil {
        return nil, err
    }
    
    // Получаем список лет в диапазоне, по которым есть данные
    var years []int
    yearsQuery := h.DB.Model(&models.TTTR{}).
        Select("DISTINCT EXTRACT(YEAR FROM date) as year").
        Where("date >= ? AND date <= ?", startDate, endDate).
        Where("EXTRACT(MONTH FROM date) = ?", month).
        Where("EXTRACT(YEAR FROM date) BETWEEN ? AND ?", startYear, endYear)
    
    // Добавляем фильтр по индексу в запрос для лет
    if len(parsedIndexes) > 0 {
        yearsQuery = yearsQuery.Where("index IN ?", parsedIndexes)
    }
    
    yearsQuery = yearsQuery.Order("year ASC")
    
    rows, err := yearsQuery.Rows()
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    
    for rows.Next() {
        var year int
        if err := rows.Scan(&year); err != nil {
            return nil, err
        }
        years = append(years, year)
    }
    
    // Получаем название месяца
    monthName := time.Month(month).String()
    
    response := &MonthlyTemperatureResponse{
        Month:              month,
        MonthName:          monthName,
        StartYear:          startYear,
        EndYear:            endYear,
        AverageTemperature: round(result.AverageTemp, 2),
        MinTemperature:     round(result.MinTemp, 2),
        MaxTemperature:     round(result.MaxTemp, 2),
        DaysCount:          result.DaysCount,
        YearsCount:         len(years),
        YearsInRange:       years,
    }
    
    // Добавляем индексы, если они были указаны
    if len(parsedIndexes) > 0 {
        response.Indexes = parsedIndexes
    }
    
    return response, nil
}

// Вспомогательная функция для округления
func round(value float64, precision int) float64 {
	if precision <= 0 {
		return value
	}

	multiplier := 1.0
	for i := 0; i < precision; i++ {
		multiplier *= 10
	}

	return float64(int(value*multiplier+0.5)) / multiplier
}
