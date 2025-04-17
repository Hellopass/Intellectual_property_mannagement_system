package dispose

import (
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
