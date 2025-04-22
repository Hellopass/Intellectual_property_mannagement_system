package service

import (
	"github.com/gin-gonic/gin"
	"intellectual_property/pkg/models"
	"intellectual_property/pkg/utils"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Re struct {
	article  models.Article
	filepath []string
}

// CreateArticle 创建著作
// @Summary 创建新著作
// @Tags 著作管理
// @Accept multipart/form-data
// @Param Authorization header string true "Bearer Token"
// @Param articleType formData string true "著作类型（中文或代码）"
// @Param title formData string true "著作全称"
// @Param abstract formData string true "详细摘要"
// @Param authors formData []int true "作者ID列表" collectionFormat(multi)
// @Param firstAuthorId formData int true "第一作者ID"
// @Param files formData []file true "相关文件"
// @Success 201 {object} Response "创建成功"
// @Failure 400 {object} Response "参数错误"
// @Failure 401 {object} Response "认证失败"
// @Failure 500 {object} Response "服务器错误"
// @Router /articles [post]
func CreateArticle(c *gin.Context) {
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
	articleTypeStr := c.PostForm("articleType")
	title := c.PostForm("title")
	abstract := c.PostForm("abstract")
	firstAuthorID, _ := strconv.Atoi(c.PostForm("firstAuthorId"))
	authorIDs := utils.ConvertStringSliceToInt(c.PostFormArray("authors"))

	// 3. 参数验证
	if title == "" || abstract == "" || len(authorIDs) == 0 {
		Resp(c, false, http.StatusBadRequest, "必要参数不能为空", nil)
		return
	}

	// 3. 创建著作记录
	articleType, err := models.ParseArticleType(articleTypeStr)
	if err != nil {
		Resp(c, false, http.StatusBadRequest, "无效的著作类型", nil)
		return
	}
	article, err := models.NewArticle(
		articleType,
		title,
		authorIDs,
		firstAuthorID,
		abstract,
	)
	article.ApplyDate = time.Now()
	if err != nil {
		logger.Error("参数验证失败: " + err.Error())
		Resp(c, false, http.StatusBadRequest, err.Error(), nil)
		return
	}

	if err := article.CreateArticleService(); err != nil {
		logger.Error("数据库操作失败: " + err.Error())
		Resp(c, false, http.StatusInternalServerError, "创建著作失败", nil)
		return
	}
	//4，文件上传
	form, err2 := c.MultipartForm()
	if err2 != nil {
		logger.Error(err2.Error())
		Resp(c, false, CodeError, "上传相关文件失败", nil)
		return
	}
	Files := form.File["files"] //拿到文件
	var directory string
	directory = strconv.Itoa(article.ID)
	url := utils.NgX.LocationArticle + "/" + directory
	for _, f := range Files {
		//保存文件 地址---原始地址/申请号文件/文件

		err = c.SaveUploadedFile(f, url+"/"+f.Filename)
		if err != nil {
			logger.Error(err.Error())
			Resp(c, false, CodeError, "上传相关文件失败", nil)
			return
		}
	}
	//在更新一下数据库--只跟新文件位置
	article.AttachmentUrl = url
	err3 := article.UpdateArticleUrl()
	if err3 != nil {
		logger.Error(err3.Error())
		Resp(c, false, CodeError, "上传相关文件失败", nil)
		return
	}
	// 5. 返回成功响应
	Resp(c, true, http.StatusCreated, "著作创建成功", gin.H{
		"articleId":   article.ID,
		"title":       article.Title,
		"firstAuthor": article.FirstAuthorID,
		"fileCount":   len(Files),
	})
}

func GetAllArticles(c *gin.Context) {
	page := c.Query("page")
	pageSize := c.Query("pageSize")
	keyword := c.Query("search")
	status := c.Query("status")
	types := c.Query("type")
	var (
		a int
		b int
	)
	if types == "" || status == "" {
		a = 0
		b = 0
	} else {
		a, _ = strconv.Atoi(status)
		b, _ = models.ParseArticleType(types)
	}
	page_a, _ := strconv.Atoi(page)
	page_size_a, _ := strconv.Atoi(pageSize)

	articles, total, err := models.GetAllArticles(keyword, a, b, page_a, page_size_a, true)
	if err != nil {
		logger.Error(err.Error())
		Resp(c, false, http.StatusInternalServerError, "获取失败", nil)
		return
	}

	Resp(c, true, http.StatusOK, "获取成功", gin.H{
		"articles": articles,
		"total":    total,
	})
}

// GetArticleFile 获取文件地址
func GetArticleFile(c *gin.Context) {
	value := c.Query("id")
	atoi, _ := strconv.Atoi(value)
	file, err := models.GetArticleFile(atoi)
	if err != nil {
		logger.Error(err.Error())
		Resp(c, false, http.StatusInternalServerError, "获取失败", nil)
		return
	}
	Resp(c, true, http.StatusOK, "查询成功", file)
}

// DeleteArticle 删除著作
func DeleteArticle(c *gin.Context) {
	id := c.Query("id")

	atoi, _ := strconv.Atoi(id)
	if err := models.DeleteArticle(atoi); err != nil {
		if err != nil {
			logger.Error(err.Error())
			Resp(c, false, http.StatusInternalServerError, "删除失败", nil)
			return
		}
	}
	Resp(c, true, http.StatusOK, "删除成功", nil)

}

// UpdateArticleStatus 更新审核状态
func UpdateArticleStatus(c *gin.Context) {
	value := c.PostForm("article_id")
	form := c.PostForm("reviewer_id")
	comment := c.PostForm("comment")
	v := c.PostForm("status")
	article_id, _ := strconv.Atoi(value)
	reviewer_id, _ := strconv.Atoi(form)
	status, _ := strconv.Atoi(v)
	//审核没通过
	if err := models.UpdateArticleStatus(article_id, reviewer_id, comment, status); err != nil {
		logger.Error(err.Error())
		Resp(c, false, http.StatusInternalServerError, "更新失败", nil)
		return
	}
	Resp(c, true, http.StatusOK, "更新成功", nil)

}
