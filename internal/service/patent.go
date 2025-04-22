package service

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"intellectual_property/pkg/models"
	"intellectual_property/pkg/utils"
	"net/http"
	"strconv"
	"time"
)

// AddPatent 新建申请
func AddPatent(c *gin.Context) {
	//拿到基础信息
	var p models.Patent
	if err := c.ShouldBind(&p); err != nil {
		logger.Error(err.Error())
		Resp(c, false, CodeError, "新建申请失败", nil)
		return
	}

	Time := time.Now()
	currentYear := Time.Year()                                          // 获取当前年份
	yearStr := strconv.Itoa(currentYear)                                // 将年份转换为字符串
	applyNo, err2 := utils.GenerateApplyNo("CN", yearStr, p.PatentType) //申请申请号
	if err2 != nil {
		logger.Error(err2.Error())
		Resp(c, false, CodeError, "新建申请失败", nil)
		return
	}
	format := Time.Format("2006-01-02")
	parse, err22 := time.Parse("2006-01-02", format)
	if err22 != nil {
		logger.Error(err22.Error())
		return
	}
	p.ApplyData = parse //申请时间
	p.ApplyNo = applyNo //申请号
	id, err2 := models.GetUserByID(p.UserID)
	if err2 != nil {
		logger.Error(err2.Error())
		Resp(c, false, CodeError, "新建申请失败", nil)
		return
	}
	p.User = id
	if err := models.CreatePatent(&p); err != nil {
		logger.Error(err.Error())
		Resp(c, false, CodeError, "新建申请失败", nil)
		return
	}

	//成功返回申请号
	Resp(c, true, http.StatusOK, "新建申请成功", gin.H{
		"apply_no": applyNo,
	})
	return
}

// UploadPatentFile 上传相关文件--多文件上传
func UploadPatentFile(c *gin.Context) {
	value := c.PostForm("apply_no") //申请号
	form, err := c.MultipartForm()
	if err != nil {
		logger.Error(err.Error())
		Resp(c, false, CodeError, "上传相关文件失败", nil)
		return
	}
	Files := form.File["file"] //拿到文件
	for _, f := range Files {
		//保存文件 地址---原始地址/申请号文件/文件
		err = c.SaveUploadedFile(f, utils.NgX.LocationDocs+"/"+value+"/"+f.Filename)
		if err != nil {
			logger.Error(err.Error())
			Resp(c, false, CodeError, "上传相关文件失败", nil)
			return
		}
	}

	Resp(c, true, http.StatusOK, "上传文件成功", nil)
}

// FindPatentS 获取所有专利信息
func FindPatentS(c *gin.Context) {
	models.CreatePatentTable()
	information, err := models.GetPatentInformation()
	if err != nil {
		logger.Error(err.Error())
		Resp(c, false, CodeError, "查询失败", nil)
		return
	}
	for i, v := range information {
		id, err := models.GetUserByID(v.UserID)
		if err != nil {
			logger.Error(err.Error())
			Resp(c, false, CodeError, "查询失败", nil)
			return
		}
		information[i].User = id
	}
	Resp(c, true, http.StatusOK, "查询成功", information)
}

// FindPatentFuzzy 模糊查询
func FindPatentFuzzy(c *gin.Context) {
	keyword := c.Query("keyword")
	s := c.Query("status")
	//转int
	if s == "" {
		s = "0"
	}
	status, err := strconv.Atoi(s)
	if err != nil {
		logger.Error(err.Error())
		Resp(c, false, CodeError, "查询失败", nil)
		return
	}
	//模糊查询
	fuzzy, err1 := models.FindPatentFuzzy(keyword, status)
	if err1 != nil {
		logger.Error(err1.Error())
		Resp(c, false, CodeError, "查询失败", nil)
		return
	}
	Resp(c, true, http.StatusOK, "查询成功", fuzzy)
}

// DelPatent 删除专利信息
func DelPatent(c *gin.Context) {
	applyNno := c.Query("apply_no")
	//根据申请号删除
	if err := models.DeletePatent(applyNno); err != nil {
		logger.Error(err.Error())
		Resp(c, false, CodeError, "删除失败", nil)
		return
	}
	Resp(c, true, http.StatusOK, "删除成功", nil)
}

// GetPatentFile 申请号对应的文件
func GetPatentFile(c *gin.Context) {
	applicationNo := c.Query("applicationNo") //拿到申请号
	filepath, err := models.GetPatentFile(applicationNo)
	if err != nil {
		logger.Error(err.Error())
		Resp(c, false, CodeError, "获取文件失败", nil)
		return
	}
	Resp(c, true, http.StatusOK, "获取文件成功", filepath)

}

