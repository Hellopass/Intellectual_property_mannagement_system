package api

import (
	"github.com/gin-gonic/gin"
	"intellectual_property/internal/service"
)

func initAlipay(r *gin.Engine) {
	g := r.Group("/pay")

	//回调函数
	g.GET("/tosuccess", service.AlipayToSuccess)
}
