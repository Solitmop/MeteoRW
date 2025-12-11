// internal/handlers/region_handler.go
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

type RegionHandler struct {
    DB *gorm.DB
}

func NewRegionHandler(db *gorm.DB) *RegionHandler {
    return &RegionHandler{DB: db}
}

// CreateRegion создает новый регион
// @Summary Create a new region
// @Description Create a new region in the database
// @Tags regions
// @Accept json
// @Produce json
// @Param region body models.RegionCreateRequest true "Region object"
// @Success 201 {object} models.Region
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /regions [post]
func (h *RegionHandler) CreateRegion(c *gin.Context) {
    var input models.RegionCreateRequest
    
    if err := validator.BindAndValidate(c, &input); err != nil {
        c.JSON(http.StatusBadRequest, validator.FormatError(err))
        return
    }
    
    region := models.Region{
        Name: strings.TrimSpace(input.Name),
    }
    
    result := h.DB.Create(&region)
    if result.Error != nil {
        if validator.IsDuplicateKeyError(result.Error) {
            c.JSON(http.StatusConflict, gin.H{
                "error":   "Region already exists",
                "message": "A region with this name already exists",
            })
            return
        }
        
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Database error",
            "details": result.Error.Error(),
        })
        return
    }
    
    c.JSON(http.StatusCreated, region)
}

// GetRegion возвращает регион по ID
// @Summary Get region by ID
// @Description Get a single region by its ID
// @Tags regions
// @Produce json
// @Param id path int true "Region ID"
// @Success 200 {object} models.Region
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /regions/{id} [get]
func (h *RegionHandler) GetRegion(c *gin.Context) {
    id := c.Param("id")
    
    regionID, err := validator.ValidateID(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   "Invalid ID",
            "message": err.Error(),
        })
        return
    }
    
    var region models.Region
    if err := h.DB.Preload("Areas").First(&region, regionID).Error; err != nil {
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
    
    c.JSON(http.StatusOK, region)
}

// GetRegions возвращает все регионы с пагинацией
// @Summary Get all regions
// @Description Get a list of all regions with pagination and filtering
// @Tags regions
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of items per page" default(20)
// @Param name query string false "Filter by name"
// @Param sort_by query string false "Sort field" default(id)
// @Param order query string false "Sort order" Enums(asc, desc) default(asc)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /regions [get]
func (h *RegionHandler) GetRegions(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
    
    if page < 1 {
        page = 1
    }
    if limit < 1 || limit > 100 {
        limit = 20
    }
    
    offset := (page - 1) * limit
    
    var regions []models.Region
    var total int64
    
    // Поиск по имени
    query := h.DB.Model(&models.Region{})
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
    query.Count(&total)
    result := query.Offset(offset).Limit(limit).Find(&regions)
    
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Database error",
            "details": result.Error.Error(),
        })
        return
    }
    
    response := gin.H{
        "data": regions,
        "pagination": gin.H{
            "page":       page,
            "limit":      limit,
            "total":      total,
            "totalPages": (int(total) + limit - 1) / limit,
        },
    }
    
    c.JSON(http.StatusOK, response)
}

// UpdateRegion обновляет регион
// @Summary Update region
// @Description Update an existing region by ID
// @Tags regions
// @Accept json
// @Produce json
// @Param id path int true "Region ID"
// @Param region body models.RegionUpdateRequest true "Region update object"
// @Success 200 {object} models.Region
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /regions/{id} [put]
func (h *RegionHandler) UpdateRegion(c *gin.Context) {
    id := c.Param("id")
    
    regionID, err := validator.ValidateID(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   "Invalid ID",
            "message": err.Error(),
        })
        return
    }
    
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
    
    var input models.RegionUpdateRequest
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
    
    result := h.DB.Model(&region).Updates(updates)
    if result.Error != nil {
        if validator.IsDuplicateKeyError(result.Error) {
            c.JSON(http.StatusConflict, gin.H{
                "error":   "Update conflict",
                "message": "A region with this name already exists",
            })
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Database error",
            "details": result.Error.Error(),
        })
        return
    }
    
    h.DB.First(&region, regionID)
    c.JSON(http.StatusOK, region)
}

// DeleteRegion удаляет регион
// @Summary Delete region
// @Description Delete a region by ID
// @Tags regions
// @Produce json
// @Param id path int true "Region ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /regions/{id} [delete]
func (h *RegionHandler) DeleteRegion(c *gin.Context) {
    id := c.Param("id")
    
    regionID, err := validator.ValidateID(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error":   "Invalid ID",
            "message": err.Error(),
        })
        return
    }
    
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
    
    // Проверяем, есть ли связанные области
    var areasCount int64
    h.DB.Model(&models.Area{}).Where("region_id = ?", regionID).Count(&areasCount)
    
    if areasCount > 0 {
        c.JSON(http.StatusConflict, gin.H{
            "error":   "Cannot delete region",
            "message": "Region has associated areas. Delete them first.",
            "areas_count": areasCount,
        })
        return
    }
    
    result := h.DB.Delete(&region)
    if result.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error":   "Database error",
            "details": result.Error.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "message": "Region deleted successfully",
        "id":      regionID,
    })
}