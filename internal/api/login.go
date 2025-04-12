package api

import (
	"github.com/gin-gonic/gin"
	"intellectual_property/internal/dispose"
)

func initLogin(r *gin.Engine) {
	//登录接口
	r.POST("/login", dispose.Login)

}
