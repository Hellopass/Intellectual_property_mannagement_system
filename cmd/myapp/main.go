package main

import (
	"intellectual_property/internal/api"
	"intellectual_property/pkg/models"
	"intellectual_property/pkg/utils"
)

func main() {
	// 创建预配置的引擎
	r := api.NewEngine("debug") // 参数："debug" 或 "release"
	utils.DB.AutoMigrate(
		&models.User{},
		&models.Patent{},
		&models.PatentAuthor{},
		&models.PatentFee{},
		&models.TrademarkAuthor{},
		&models.Trademark{},
		&models.TrademarkFee{},
		&models.Article{},
		&models.ArticleAuthor{},
		&models.ArticleFee{},
		&models.RouteStats{},
	)

	// 添加路由
	api.InitApi(r)

	// 启动服务
	r.Run(":8080")
}
