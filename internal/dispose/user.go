package dispose

import (
	"github.com/gin-gonic/gin"
	"intellectual_property/pkg/models"
	"intellectual_property/pkg/utils"
	"net/http"
)

const (
	UserExists  = iota + 1000 //用户存在 1000
	SystemError               //系统错误
	CodeError
)

// 日志
var logger = utils.Logger

// Resp 返回体
var Resp = utils.Resp

// AddUser 注册用户
func AddUser(c *gin.Context) {
	username := c.PostForm("name")
	email := c.PostForm("email")
	code := c.PostForm("verificationCode") //需要从redis中读取验证码
	password := c.PostForm("password")
	idcard := c.PostForm("idCard")

	//根据email查询用户是否存在
	_, err := models.GetUserByEmail(email)
	if err != nil {
		logger.Error(err.Error())
		Resp(c, false, UserExists, "用户存在", "")
	}
	//首先需要验证code是否正确
	//如果不存在
	re_code, err := utils.RedisGet(email)
	if err != nil {
		logger.Error(err.Error())
		Resp(c, false, SystemError, "系统错误", "")
	}
	if re_code != code {
		Resp(c, false, CodeError, "验证码错误", "")
	}
	//验证完删除缓存
	err = utils.RedisDel(email)
	if err != nil {
		logger.Error(err.Error())
	}
	//拿到加密密码
	salt, err := utils.GenerateSecureSalt()
	if err != nil {
		logger.Error(err.Error())
		Resp(c, false, SystemError, "系统错误", "")
	}
	EnPassword, err1 := utils.SecureHashWithSalt(password, salt)
	if err1 != nil {
		logger.Error(err1.Error())
		Resp(c, false, SystemError, "系统错误", "")
	}
	//创建用户
	u := models.User{
		UserName: username,
		Email:    email,
		Password: EnPassword,
		IDCard:   idcard,
	}

	//添加用户
	err2 := models.CreateUser(&u)
	if err2 != nil {
		logger.Error(err2.Error())
		Resp(c, false, SystemError, "添加失败", "")
	}
	//拿到用户id
	//根据email拿到信息
	us, err4 := models.GetUserByEmail(email)
	if err4 != nil {
		logger.Error(err4.Error())
		Resp(c, false, SystemError, "添加失败", "")
	}

	//生成jwt
	token, err3 := utils.GenerateToken(us.ID, us.Authority, us.UserName)
	if err3 != nil {
		logger.Error(err3.Error())
		Resp(c, false, http.StatusInternalServerError, "生成token错误", "")
	}

	//返回成功json--返回jwt
	Resp(c, true, http.StatusOK, "注册成功", gin.H{
		"token": token,
	})
}
