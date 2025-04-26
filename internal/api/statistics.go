package api

import (
	"github.com/gin-gonic/gin"
	"intellectual_property/internal/service"
)

func initStatistics(r *gin.Engine) {
	group := r.Group("/statistics")

	group.GET("/get_all", service.GetAllStatistic)
}
