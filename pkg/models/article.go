package models

import (
	"errors"
	"gorm.io/gorm"
	"intellectual_property/pkg/utils"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Article 著作核心模型
// 包含著作基本信息和审批流程管理字段
type Article struct {
	ID            int       `json:"id" gorm:"primaryKey;autoIncrement;type:bigint"`
	ArticleType   int       `json:"article_type" gorm:"type:int;comment:著作类型代码"`
	Title         string    `json:"title" gorm:"type:varchar(255);comment:著作全称"`
	Abstract      string    `json:"abstract" gorm:"type:text;comment:详细摘要"`
	ApplyDate     time.Time `json:"apply_date" gorm:"type:date;comment:申请日期"`
	AttachmentUrl string    `json:"attachment_url" gorm:"type:varchar(255);comment:附件存储路径"`

	// 作者管理字段
	FirstAuthorID int             `json:"first_author_id" gorm:"type:bigint;comment:第一作者ID"`
	FirstAuthor   User            `json:"first_author" gorm:"foreignKey:FirstAuthorID;comment:第一作者详细信息"`
	Authors       []ArticleAuthor `json:"authors" gorm:"foreignKey:ArticleID;comment:所有作者关联记录"`

	// 审批流程字段
	CurrentStep    int `json:"current_step" gorm:"type:int;comment:当前审批步骤(1=初审,2=终审)"`
	ApprovalStatus int `json:"approval_status" gorm:"type:int;comment:整体审批状态(0=进行中,1=通过,2=驳回)"`

	// 初审相关字段
	InitialReviewerID int       `json:"initial_reviewer_id" gorm:"type:bigint;comment:初审人ID"`
	InitialComment    string    `json:"initial_comment" gorm:"type:text;comment:初审意见"`
	InitialSubmitTime time.Time `json:"initial_submit_time" gorm:"type:datetime;comment:初审提交时间"`
	InitialStatus     bool      `json:"initial_status" gorm:"comment:初审状态"`

	// 终审相关字段
	FinalReviewerID int       `json:"final_reviewer_id" gorm:"type:bigint;comment:终审人ID"`
	FinalComment    string    `json:"final_comment" gorm:"type:text;comment:终审意见"`
	FinalSubmitTime time.Time `json:"final_submit_time" gorm:"type:datetime;comment:终审提交时间"`
	FinalStatus     bool      `json:"final_status" gorm:"comment:终审状态"`

	CreatedAt         time.Time
	UpdatedAt         time.Time
	ApplicationNumber string `json:"application_number" gorm:"type:varchar(255);comment:著作申请号"`
}

// ArticleAuthor 著作-作者关联模型
// 记录著作与作者的关联关系及作者角色
type ArticleAuthor struct {
	ArticleID     int  `json:"article_id" gorm:"primaryKey;type:bigint;comment:著作ID"`
	UserID        int  `json:"user_id" gorm:"primaryKey;type:bigint;comment:用户ID"`
	IsFirstAuthor bool `json:"is_first_author" gorm:"comment:是否第一作者"`
}

// ArticleDB 全局数据库连接实例
// 使用utils包中初始化的数据库连接，用于执行数据库操作
var ArticleDB *gorm.DB = utils.DB

// NewArticle 创建著作实例
// 参数：
//   - articleType: 著作类型代码
//   - title: 著作标题
//   - authorIDs: 所有作者ID列表
//   - firstAuthorID: 第一作者ID（必须包含在authorIDs中）
//   - abstract: 摘要内容
//   - attachmentUrl: 附件地址
func NewArticle(
	articleType int,
	title string,
	authorIDs []int,
	firstAuthorID int,
	abstract string,
) (*Article, *ArticleFee, error) {

	// 验证第一作者合法性
	if !contains(authorIDs, firstAuthorID) {
		return nil, nil, errors.New("第一作者必须包含在作者列表中")
	}

	// 初始化作者关联记录
	var authors []ArticleAuthor
	for _, uid := range authorIDs {
		authors = append(authors, ArticleAuthor{
			UserID:        uid,
			IsFirstAuthor: uid == firstAuthorID,
		})
	}
	// 初始化年费对象
	var reviewFee float64
	switch articleType {
	case PracticalInvention:
		reviewFee = 1
	case PatentInvention:
		reviewFee = 2
	case AppearanceDesign:
		reviewFee = 1
	}

	patentFee := &ArticleFee{
		ArticleID:       0, // 后续在创建专利时更新
		ReviewFee:       reviewFee,
		IsPaid:          false,
		CreatedAt:       time.Now(),
		PaymentDeadline: time.Now().AddDate(0, 1, 0), //一个月
		Status:          0,
	}
	return &Article{
		ArticleType:    articleType,
		Title:          title,
		Abstract:       abstract,
		FirstAuthorID:  firstAuthorID,
		Authors:        authors,
		ApplyDate:      time.Now(),
		CurrentStep:    1,
		ApprovalStatus: 0,
	}, patentFee, nil
}

// contains 检查切片是否包含指定元素
func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// CreateArticleService 创建著作服务方法
// 使用事务保证数据一致性：
// 1. 创建主记录
// 2. 创建作者关联记录
// 3. 更新第一作者外键
func (article *Article) CreateArticleService(articlefee *ArticleFee) error {
	return ArticleDB.Transaction(func(tx *gorm.DB) error {
		// 创建主记录
		if err := tx.Create(article).Error; err != nil {
			return err
		}

		// 更新专利年费对象的 PatentID
		articlefee.ArticleID = article.ID

		// 保存专利年费记录
		if err := tx.Create(articlefee).Error; err != nil {
			return err
		}
		// 清理旧关联记录
		if err := tx.Where("article_id = ?", article.ID).Delete(&ArticleAuthor{}).Error; err != nil {
			return err
		}

		// 准备关联记录
		authors := make([]ArticleAuthor, len(article.Authors))
		for i, a := range article.Authors {
			authors[i] = ArticleAuthor{
				ArticleID:     article.ID,
				UserID:        a.UserID,
				IsFirstAuthor: a.UserID == article.FirstAuthorID,
			}
		}

		// 批量创建新关联记录
		if err := tx.Create(authors).Error; err != nil {
			return err
		}

		return nil
	})
}

// UpdateArticleUrl 更新资源url
func (article *Article) UpdateArticleUrl() error {
	return ArticleDB.Model(&Article{}).Where("id =?", article.ID).Select("attachment_url").Updates(article).Error
}

// GetAllArticles 获取所有著作及其关联信息
// 支持分页和预加载关联数据
// 参数：
//   - page: 页码（从1开始）
//   - pageSize: 每页数量
//   - preloadAuthors: 是否预加载作者详细信息
//
// 返回：
//   - []Article 著作列表
//   - int 总记录数
//   - error 错误信息
func GetAllArticles(
	keyword string,
	approvalStatus int,
	articleType int,
	page int,
	pageSize int,
	preloadAuthors bool,
) ([]Article, int64, error) {
	var articles []Article
	var total int64

	query := ArticleDB.Model(&Article{})

	// 关键词模糊查询（修复后）
	if keyword != "" {
		keyword = "%" + keyword + "%"
		query = query.Where(
			"title LIKE ? OR EXISTS ("+
				"SELECT 1 FROM article_authors "+
				"JOIN users ON article_authors.user_id = users.id "+
				"WHERE article_authors.article_id = articles.id "+
				"AND (users.user_name LIKE ?)"+
				")",
			keyword,
			keyword,
		)
	}

	// 审核状态过滤
	if approvalStatus >= 0 {
		query = query.Where("approval_status = ?", approvalStatus)
	}

	// 著作类型过滤
	if articleType > 0 {
		query = query.Where("article_type = ?", articleType)
	}

	// 预加载作者信息
	if preloadAuthors {
		query = query.
			Preload("FirstAuthor").
			Preload("Authors")
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页处理
	if page > 0 && pageSize > 0 {
		query = query.Offset((page - 1) * pageSize).Limit(pageSize)
	}

	// 获取结果
	if err := query.Order("created_at DESC").Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}

// GetArticleFile 根据第文章id来查询所有文件地址
func GetArticleFile(id int) ([]string, error) {
	var ans []string
	var article Article
	PatentDB.Model(&Article{}).Select("application_number,attachment_url").Where("id =?", id).Find(&article)

	baseurl := utils.NgX.LocationArticle + "/" + article.ApplicationNumber
	dir, err := os.ReadDir(baseurl)
	if err != nil {
		return nil, err
	}
	for _, file := range dir {
		ans = append(ans, "/article/"+article.ApplicationNumber+"/"+file.Name())
	}
	return ans, nil
}

// DeleteArticle 删除著作及其关联数据
// 参数：
//   - id: 著作ID
//
// 返回：
//   - error 错误信息
func DeleteArticle(id int) error {
	// 开启数据库事务
	err := ArticleDB.Transaction(func(tx *gorm.DB) error {
		// 1. 删除作者关联记录
		if err := tx.Where("article_id = ?", id).Delete(&ArticleAuthor{}).Error; err != nil {
			return err
		}

		// 2. 删除主记录（修正后的版本）
		if err := tx.Where("id = ?", id).Delete(&Article{}).Error; err != nil {
			return err
		}

		return nil
	})

	// 3. 清理文件系统（在事务成功后执行）
	if err == nil {
		if cleanErr := cleanArticleFiles(id); cleanErr != nil {
			log.Printf("文件清理失败: %v", cleanErr)
		}
	}

	return err
}

// cleanArticleFiles 清理著作相关文件
func cleanArticleFiles(id int) error {
	// 构建安全路径
	safePath := filepath.Join(utils.NgX.LocationArticle, "/", strconv.Itoa(id))

	// 验证路径是否在允许的根目录下（防止路径遍历攻击）
	if !strings.HasPrefix(filepath.Clean(safePath)+string(filepath.Separator),
		filepath.Clean(utils.NgX.LocationArticle)+string(filepath.Separator)) {
		return errors.New("非法路径")
	}

	// 删除目录及其内容
	return os.RemoveAll(safePath)
}

// GetArticleById 根据id查询Article
func GetArticleById(id int) (Article, error) {
	var article Article
	if err := ArticleDB.Model(&Article{}).Select("current_step", "id").Where("id = ?", id).Find(&article).Error; err != nil {
		return Article{}, err
	}
	return article, nil
}

// UpdateArticleStatus 更新状态
func UpdateArticleStatus(
	articleId int,
	ReviewerID int,
	Comment string,
	Status int, // 0=驳回，1=通过
) error {
	// 查询完整文章信息
	article, err := GetArticleById(articleId)
	if err != nil {
		return err
	}

	now := time.Now()
	updates := make(map[string]interface{})

	if article.CurrentStep == ApprovalStepInitial {
		// 初审逻辑
		article.InitialReviewerID = ReviewerID
		article.InitialComment = Comment
		article.InitialSubmitTime = now
		article.InitialStatus = Status == 1

		if Status == 1 {
			// 初审通过，进入终审
			article.CurrentStep = ApprovalStepFinal
		} else {
			// 初审驳回
			article.ApprovalStatus = ApprovalStatusRejected
		}

		updates = map[string]interface{}{
			"current_step":        article.CurrentStep,
			"approval_status":     article.ApprovalStatus,
			"initial_reviewer_id": article.InitialReviewerID,
			"initial_comment":     article.InitialComment,
			"initial_submit_time": article.InitialSubmitTime,
			"initial_status":      article.InitialStatus,
		}
	} else if article.CurrentStep == ApprovalStepFinal {
		// 终审逻辑
		article.FinalReviewerID = ReviewerID
		article.FinalComment = Comment
		article.FinalSubmitTime = now
		article.FinalStatus = (Status == 1)

		if Status == 1 {
			// 终审通过
			article.ApprovalStatus = ApprovalStatusApproved
		} else {
			// 终审驳回
			article.ApprovalStatus = ApprovalStatusRejected
		}

		updates = map[string]interface{}{
			"approval_status":   article.ApprovalStatus,
			"final_reviewer_id": article.FinalReviewerID,
			"final_comment":     article.FinalComment,
			"final_submit_time": article.FinalSubmitTime,
			"final_status":      article.FinalStatus,
		}
	}

	// 执行数据库更新
	if err := ArticleDB.Model(&Article{}).
		Where("id = ?", articleId).
		Updates(updates).Error; err != nil {
		return err
	}

	return nil
}
