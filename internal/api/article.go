package api

import (
	"intellectual_property/internal/service"

	"github.com/gin-gonic/gin"
)

func initArticle(r *gin.Engine) {
	group := r.Group("/article")

	//新增申请
	group.POST("/add", service.CreateArticle)

	//获取信息并模糊查询
	group.GET("/get_articles", service.GetAllArticles)
	//查询所有的文件
	group.GET("/get_files", service.GetArticleFile)

	//删除文件
	group.DELETE("/del_article", service.DeleteArticle)

	//跟新审核状态
	group.PUT("/update_article_aduit", service.UpdateArticleStatus)

	// 获取所有著作年费
	group.GET("/get_fee_all", service.GetAllArticleFees)

	// 获取本月著作年费统计信息
	group.GET("/get_monthly_fee_stats", service.GetMonthlyArticleFeeStatsService)
}
