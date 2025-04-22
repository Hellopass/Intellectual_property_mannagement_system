package api

import (
	"github.com/gin-gonic/gin"
	"intellectual_property/internal/service"
)

func initPatent(r *gin.Engine) {
	group := r.Group("/patent")
	//group.Use(utils.JWTMiddleware())

	//新建申请
	group.POST("/add", service.AddPatent)

	//上传文件
	group.POST("/upload_file", service.UploadPatentFile)

	//获取所有专利信息
	group.GET("/find", service.FindPatentS)

	//模糊查询
	group.GET("/find_fuzzy", service.FindPatentFuzzy)

	//根据申请号拿到所有的文件
	group.GET("/find_file", service.GetPatentFile)

	//删除专利信息
	group.DELETE("/delete", service.DelPatent)

	//跟新状态
	group.PUT("/update_status", service.UpdateStatus)

	//专利年费
	group.GET("/get_fee_statistics", service.GetFeeStatistics)

	//获取年费信息
	group.GET("/get_fee_all", service.GetAllPatentFees)

	//更新金额
	group.PUT("/update_amount", service.UpdatePatentFeeByApplyNo)

	//年费模糊查询
	group.GET("/get_fee_fuzzy", service.GetPatentFeesByFilters)

	//专利分析
	group.GET("/get_analysis", service.GetAnalysis)
}
