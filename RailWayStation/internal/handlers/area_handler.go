// internal/handlers/area_handler.go
package handlers

import (
	"Railwaystation/internal/models"
	"Railwaystation/pkg/validator"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AreaHandler struct {
	DB *gorm.DB
}

func NewAreaHandler(db *gorm.DB) *AreaHandler {
	return &AreaHandler{DB: db}
}

// CreateArea создает новую область
// @Summary Create a new area
// @Description Create a new area in the database
// @Tags areas
// @Accept json
// @Produce json
// @Param area body models.AreaCreateRequest true "Area object"
// @Success 201 {object} models.Area
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /areas [post]
func (h *AreaHandler) CreateArea(c *gin.Context) {
	var input models.AreaCreateRequest

	if err := validator.BindAndValidate(c, &input); err != nil {
		c.JSON(http.StatusBadRequest, validator.FormatError(err))
		return
	}

	// Проверяем существование региона
	var region models.Region
	if err := h.DB.First(&region, input.RegionID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Region not found",
				"message": "Region with specified ID does not exist",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": err.Error(),
		})
		return
	}

	area := models.Area{
		Name:     strings.TrimSpace(input.Name),
		RegionID: input.RegionID,
	}

	result := h.DB.Create(&area)
	if result.Error != nil {
		if validator.IsDuplicateKeyError(result.Error) {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Area already exists",
				"message": "An area with this name already exists in this region",
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
	h.DB.Preload("Region").First(&area, area.ID)

	c.JSON(http.StatusCreated, area)
}

// GetArea возвращает область по ID
// @Summary Get area by ID
// @Description Get a single area by its ID
// @Tags areas
// @Produce json
// @Param id path int true "Area ID"
// @Success 200 {object} models.Area
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /areas/{id} [get]
func (h *AreaHandler) GetArea(c *gin.Context) {
	id := c.Param("id")

	areaID, err := validator.ValidateID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID",
			"message": err.Error(),
		})
		return
	}

	var area models.Area
	if err := h.DB.Preload("Region").Preload("Districts").First(&area, areaID).Error; err != nil {
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

	c.JSON(http.StatusOK, area)
}

// GetAreas возвращает все области с фильтрацией
// @Summary Get all areas
// @Description Get a list of all areas with pagination and filtering
// @Tags areas
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(20)
// @Param region_id query int false "Filter by region ID"
// @Param name query string false "Filter by name"
// @Param sort_by query string false "Sort field" default(id)
// @Param order query string false "Sort order" Enums(asc, desc) default(asc)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /areas [get]
func (h *AreaHandler) GetAreas(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	var areas []models.Area
	var total int64

	query := h.DB.Preload("Region")

	// Фильтр по региону
	if regionID := c.Query("region_id"); regionID != "" {
		if id, err := strconv.ParseUint(regionID, 10, 32); err == nil {
			query = query.Where("region_id = ?", id)
		}
	}

	// Поиск по имени
	if name := c.Query("name"); name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}

	// Сортировка
	sortBy := c.DefaultQuery("sort_by", "id")
	order := c.DefaultQuery("order", "asc")
	if order != "asc" && order != "desc" {
		order = "asc"
	}
	query = query.Order(sortBy + " " + order)

	// Подсчет и получение
	query.Model(&models.Area{}).Count(&total)
	result := query.Offset(offset).Limit(limit).Find(&areas)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": result.Error.Error(),
		})
		return
	}

	response := gin.H{
		"data": areas,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": (int(total) + limit - 1) / limit,
		},
	}

	c.JSON(http.StatusOK, response)
}

// GetAreasByRegion возвращает области по региону
// @Summary Get areas by region
// @Description Get all areas belonging to a specific region
// @Tags areas
// @Produce json
// @Param region_id path int true "Region ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /areas/region/{region_id} [get]
func (h *AreaHandler) GetAreasByRegion(c *gin.Context) {
	regionIDStr := c.Param("region_id")

	regionID, err := validator.ValidateID(regionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid region ID",
			"message": err.Error(),
		})
		return
	}

	// Проверяем существование региона
	var region models.Region
	if err := h.DB.First(&region, regionID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Region not found",
				"message": "Region with specified ID does not exist",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": err.Error(),
		})
		return
	}

	var areas []models.Area
	result := h.DB.Where("region_id = ?", regionID).Preload("Region").Find(&areas)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"region": region,
		"areas":  areas,
		"count":  len(areas),
	})
}

// UpdateArea обновляет область
// @Summary Update area
// @Description Update an existing area by ID
// @Tags areas
// @Accept json
// @Produce json
// @Param id path int true "Area ID"
// @Param area body models.AreaUpdateRequest true "Area update object"
// @Success 200 {object} models.Area
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /areas/{id} [put]
func (h *AreaHandler) UpdateArea(c *gin.Context) {
	id := c.Param("id")

	areaID, err := validator.ValidateID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID",
			"message": err.Error(),
		})
		return
	}

	var area models.Area
	if err := h.DB.First(&area, areaID).Error; err != nil {
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

	var input models.AreaUpdateRequest
	if err := validator.BindAndValidate(c, &input); err != nil {
		c.JSON(http.StatusBadRequest, validator.FormatError(err))
		return
	}

	// Проверяем новый регион если указан
	if input.RegionID != nil {
		var region models.Region
		if err := h.DB.First(&region, *input.RegionID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Region not found",
					"message": "New region with specified ID does not exist",
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
	if input.RegionID != nil {
		updates["region_id"] = *input.RegionID
	}

	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "No data to update",
			"message": "Provide at least one field to update",
		})
		return
	}

	result := h.DB.Model(&area).Updates(updates)
	if result.Error != nil {
		if validator.IsDuplicateKeyError(result.Error) {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Update conflict",
				"message": "An area with this name already exists in the region",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": result.Error.Error(),
		})
		return
	}

	h.DB.Preload("Region").First(&area, areaID)
	c.JSON(http.StatusOK, area)
}

// DeleteArea удаляет область
// @Summary Delete area
// @Description Delete an area by ID
// @Tags areas
// @Produce json
// @Param id path int true "Area ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /areas/{id} [delete]
func (h *AreaHandler) DeleteArea(c *gin.Context) {
	id := c.Param("id")

	areaID, err := validator.ValidateID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid ID",
			"message": err.Error(),
		})
		return
	}

	var area models.Area
	if err := h.DB.First(&area, areaID).Error; err != nil {
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

	// Проверяем, есть ли связанные районы
	var districtsCount int64
	h.DB.Model(&models.District{}).Where("area_id = ?", areaID).Count(&districtsCount)

	if districtsCount > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error":           "Cannot delete area",
			"message":         "Area has associated districts. Delete them first.",
			"districts_count": districtsCount,
		})
		return
	}

	result := h.DB.Delete(&area)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"details": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Area deleted successfully",
		"id":      areaID,
	})
}
