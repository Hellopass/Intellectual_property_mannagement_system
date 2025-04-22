package api

import (
	"github.com/gin-gonic/gin"
	"intellectual_property/internal/service"
)

func initEmail(r *gin.Engine) {
	//发送邮件验证码
	r.POST("/email", service.SendEmail)
}
