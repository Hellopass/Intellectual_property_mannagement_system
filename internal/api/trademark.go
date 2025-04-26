package api

import (
	"intellectual_property/internal/service"

	"github.com/gin-gonic/gin"
)

func initTrademark(r *gin.Engine) {
	group := r.Group("/trademark")

	// 新增商标申请
	group.POST("/add", service.CreateTrademark)

	// 获取商标信息并模糊查询
	group.GET("/get_trademarks", service.GetAllTrademarks)
	// 查询所有的商标文件
	group.GET("/get_files", service.GetTrademarkFile)

	// 删除商标
	group.DELETE("/del_trademark", service.DeleteTrademark)

	// 更新审核状态
	group.PUT("/update_trademark_status", service.UpdateTrademarkStatus)

	// 获取所有商标年费
	group.GET("/get_fee_all", service.GetAllTrademarkFees)

	// 获取本月商标年费统计信息
	group.GET("/get_monthly_fee_stats", service.GetMonthlyTrademarkFeeStatsService)
}
