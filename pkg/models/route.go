package models

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"intellectual_property/pkg/utils"
	"time"
)

// RouteStats 存储路由的统计信息
type RouteStats struct {
	ID        int    `json:"id" gorm:"AUTO_INCREMENT;primary_key"`             //主键
	Method    string `json:"method" gorm:"column:method;type:varchar(20)"`     // 请求方法
	Path      string `json:"path" gorm:"column:path;type:varchar(255)"`        // 请求地址
	Handler   string `json:"handler" gorm:"column:handler;type:varchar(255)"`  //所在文件地址
	Status    bool   `json:"status" gorm:"column:status;type:bool"`            //是否启用
	Count     int64  `json:"count"  gorm:"column:count;type:int"`              //调用次数
	TotalTime int64  `json:"total_time" gorm:"column:total_time;type:bigint"`  //总的调用时间
	Average   int64  `json:"average_ms"  gorm:"column:average_ms;type:bigint"` //平均耗时
}

// MiddlewareRoute Gin中间件，用于统计请求信息
func MiddlewareRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录开始时间
		start := time.Now()

		// 处理请求
		c.Next()

		// 计算耗时
		duration := time.Since(start)

		// 获取路由信息
		method := c.Request.Method
		path := c.FullPath()
		if path == "" {
			path = "404" // 处理未匹配路由
		}

		//保存到数据库
		//首先根据path查询
		var collect *RouteStats
		if err := utils.DB.Model(&RouteStats{}).Where("path = ?", path).Find(&collect).Error; err != nil {
			utils.Logger.Error("数据库错误")
			return
		}
		//不存在
		if collect.ID == 0 {
			collect.Path = path
			collect.Method = method
			collect.Count++
			collect.TotalTime = int64(duration)
			collect.Average = collect.TotalTime / collect.Count
			collect.Status = true
			//	handle在前面更新
			if err := utils.DB.Debug().Model(collect).Create(&collect).Error; err != nil {
				utils.Logger.Error("数据库错误")
				return
			}
		} else {
			//更新
			collect.Count++
			collect.TotalTime += int64(duration)
			collect.Average = collect.TotalTime / collect.Count
			if err := utils.DB.Debug().Model(collect).Updates(&collect).Error; err != nil {
				utils.Logger.Error("数据库错误")
				return
			}
		}

	}

}

// InterceptRoute 这里是实现对特点路由的拦截
// 对不符合的进行拦截
func InterceptRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		//首先查询所有的路由，查看status状态
		path := c.FullPath()
		var r RouteStats
		if err := utils.DB.Debug().Where("path=?", path).Find(&r).Error; err != nil {
			utils.Logger.Error(err.Error())
			return
		}
		fmt.Println(r.Status)
		if r.Status {
			c.Next()
		}
		c.Abort()
	}
}

func RouteStatusChange(status bool, path string) error {
	return utils.DB.Debug().Model(&RouteStats{}).Where("path=?", path).Update("status", status).Error
}
