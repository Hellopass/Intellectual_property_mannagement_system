package api

import (
	"github.com/gin-gonic/gin"
	"intellectual_property/internal/service"
)

func initUser(r *gin.Engine) {
	group := r.Group("/user")

	//查询所有用户
	group.GET("/find_all", service.GetSimpleUsers)

	//查询用户-id
	group.GET("/find", service.FindUsersByID)

	//上传头像
	group.POST("/upload_avatar", service.UploadAvatar)
	//增加用户
	group.POST("/add", service.AddUser)

	//修改用户
	group.PUT("/edit", service.ModifyTheUser)

}
