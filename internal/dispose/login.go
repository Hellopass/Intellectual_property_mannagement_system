package dispose

import (
	"github.com/gin-gonic/gin"
	"intellectual_property/pkg/models"
	"intellectual_property/pkg/utils"
	"net/http"
)

// Login 登录
func Login(c *gin.Context) {
	email := c.PostForm("username")    //这里的username 表示的是有邮件地址
	password := c.PostForm("password") //密码

	//eamil作为查询来查询User信息
	user, err := models.GetUserByEmail(email)
	if err != nil {
		utils.Logger.Error(err.Error())
		Resp(c, false, UserNoExists, "用户不存在", "")
		return
	}
	//表示存在
	if user.ID != 0 {
		//验证密码是否正确
		//拿到研盐值
		slat := utils.GetSlat(user.Password)
		if utils.VerifyPassword(password, slat, user.Password) {
			//密码正确
			//生成token
			token, err3 := utils.GenerateToken(user.ID, user.Authority, user.UserName)
			if err3 != nil {
				logger.Error(err3.Error())
				return
			}
			Resp(c, true, http.StatusOK, "登录成功", gin.H{
				"token":   token,
				"success": true,
			})
			return
		} else {
			//密码错误
			Resp(c, false, http.StatusOK, "密码错误", "")
			return
		}
	}
	Resp(c, false, http.StatusOK, "用户不存在", "")
}
