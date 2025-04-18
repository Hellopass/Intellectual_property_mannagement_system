package api

import (
	"github.com/gin-gonic/gin"
	"intellectual_property/internal/dispose"
)

func initPatent(r *gin.Engine) {
	group := r.Group("/patent")
	//group.Use(utils.JWTMiddleware())

	//新建申请
	group.POST("/add", dispose.AddPatent)

	//上传文件
	group.POST("/upload_file", dispose.UploadPatentFile)

	//获取所有专利信息
	group.GET("/find", dispose.FindPatentS)

	//模糊查询
	group.GET("/find_fuzzy", dispose.FindPatentFuzzy)

	//根据申请号拿到所有的文件
	group.GET("/find_file", dispose.GetPatentFile)

	//删除专利信息
	group.DELETE("/delete", dispose.DelPatent)

	//跟新状态
	group.PUT("/update_status", dispose.UpdateStatus)

	//专利年费
	group.GET("/get_fee_statistics", dispose.GetFeeStatistics)

	//获取年费信息
	group.GET("/get_fee_all", dispose.GetAllPatentFees)

	//更新金额
	group.PUT("/update_amount", dispose.UpdatePatentFeeByApplyNo)

	//年费模糊查询
	group.GET("/get_fee_fuzzy", dispose.GetPatentFeesByFilters)
}
