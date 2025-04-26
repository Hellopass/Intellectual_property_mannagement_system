package api

import (
	"github.com/gin-gonic/gin"
	"intellectual_property/internal/service"
)

func initRoute(r *gin.Engine) {
	group := r.Group("/route")

	//接口访问信息
	group.GET("/interface_info", service.StatsHandler)

	//接口状态改变
	group.PUT("/interface_status", service.RouteStatusChange)
}
