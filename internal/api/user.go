package api

import (
	"github.com/gin-gonic/gin"
	"intellectual_property/internal/dispose"
)

func initUser(r *gin.Engine) {
	group := r.Group("/user")

	//查询用户-id
	group.GET("/find", dispose.FindUsersByID)

	//上传头像
	group.POST("/upload_avatar", dispose.UploadAvatar)
	//增加用户
	group.POST("/add", dispose.AddUser)

	//修改用户
	group.PUT("/edit", dispose.ModifyTheUser)

}
