package api

import (
	"intellectual_property/internal/service"

	"github.com/gin-gonic/gin"
)

func initPatent(r *gin.Engine) {
	group := r.Group("/patent")

	// 新建申请
	group.POST("/add", service.CreatePatent)

	// 获取信息并模糊查询
	group.GET("/get_patents", service.GetAllPatents)
	// 查询所有的文件
	group.GET("/get_files", service.GetPatentFile)

	// 删除专利信息
	group.DELETE("/del_patent", service.DeletePatent)

	// 更新审核状态
	group.PUT("/update_patent_status", service.UpdatePatentStatus)

	// 获取所有专利年费
	group.GET("/get_fee_all", service.GetAllPatentFees)

	// 获取本月专利年费统计信息
	group.GET("/get_monthly_fee_stats", service.GetMonthlyPatentFeeStatsService)
}
