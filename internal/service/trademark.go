package service

import (
	"intellectual_property/pkg/models"
	"intellectual_property/pkg/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateTrademark 创建商标

func CreateTrademark(c *gin.Context) {
	// 1. 身份验证
	authHeader := c.GetHeader("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	_, err := utils.ParseToken(tokenString)
	if err != nil {
		logger.Error("Token解析失败: " + err.Error())
		Resp(c, false, http.StatusUnauthorized, "无效的访问令牌", nil)
		return
	}

	// 2. 解析表单数据
	trademarkTypeStr := c.PostForm("trademarkType")
	title := c.PostForm("title")
	abstract := c.PostForm("abstract")
	firstAuthorID, _ := strconv.Atoi(c.PostForm("firstAuthorId"))
	authorIDs := utils.ConvertStringSliceToInt(c.PostFormArray("authors"))

	// 3. 参数验证
	if title == "" || abstract == "" || len(authorIDs) == 0 {
		Resp(c, false, http.StatusBadRequest, "必要参数不能为空", nil)
		return
	}

	// 3. 创建商标记录
	trademarkType, err := models.ParseTrademarkType(trademarkTypeStr)
	if err != nil {
		Resp(c, false, http.StatusBadRequest, "无效的商标类型", nil)
		return
	}
	// 生成申请号
	number, err := utils.GenerateTrademarkApplicationNumber("CN", strconv.Itoa(time.Now().Year()), trademarkType)
	if err != nil {
		utils.Logger.Error(err.Error())
		Resp(c, false, http.StatusBadRequest, err.Error(), nil)
		return
	}
	trademark, trademarkfee, err := models.NewTrademark(
		trademarkType,
		title,
		authorIDs,
		firstAuthorID,
		abstract,
	)
	trademark.ApplicationNumber = number
	trademark.ApplyDate = time.Now()
	if err != nil {
		logger.Error("参数验证失败: " + err.Error())
		Resp(c, false, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := trademark.CreateTrademarkService(trademarkfee); err != nil {
		logger.Error("数据库操作失败: " + err.Error())
		Resp(c, false, http.StatusInternalServerError, "创建商标失败", nil)
		return
	}
	//4，文件上传
	form, err2 := c.MultipartForm()
	if err2 != nil {
		logger.Error(err2.Error())
		Resp(c, false, http.StatusInternalServerError, "上传相关文件失败", nil)
		return
	}
	Files := form.File["files"] //拿到文件
	url := utils.NgX.LocationTrademark + "/" + trademark.ApplicationNumber
	for _, f := range Files {
		//保存文件 地址---原始地址/申请号文件/文件

		err = c.SaveUploadedFile(f, url+"/"+f.Filename)
		if err != nil {
			logger.Error(err.Error())
			Resp(c, false, http.StatusInternalServerError, "上传相关文件失败", nil)
			return
		}
	}
	//在更新一下数据库--只跟新文件位置
	trademark.AttachmentUrl = url
	err3 := trademark.UpdateTrademarkUrl()
	if err3 != nil {
		logger.Error(err3.Error())
		Resp(c, false, http.StatusInternalServerError, "上传相关文件失败", nil)
		return
	}
	// 5. 返回成功响应
	Resp(c, true, http.StatusCreated, "商标创建成功", gin.H{
		"trademarkId": trademark.ID,
		"title":       trademark.Title,
		"firstAuthor": trademark.FirstAuthorID,
		"fileCount":   len(Files),
	})
}

// GetAllTrademarks 获取所有商标
func GetAllTrademarks(c *gin.Context) {
	page := c.Query("page")
	pageSize := c.Query("pageSize")
	keyword := c.Query("search")
	status := c.Query("status")
	types := c.Query("type")
	var (
		a int //状态
		b int
	)

	if status == "" {
		a = -1
	} else {
		a, _ = strconv.Atoi(status)
	}
	if types == "" {
		b = 0
	} else {
		b, _ = strconv.Atoi(types)
	}
	page_a, _ := strconv.Atoi(page)
	page_size_a, _ := strconv.Atoi(pageSize)

	trademarks, total, err := models.GetAllTrademarks(keyword, a, b, page_a, page_size_a, true)

	if err != nil {
		logger.Error(err.Error())
		Resp(c, false, http.StatusInternalServerError, "获取失败", nil)
		return
	}

	Resp(c, true, http.StatusOK, "获取成功", gin.H{
		"trademarks": trademarks,
		"total":      total,
	})
}

// GetTrademarkFile 获取文件地址
func GetTrademarkFile(c *gin.Context) {
	value := c.Query("id")
	atoi, _ := strconv.Atoi(value)
	file, err := models.GetTrademarkFile(atoi)
	if err != nil {
		logger.Error(err.Error())
		Resp(c, false, http.StatusInternalServerError, "获取失败", nil)
		return
	}
	Resp(c, true, http.StatusOK, "查询成功", file)
}

// DeleteTrademark 删除商标
func DeleteTrademark(c *gin.Context) {
	id := c.Query("id")

	atoi, _ := strconv.Atoi(id)
	if err := models.DeleteTrademark(atoi); err != nil {
		logger.Error(err.Error())
		Resp(c, false, http.StatusInternalServerError, "删除失败", nil)
		return
	}
	Resp(c, true, http.StatusOK, "删除成功", nil)
}

// UpdateTrademarkStatus 更新审核状态
func UpdateTrademarkStatus(c *gin.Context) {
	value := c.PostForm("trademark_id")
	form := c.PostForm("reviewer_id")
	comment := c.PostForm("comment")
	v := c.PostForm("status")
	trademark_id, _ := strconv.Atoi(value)
	reviewer_id, _ := strconv.Atoi(form)
	status, _ := strconv.Atoi(v)
	//审核没通过
	if err := models.UpdateTrademarkStatus(trademark_id, reviewer_id, comment, status); err != nil {
		logger.Error(err.Error())
		Resp(c, false, http.StatusInternalServerError, "更新失败", nil)
		return
	}
	Resp(c, true, http.StatusOK, "更新成功", nil)
}

// GetAllTrademarkFees 获取所有商标年费服务方法
func GetAllTrademarkFees(c *gin.Context) {
	page := c.Query("page")
	pageSize := c.Query("pageSize")
	keyword := c.Query("keyword")
	statusStr := c.Query("status")

	var (
		status      int
		pageInt     int
		pageSizeInt int
		err         error
	)

	// 解析状态
	if statusStr == "" {
		status = -1
	} else {
		status, err = strconv.Atoi(statusStr)
		if err != nil {
			utils.Logger.Error("状态参数解析失败: " + err.Error())
			Resp(c, false, http.StatusBadRequest, "无效的状态参数", nil)
			return
		}
	}

	// 解析分页参数
	pageInt, err = strconv.Atoi(page)
	if err != nil {
		pageInt = 1
	}
	pageSizeInt, err = strconv.Atoi(pageSize)
	if err != nil {
		pageSizeInt = 10
	}

	fees, total, err := models.GetAllTrademarkFees(keyword, status, pageInt, pageSizeInt)
	if err != nil {
		utils.Logger.Error(err.Error())
		Resp(c, false, http.StatusInternalServerError, "获取商标年费失败", nil)
		return
	}

	Resp(c, true, http.StatusOK, "获取商标年费成功", gin.H{
		"fees":  fees,
		"total": total,
	})
}

// GetMonthlyTrademarkFeeStatsService 获取本月商标年费统计信息服务方法
func GetMonthlyTrademarkFeeStatsService(c *gin.Context) {
	pendingCount, paidAmount, overdueCount, totalAmount, err := models.GetMonthlyTrademarkFeeStats()
	if err != nil {
		utils.Logger.Error(err.Error())
		Resp(c, false, http.StatusInternalServerError, "获取本月商标年费统计信息失败", nil)
		return
	}

	Resp(c, true, http.StatusOK, "获取本月商标年费统计信息成功", gin.H{
		"pending_count": pendingCount,
		"paid_amount":   paidAmount,
		"overdue_count": overdueCount,
		"total_amount":  totalAmount,
	})
}
