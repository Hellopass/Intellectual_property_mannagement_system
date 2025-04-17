package api

import "github.com/gin-gonic/gin"

func InitApi(r *gin.Engine) {
	initUser(r)
	initLogin(r)
	initEmail(r)
	initPatent(r)
}
