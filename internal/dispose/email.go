package dispose

import (
	"github.com/gin-gonic/gin"
	"intellectual_property/pkg/utils"
	"net/http"
)

func SendEmail(c *gin.Context) {
	value := c.PostForm("email")
	err := utils.SendAddUserEmailCode(value)
	if err != nil {
		logger.Error(err.Error())
		Resp(c, false, SystemError, "验证码发送失败", "")
	}
	Resp(c, true, http.StatusOK, "验证码发送成功", "")
}
