package api

import (
	"github.com/gin-gonic/gin"

	"intellectual_property/internal/service"
)

func initLogin(r *gin.Engine) {
	//登录接口
	r.POST("/login", service.Login)

}
