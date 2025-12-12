package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"meteodata/internal/models"
)

// StatisticsHandler обработчик для статистики
type StatisticsHandler struct {
	*BaseHandler
}

// NewStatisticsHandler создает новый обработчик для статистики
func NewStatisticsHandler(db *gorm.DB) *StatisticsHandler {
	return &StatisticsHandler{
		BaseHandler: NewBaseHandler(db),
	}
}

// @Summary Статистика по измерениям
// @Description Возвращает статистику по всем типам измерений
// @Tags statistics
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /statistics [get]
func (h *StatisticsHandler) GetStatistics(c *gin.Context) {
	stats := gin.H{
		"regular": gin.H{},
		"led":     gin.H{},
		"tttr":    gin.H{},
		"snow":    gin.H{},
	}

	var regularCount, ledCount, tttrCount, snowCount int64

	h.DB.Model(&models.Regular{}).Count(&regularCount)
	h.DB.Model(&models.LED{}).Count(&ledCount)
	h.DB.Model(&models.TTTR{}).Count(&tttrCount)
	h.DB.Model(&models.SNOW{}).Count(&snowCount)

	stats["regular"] = gin.H{"count": regularCount}
	stats["led"] = gin.H{"count": ledCount}
	stats["tttr"] = gin.H{"count": tttrCount}
	stats["snow"] = gin.H{"count": snowCount}

	c.JSON(http.StatusOK, gin.H{
		"statistics": stats,
		"timestamp":  time.Now().Unix(),
	})
}
