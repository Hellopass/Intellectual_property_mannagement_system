package api

import (
	"github.com/gin-gonic/gin"
	"intellectual_property/pkg/models"
	"intellectual_property/pkg/utils"
)

var Resp = utils.Resp

func InitApi(r *gin.Engine) {
	r.Use(models.MiddlewareRoute()) //使用中间件来实现对路由信息的统计
	r.Use(models.InterceptRoute())
	initUser(r)
	initLogin(r)
	initStatistics(r)
	initEmail(r)
	initPatent(r)
	initArticle(r)
	// 新增初始化商标路由
	initTrademark(r)
	initRoute(r)
	initAlipay(r) //支付宝支付
	//拿到所有信息 --支持分页查询
	routes := r.Routes()
	for _, v := range routes {
		//先根据path查询一下是否存在
		var ro *models.RouteStats
		if err := utils.DB.Debug().Model(&models.RouteStats{}).Where("path = ?", v.Path).Find(&ro).Error; err != nil {
			utils.Logger.Error(err.Error())
			return
		}

		//不存在
		if ro.ID == 0 {
			ro.Path = v.Path
			ro.Method = v.Method
			ro.Handler = v.Handler
			ro.Status = true
			if err := utils.DB.Model(ro).Create(ro).Error; err != nil {
				utils.Logger.Error("数据库错误")
				return
			}
		} else {
			ro.Status = true
			ro.Handler = v.Handler
			if err := utils.DB.Model(ro).Updates(ro).Error; err != nil {
				utils.Logger.Error("数据库错误")
				return
			}
		}

	}

}
