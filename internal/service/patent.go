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

// CreatePatent 创建专利
func CreatePatent(c *gin.Context) {
	// 1. 身份验证
	authHeader := c.GetHeader("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	_, err := utils.ParseToken(tokenString)
	if err != nil {
		utils.Logger.Error("Token解析失败: " + err.Error())
		Resp(c, false, http.StatusUnauthorized, "无效的访问令牌", nil)
		return
	}

	// 2. 解析表单数据
	patentTypeStr := c.PostForm("patentType")
	patentType, err := strconv.Atoi(patentTypeStr)
	if err != nil {
		utils.Logger.Error("专利类型解析失败: " + err.Error())
		Resp(c, false, http.StatusBadRequest, "无效的专利类型", nil)
		return
	}
	title := c.PostForm("title")
	abstract := c.PostForm("abstract")
	firstAuthorID, _ := strconv.Atoi(c.PostForm("firstAuthorId"))
	authorIDs := utils.ConvertStringSliceToInt(c.PostFormArray("authors"))

	// 3. 参数验证
	if title == "" || abstract == "" || len(authorIDs) == 0 {
		Resp(c, false, http.StatusBadRequest, "必要参数不能为空", nil)
		return
	}
	//生成申请号
	// 将 time.Now().Year() 的 int 类型转换为 string 类型
	number, err0 := utils.GenerateApplicationNumber("CN", strconv.Itoa(time.Now().Year()), patentType)
	if err0 != nil {
		utils.Logger.Error(err0.Error())
		Resp(c, false, http.StatusBadRequest, err0.Error(), nil)
		return
	}
	// 3. 创建专利记录
	patent, patentFee, err := models.NewPatent(
		patentType,
		title,
		authorIDs,
		firstAuthorID,
		abstract,
		number,
	)
	patent.ApplyDate = time.Now()

	if err != nil {
		utils.Logger.Error("参数验证失败: " + err.Error())
		Resp(c, false, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := patent.CreatePatentService(patentFee); err != nil {
		utils.Logger.Error("数据库操作失败: " + err.Error())
		Resp(c, false, http.StatusInternalServerError, "创建专利失败", nil)
		return
	}
	//4，文件上传
	form, err2 := c.MultipartForm()
	if err2 != nil {
		utils.Logger.Error(err2.Error())
		Resp(c, false, http.StatusInternalServerError, "上传相关文件失败", nil)
		return
	}
	Files := form.File["files"] //拿到文件
	url := utils.NgX.LocationPatent + "/" + patent.ApplicationNumber
	for _, f := range Files {
		//保存文件 地址---原始地址/申请号文件/文件

		err = c.SaveUploadedFile(f, url+"/"+f.Filename)
		if err != nil {
			utils.Logger.Error(err.Error())
			Resp(c, false, http.StatusInternalServerError, "上传相关文件失败", nil)
			return
		}
	}
	//在更新一下数据库--只跟新文件位置
	patent.AttachmentUrl = url
	err3 := patent.UpdatePatentUrl()
	if err3 != nil {
		utils.Logger.Error(err3.Error())
		Resp(c, false, http.StatusInternalServerError, "上传相关文件失败", nil)
		return
	}
	// 5. 返回成功响应
	Resp(c, true, http.StatusCreated, "专利创建成功", gin.H{
		"patentId":    patent.ID,
		"title":       patent.Title,
		"firstAuthor": patent.FirstAuthorID,
		"fileCount":   len(Files),
	})
}

// GetAllPatents 获取所有专利
func GetAllPatents(c *gin.Context) {
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

	patents, total, err := models.GetAllPatents(keyword, a, b, page_a, page_size_a, true)

	if err != nil {
		utils.Logger.Error(err.Error())
		Resp(c, false, http.StatusInternalServerError, "获取失败", nil)
		return
	}

	Resp(c, true, http.StatusOK, "获取成功", gin.H{
		"patents": patents,
		"total":   total,
	})
}

// GetPatentFile 获取文件地址
func GetPatentFile(c *gin.Context) {
	value := c.Query("id")
	atoi, _ := strconv.Atoi(value)
	file, err := models.GetPatentFile(atoi)
	if err != nil {
		utils.Logger.Error(err.Error())
		Resp(c, false, http.StatusInternalServerError, "获取失败", nil)
		return
	}
	Resp(c, true, http.StatusOK, "查询成功", file)
}

// DeletePatent 删除专利
func DeletePatent(c *gin.Context) {
	id := c.Query("id")

	atoi, _ := strconv.Atoi(id)
	if err := models.DeletePatent(atoi); err != nil {
		utils.Logger.Error(err.Error())
		Resp(c, false, http.StatusInternalServerError, "删除失败", nil)
		return
	}
	Resp(c, true, http.StatusOK, "删除成功", nil)
}

// UpdatePatentStatus 更新审核状态
func UpdatePatentStatus(c *gin.Context) {
	value := c.PostForm("patent_id")
	form := c.PostForm("reviewer_id")
	comment := c.PostForm("comment")
	v := c.PostForm("status")
	patent_id, _ := strconv.Atoi(value)
	reviewer_id, _ := strconv.Atoi(form)
	status, _ := strconv.Atoi(v)
	//审核没通过
	if err := models.UpdatePatentStatus(patent_id, reviewer_id, comment, status); err != nil {
		utils.Logger.Error(err.Error())
		Resp(c, false, http.StatusInternalServerError, "更新失败", nil)
		return
	}
	Resp(c, true, http.StatusOK, "更新成功", nil)
}

// GetAllPatentFees 获取所有专利年费服务方法
func GetAllPatentFees(c *gin.Context) {
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

	fees, total, err := models.GetAllPatentFees(keyword, status, pageInt, pageSizeInt)
	if err != nil {
		utils.Logger.Error(err.Error())
		Resp(c, false, http.StatusInternalServerError, "获取专利年费失败", nil)
		return
	}

	Resp(c, true, http.StatusOK, "获取专利年费成功", gin.H{
		"fees":  fees,
		"total": total,
	})
}

// GetMonthlyPatentFeeStatsService 获取本月专利年费统计信息的服务方法
func GetMonthlyPatentFeeStatsService(c *gin.Context) {
	pendingCount, paidAmount, overdueCount, totalAmount, err := models.GetMonthlyPatentFeeStats()
	if err != nil {
		utils.Logger.Error(err.Error())
		Resp(c, false, http.StatusInternalServerError, "获取本月专利年费统计信息失败", nil)
		return
	}

	Resp(c, true, http.StatusOK, "获取本月专利年费统计信息成功", gin.H{
		"pending_count": pendingCount,
		"paid_amount":   paidAmount,
		"overdue_count": overdueCount,
		"total_amount":  totalAmount,
	})
}
