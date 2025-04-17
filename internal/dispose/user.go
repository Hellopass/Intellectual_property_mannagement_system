package dispose

import (
	"github.com/gin-gonic/gin"
	"intellectual_property/pkg/models"
	"intellectual_property/pkg/utils"
	"net/http"
	"strconv"
)

// 封装返回码
const (
	UserExists = iota + 600 //用户存在 1000
	UserNoExists
	SystemError //系统错误
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
	gender := c.PostForm("gender")
	//根据email查询用户是否存在
	_, err := models.GetUserByEmail(email)
	if err != nil {
		logger.Error(err.Error())
		Resp(c, false, UserExists, "用户存在", "")
		return
	}
	//首先需要验证code是否正确
	//如果不存在
	re_code, err := utils.RedisGet(email)
	if err != nil {
		logger.Error(err.Error())
		Resp(c, false, SystemError, "系统错误", "")
		return
	}
	if re_code != code {
		Resp(c, false, http.StatusBadRequest, "验证码错误", "")
		return
	}
	//验证完删除缓存
	err = utils.RedisDel(email)
	if err != nil {
		logger.Error(err.Error())
		Resp(c, false, SystemError, "系统错误", "")
		return
	}
	//拿到加密密码
	salt, err := utils.GenerateSecureSalt()
	if err != nil {
		logger.Error(err.Error())
		Resp(c, false, SystemError, "系统错误", "")
		return
	}
	EnPassword, err1 := utils.SecureHashWithSalt(password, salt)
	if err1 != nil {
		logger.Error(err1.Error())
		Resp(c, false, SystemError, "系统错误", "")
		return
	}
	//创建用户
	u := models.User{
		UserName: username,
		Email:    email,
		Password: salt + ":" + EnPassword,
		IDCard:   idcard,
		Sex:      gender,
		Status:   "1",
	}

	//添加用户
	err2 := models.CreateUser(&u)
	if err2 != nil {
		logger.Error(err2.Error())
		Resp(c, false, SystemError, "添加失败", "")
		return
	}

	//返回成功json-
	Resp(c, true, http.StatusOK, "注册成功", gin.H{})
}

// FindUsersByID 查询用户-id
func FindUsersByID(c *gin.Context) {
	id := c.Query("user_id")
	//查询id
	ids, err2 := strconv.Atoi(id)
	if err2 != nil {
		logger.Error(err2.Error())
		Resp(c, false, SystemError, "系统错误", "")
		return
	}
	user, err := models.GetUserByID(ids)
	if err != nil {
		logger.Error(err.Error())
		Resp(c, false, SystemError, "系统错误", "")
		return
	}
	Resp(c, true, http.StatusOK, "查询成功", user)
}

// ModifyTheUser 修改用户
func ModifyTheUser(c *gin.Context) {
	uid := c.PostForm("user_id")             //用户id
	last_degree := c.PostForm("last_degree") //最高学历
	political := c.PostForm("political")     //政治面貌
	cour := c.PostForm("cour")               //一级学科
	dep_id := c.PostForm("dep_id")           //部门id
	research := c.PostForm("research")       //研究方向
	unit := c.PostForm("unit")               //所属学院
	tech_ip := c.PostForm("tech_ip")         //技术职称
	avatarUrl := c.PostForm("avatar_url")
	user_id, err := strconv.Atoi(uid)
	if err != nil {
		logger.Error(err.Error())
		Resp(c, false, http.StatusBadRequest, "修改失败", gin.H{
			"error": err.Error(),
		})
		return
	}
	dep_ID, err1 := strconv.Atoi(dep_id)
	if err1 != nil {
		logger.Error(err1.Error())
		Resp(c, false, http.StatusBadRequest, "修改失败", gin.H{
			"error": err1.Error(),
		})
		return
	}

	us := models.User{
		ID:         user_id,
		DepID:      dep_ID,
		Political:  political,
		Unit:       unit,
		LastDegree: last_degree,
		TechIP:     tech_ip,
		Cour:       cour,
		Research:   research,
		AvatarUrl:  avatarUrl,
	}
	//更新用户
	err2 := models.UpdateUser(&us)
	if err2 != nil {
		logger.Error(err2.Error())
		Resp(c, false, http.StatusBadRequest, "修改失败", gin.H{
			"error": err2.Error(),
		})
		return
	}
	//修改成功
	Resp(c, true, http.StatusOK, "修改成功", "")
}

// UploadAvatar 上传头像
func UploadAvatar(c *gin.Context) {
	//实现头像文件上传并拿到地址
	file, err4 := c.FormFile("avatar")
	id := c.PostForm("user_id")
	if err4 != nil {
		logger.Error(err4.Error())
		Resp(c, false, http.StatusBadRequest, "修改失败", gin.H{
			"error": err4.Error(),
		})
		return
	}
	err := c.SaveUploadedFile(file, utils.NgX.LocationAvatar+"/"+file.Filename)
	if err != nil {
		logger.Error(err.Error())
		Resp(c, false, http.StatusBadRequest, "修改失败", gin.H{
			"error": err.Error(),
		})
		return
	}
	//拼接地址
	avatarUrl := utils.NgX.Url + "/" + "avatar" + "/" + file.Filename
	//上传数据库
	atoi, err4 := strconv.Atoi(id)
	if err4 != nil {
		logger.Error(err.Error())
		Resp(c, false, http.StatusBadRequest, "修改失败", gin.H{
			"error": err4.Error(),
		})
		return
	}
	err4 = models.UploadAvatar(atoi, avatarUrl)
	if err4 != nil {
		logger.Error(err.Error())
		Resp(c, false, http.StatusBadRequest, "修改失败", gin.H{
			"error": err4.Error(),
		})
		return
	}

	Resp(c, true, http.StatusOK, "上传成功", avatarUrl)
}
