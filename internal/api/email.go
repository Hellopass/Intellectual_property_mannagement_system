package api

import (
	"github.com/gin-gonic/gin"
	"intellectual_property/internal/dispose"
)

func initEmail(r *gin.Engine) {
	//发送邮件验证码
	r.POST("/email", dispose.SendEmail)
}
