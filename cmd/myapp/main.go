package main

import (
	"intellectual_property/internal/api"
)

func main() {
	// 创建预配置的引擎
	r := api.NewEngine("debug") // 参数："debug" 或 "release"

	// 添加路由
	api.InitApi(r)
	//r.GET("/", func(c *gin.Context) {
	//	c.JSON(200, gin.H{"message": "Hello World"})
	//})

	// 启动服务
	r.Run(":8080")
}

//func main() {
//	r := gin.Default()
//
//	// 登录接口（生成 Token）
//	r.POST("/login", func(c *gin.Context) {
//		// 实际应验证用户名密码
//		token, err := utils.GenerateToken("123")
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, gin.H{"error": "无法生成令牌"})
//			return
//		}
//		c.JSON(http.StatusOK, gin.H{"token": token})
//	})
//	// 刷新令牌机制
//	r.POST("/refresh", utils.JWTMiddleware(), func(c *gin.Context) {
//		userID := c.MustGet("userID").(string)
//		newToken, _ := utils.GenerateToken(userID)
//		c.JSON(http.StatusOK, gin.H{"token": newToken})
//	})
//
//	// 需要认证的路由
//	authGroup := r.Group("/api")
//	authGroup.Use(utils.JWTMiddleware())
//	{
//		authGroup.GET("/profile", func(c *gin.Context) {
//			userID := c.MustGet("userID").(string)
//			c.JSON(http.StatusOK, gin.H{"message": "欢迎 " + userID})
//		})
//	}
//
//	r.Run(":8080")
//}
