package main

import (
	"intellectual_property/internal/api"
	"intellectual_property/pkg/models"
)

func main() {
	// 创建预配置的引擎
	r := api.NewEngine("debug") // 参数："debug" 或 "release"
	models.CreatePatentTable()
	// 添加路由
	api.InitApi(r)

	models.AutoMigrate()
	// 启动服务
	r.Run(":8080")
	//rand.Seed(time.Now().UnixNano())
	//a := map[int]string{
	//	0: "左边",
	//	1: "右边",
	//	2: "中间"}
	//fmt.Println(a[rand.Intn(3)])
}
