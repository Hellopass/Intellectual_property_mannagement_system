package service

import (
	"intellectual_property/pkg/models"
	"intellectual_property/pkg/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// StatsHandler 返回统计访问信息
func StatsHandler(c *gin.Context) {
	// 获取分页参数
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")
	keyword := c.DefaultQuery("keyword", "")
	statusStr := c.DefaultQuery("status", "")
	// 从前端获取分割好的 path 部分
	splitPath := c.DefaultQuery("split_path", "")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 10
	}

	// 解析状态
	var status *bool
	if statusStr != "" {
		boolValue, err := strconv.ParseBool(statusStr)
		if err == nil {
			status = &boolValue
		}
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 构建查询
	query := utils.DB.Model(&models.RouteStats{})

	// 根据 keyword 进行模糊查询
	if keyword != "" {
		likeKeyword := "%" + keyword + "%"
		query = query.Where("path LIKE ? OR handler LIKE ?", likeKeyword, likeKeyword)
	}

	// 根据 status 进行查询
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// 根据前端传递的分割好的 path 进行查询
	if splitPath != "" {
		likeSplitPath := "%" + splitPath + "%"
		query = query.Where("path LIKE ?", likeSplitPath)
	}

	// 计算总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		utils.Logger.Error(err.Error())
		Resp(c, false, http.StatusInternalServerError, "获取接口信息失败", nil)
		return
	}

	// 添加排序逻辑，按 path 字段按 '/' 分割后取索引为 1 的部分排序
	var routes []models.RouteStats
	if err := query.
		Offset(offset).
		Limit(pageSize).
		Order("SUBSTRING_INDEX(SUBSTRING_INDEX(path, '/', 2), '/', -1)").
		Find(&routes).Error; err != nil {
		utils.Logger.Error(err.Error())
		Resp(c, false, http.StatusInternalServerError, "获取接口信息失败", nil)
		return
	}

	// 返回成功响应
	Resp(c, true, http.StatusOK, "获取接口信息成功", gin.H{
		"routes":    routes,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func RouteStatusChange(c *gin.Context) {
	value := c.PostForm("status")
	path := c.PostForm("path")
	status, _ := strconv.ParseBool(value)
	if err := models.RouteStatusChange(status, path); err != nil {
		utils.Logger.Error(err.Error())
		Resp(c, false, http.StatusInternalServerError, "更新失败", nil)
	}
	Resp(c, true, http.StatusOK, "更新成功", nil)
}
