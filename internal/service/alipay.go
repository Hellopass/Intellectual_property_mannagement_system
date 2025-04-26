package service

import "github.com/gin-gonic/gin"

// AlipayToSuccess 回调函数
func AlipayToSuccess(c *gin.Context) {
	c.JSON(200, "支付成功")
}

// AlipayToPatent 专利年费支付
func AlipayToPatent(c *gin.Context) {

}
