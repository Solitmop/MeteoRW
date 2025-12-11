// internal/handlers/district_handler.go
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

type DistrictHandler struct {
	DB *gorm.DB
}

func NewDistrictHandler(db *gorm.DB) *DistrictHandler {
	return &DistrictHandler{DB: db}
}

// CreateDistrict создает новый район
// @Summary Create a new district
// @Description Create a new district in the database
// @Tags districts
// @Accept json
// @Produce json
// @Param district body models.DistrictCreateRequest true "District object"
// @Success 201 {object} models.District
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /districts [post]
func (h *DistrictHandler) CreateDistrict(c *gin.Context) {
	var input models.DistrictCreateRequest

	if err := validator.BindAndValidate(c, &input); err != nil {
		c.JSON(http.StatusBadRequest, validator.FormatError(err))
		return
	}

	// Проверяем существование области
	var area models.Area
	if err := h.DB.First(&area, input.AreaID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Area not found",
				"message": "Area with specified ID does not exist",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": err.Error(),
		})
		return
	}

	district := models.District{
		Name:   strings.TrimSpace(input.Name),
		AreaID: input.AreaID,
	}

	result := h.DB.Create(&district)
	if result.Error != nil {
		if validator.IsDuplicateKeyError(result.Error) {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "District already exists",
				"message": "A district with this name already exists in this area",
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
	h.DB.Preload("Area.Region").First(&district, district.ID)

	c.JSON(http.StatusCreated, district)
}

// GetDistrict возвращает район по ID
// @Summary Get district by ID
// @Description Get a single district by its ID
// @Tags districts
// @Produce json
// @Param id path int true "District ID"
// @Success 200 {object} models.District
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /districts/{id} [get]
func (h *DistrictHandler) GetDistrict(c *gin.Context) {
	id := c.Param("id")

	districtID, err := validator.ValidateID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID",
			"message": err.Error(),
		})
		return
	}

	var district models.District
	if err := h.DB.Preload("Area.Region").Preload("RailStations").First(&district, districtID).Error; err != nil {
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

	c.JSON(http.StatusOK, district)
}

// GetDistricts возвращает все районы с фильтрацией
// @Summary Get all districts
// @Description Get a list of all districts with pagination and filtering
// @Tags districts
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(20)
// @Param area_id query int false "Filter by area ID"
// @Param region_id query int false "Filter by region ID"
// @Param name query string false "Filter by name"
// @Param sort_by query string false "Sort field" default(id)
// @Param order query string false "Sort order" Enums(asc, desc) default(asc)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /districts [get]
func (h *DistrictHandler) GetDistricts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	var districts []models.District
	var total int64

	query := h.DB.Preload("Area.Region")

	// Фильтр по области
	if areaID := c.Query("area_id"); areaID != "" {
		if id, err := strconv.ParseUint(areaID, 10, 32); err == nil {
			query = query.Where("area_id = ?", id)
		}
	}

	// Фильтр по региону через область
	if regionID := c.Query("region_id"); regionID != "" {
		query = query.Joins("JOIN areas ON areas.id = districts.area_id").
			Where("areas.region_id = ?", regionID)
	}

	// Поиск по имени
	if name := c.Query("name"); name != "" {
		query = query.Where("districts.name ILIKE ?", "%"+name+"%")
	}

	// Сортировка
	sortBy := c.DefaultQuery("sort_by", "id")
	order := c.DefaultQuery("order", "asc")
	if order != "asc" && order != "desc" {
		order = "asc"
	}
	query = query.Order("districts." + sortBy + " " + order)

	// Подсчет и получение
	query.Model(&models.District{}).Count(&total)
	result := query.Offset(offset).Limit(limit).Find(&districts)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": result.Error.Error(),
		})
		return
	}

	response := gin.H{
		"data": districts,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": (int(total) + limit - 1) / limit,
		},
	}

	c.JSON(http.StatusOK, response)
}

// GetDistrictsByArea возвращает районы по области
// @Summary Get districts by area
// @Description Get all districts belonging to a specific area
// @Tags districts
// @Produce json
// @Param area_id path int true "Area ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /districts/area/{area_id} [get]
func (h *DistrictHandler) GetDistrictsByArea(c *gin.Context) {
	areaIDStr := c.Param("area_id")

	areaID, err := validator.ValidateID(areaIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid area ID",
			"message": err.Error(),
		})
		return
	}

	// Проверяем существование области
	var area models.Area
	if err := h.DB.Preload("Region").First(&area, areaID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Area not found",
				"message": "Area with specified ID does not exist",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": err.Error(),
		})
		return
	}

	var districts []models.District
	result := h.DB.Where("area_id = ?", areaID).Preload("Area.Region").Find(&districts)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"area":      area,
		"districts": districts,
		"count":     len(districts),
	})
}

// UpdateDistrict обновляет район
// @Summary Update district
// @Description Update an existing district by ID
// @Tags districts
// @Accept json
// @Produce json
// @Param id path int true "District ID"
// @Param district body models.DistrictUpdateRequest true "District update object"
// @Success 200 {object} models.District
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /districts/{id} [put]
func (h *DistrictHandler) UpdateDistrict(c *gin.Context) {
	id := c.Param("id")

	districtID, err := validator.ValidateID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID",
			"message": err.Error(),
		})
		return
	}

	var district models.District
	if err := h.DB.First(&district, districtID).Error; err != nil {
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

	var input models.DistrictUpdateRequest
	if err := validator.BindAndValidate(c, &input); err != nil {
		c.JSON(http.StatusBadRequest, validator.FormatError(err))
		return
	}

	// Проверяем новую область если указана
	if input.AreaID != nil {
		var area models.Area
		if err := h.DB.First(&area, *input.AreaID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Area not found",
					"message": "New area with specified ID does not exist",
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
		updates["name"] = strings.TrimSpace(*input.Name)
	}
	if input.AreaID != nil {
		updates["area_id"] = *input.AreaID
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "No data to update",
			"message": "Provide at least one field to update",
		})
		return
	}

	result := h.DB.Model(&district).Updates(updates)
	if result.Error != nil {
		if validator.IsDuplicateKeyError(result.Error) {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Update conflict",
				"message": "A district with this name already exists in the area",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": result.Error.Error(),
		})
		return
	}

	h.DB.Preload("Area.Region").First(&district, districtID)
	c.JSON(http.StatusOK, district)
}

// DeleteDistrict удаляет район
// @Summary Delete district
// @Description Delete a district by ID
// @Tags districts
// @Produce json
// @Param id path int true "District ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /districts/{id} [delete]
func (h *DistrictHandler) DeleteDistrict(c *gin.Context) {
	id := c.Param("id")

	districtID, err := validator.ValidateID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID",
			"message": err.Error(),
		})
		return
	}

	var district models.District
	if err := h.DB.First(&district, districtID).Error; err != nil {
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

	// Проверяем, есть ли связанные станции
	var stationsCount int64
	h.DB.Model(&models.RailStation{}).Where("district_id = ?", districtID).Count(&stationsCount)

	if stationsCount > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error":          "Cannot delete district",
			"message":        "District has associated rail stations. Delete them first.",
			"stations_count": stationsCount,
		})
		return
	}

	result := h.DB.Delete(&district)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "District deleted successfully",
		"id":      districtID,
	})
}
