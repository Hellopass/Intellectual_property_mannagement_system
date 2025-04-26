package service

import (
	"github.com/gin-gonic/gin"
	"intellectual_property/pkg/models"
	"net/http"
)

func GetAllStatistic(c *gin.Context) {
	statistics, err := models.GetStatistics()
	if err != nil {
		logger.Error(err.Error())
		Resp(c, false, http.StatusBadRequest, "获取失败", nil)
		return
	}
	Resp(c, true, http.StatusOK, "获取成功", statistics)
}