// UpdateStatus 更新状态
func UpdateStatus(c *gin.Context) {
	applyno := c.PostForm("apply_no") //申请号
	status := c.PostForm("status")    //状态
	atoi, err := strconv.Atoi(status)
	if err != nil {
		logger.Error(err.Error())
		Resp(c, false, CodeError, "通过申请失败", nil)
		return
	}
	err = models.UpdateStatusByApplicationNumber(applyno, atoi)
	if err != nil {
		logger.Error(err.Error())
		Resp(c, false, CodeError, "通过申请失败", nil)
		return
	}

	Resp(c, true, http.StatusOK, "通过申请成功", nil)
}

// GetFeeStatistics 统计
func GetFeeStatistics(c *gin.Context) {
	statistics, err := models.GetFeeStatistics()
	if err != nil {
		logger.Error(err.Error())
		Resp(c, false, CodeError, "初始化失败", nil)
		return
	}
	Resp(c, true, http.StatusOK, "初始化成功", statistics)
}

// GetAllPatentFees 获取所有专利年费记录分页
func GetAllPatentFees(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	pagination, fees, err := models.GetAllPatentFees(page, size)
	if err != nil {
		Resp(c, false, CodeError, "查询年费失败", nil)
		return
	}
	Resp(c, true, http.StatusOK, "初始化成功", gin.H{
		"pagination": pagination,
		"fees":       fees,
	})

}

// UpdatePatentFeeByApplyNo 根据申请号更新专利费用记录
func UpdatePatentFeeByApplyNo(c *gin.Context) {
	apply_no := c.PostForm("apply_no") //申请号
	amount, err := strconv.ParseFloat(c.PostForm("amount"), 64)
	if err != nil {
		logger.Error(err.Error())
		Resp(c, false, CodeError, "更新失败", nil)
		return
	}
	err = models.UpdatePatentFeeByApplyNo(apply_no, map[string]interface{}{
		"amount": amount,
	})
	if err != nil {
		logger.Error(err.Error())
		Resp(c, false, CodeError, "更新失败", nil)
		return
	}
	Resp(c, true, http.StatusOK, "更新成功", nil)
}

// GetPatentFeesByFilters 根据状态、关键字模糊查询（申请号/专利名称）分页查询
func GetPatentFeesByFilters(c *gin.Context) {
	// 获取查询参数
	statusStr := c.DefaultQuery("status", "")       // 缴费状态
	keyword := c.DefaultQuery("keyword", "")        // 合并后的关键字查询
	pageStr := c.DefaultQuery("page", "1")          // 页码
	pageSizeStr := c.DefaultQuery("pageSize", "10") // 页长

	// 解析分页参数
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		logger.Error(fmt.Sprintf("无效页码参数: %s", pageStr))
		Resp(c, false, http.StatusBadRequest, "无效的页码", nil)
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize <= 0 {
		logger.Error(fmt.Sprintf("无效页长参数: %s", pageSizeStr))
		Resp(c, false, http.StatusBadRequest, "无效的页长", nil)
		return
	}

	// 处理状态参数
	var statusPtr *models.PaymentStatus
	if statusStr != "" {
		statusInt, err := strconv.Atoi(statusStr)
		if err != nil || statusInt < int(models.StatusUnpaid) || statusInt > int(models.StatusOverdue) {
			logger.Error(fmt.Sprintf("无效状态参数: %s", statusStr))
			Resp(c, false, http.StatusBadRequest, "无效的缴费状态", nil)
			return
		}
		status := models.PaymentStatus(statusInt)
		statusPtr = &status
	}

	// 执行查询
	information, fees, err := models.GetPatentFeesByFilters(
		statusPtr,
		keyword,
		page,
		pageSize,
	)

	// 处理错误
	if err != nil {
		logger.Error("查询失败: " + err.Error())
		Resp(c, false, http.StatusInternalServerError, "查询失败", nil)
		return
	}

	// 返回结果
	Resp(c, true, http.StatusOK, "查询成功", gin.H{
		"data":       fees,
		"pagination": information,
	})
}

// GetAnalysis 抓专利数据分析
func GetAnalysis(c *gin.Context) {
	// 获取各分析结果
	yearly, _ := models.GetYearlyTrends()
	types, _ := models.GetTypeDistribution()
	applicants, _ := models.GetTopApplicants()
	tech, _ := models.GetTechDomains()

	Resp(c, true, http.StatusOK, "查询成功", gin.H{
		"yearly_trends":     yearly,
		"type_distribution": types,
		"top_applicants":    applicants,
		"tech_domains":      tech,
	})
}
