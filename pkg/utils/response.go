package utils

import "github.com/gin-gonic/gin"

//数据格式为
/*
 {
"success": true,
"code": 200,
"message": "成功",
"data": {}
}
*/

// /---------------------------/

// Resp 封装返回格式
func Resp(c *gin.Context, success bool, code int, message string, data any) {
	c.JSON(code, gin.H{
		"success": success,
		"message": message,
		"data":    data,
	})
}
