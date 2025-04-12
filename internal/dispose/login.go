package dispose

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"intellectual_property/pkg/models"
	"intellectual_property/pkg/utils"
)

//	r.POST("/login", func(c *gin.Context) {
//		// 实际应验证用户名密码
//		token, err := utils.GenerateToken("123")
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "无法生成令牌"})
//			return
//		}
//		c.JSON(http.StatusOK, gin.H{"token": token})
//	})
//

// Login 登录
func Login(c *gin.Context) {
	email := c.PostForm("username")    //这里的username 表示的是有邮件地址
	password := c.PostForm("password") //密码
	//eamil作为查询来查询User信息
	user, err := models.GetUserByEmail(email)
	if err != nil {
		utils.Logger.Error(err.Error())
	}
	fmt.Println(user)
	fmt.Println(password)
}
