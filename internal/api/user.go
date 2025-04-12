package api

import (
	"github.com/gin-gonic/gin"
	"intellectual_property/internal/dispose"
)

func initUser(r *gin.Engine) {
	group := r.Group("/user")
	//增加用户
	group.POST("/add", dispose.AddUser)
}
