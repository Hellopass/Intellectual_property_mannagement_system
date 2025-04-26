package models

import (
	"errors"
	"intellectual_property/pkg/utils"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Trademark 商标核心模型
// 包含商标基本信息和审批流程管理字段
type Trademark struct {
	ID            int       `json:"id" gorm:"primaryKey;autoIncrement;type:bigint"`
	TrademarkType int       `json:"trademark_type" gorm:"type:int;comment:商标类型代码"`
	Title         string    `json:"title" gorm:"type:varchar(255);comment:商标全称"`
	Abstract      string    `json:"abstract" gorm:"type:text;comment:详细摘要"`
	ApplyDate     time.Time `json:"apply_date" gorm:"type:date;comment:申请日期"`
	AttachmentUrl string    `json:"attachment_url" gorm:"type:varchar(255);comment:附件存储路径"`

	// 作者管理字段
	FirstAuthorID int               `json:"first_author_id" gorm:"type:bigint;comment:第一作者ID"`
	FirstAuthor   User              `json:"first_author" gorm:"foreignKey:FirstAuthorID;comment:第一作者详细信息"`
	Authors       []TrademarkAuthor `json:"authors" gorm:"foreignKey:TrademarkID;comment:所有作者关联记录"`

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
	ApplicationNumber string `json:"application_number" gorm:"type:varchar(255);comment:商标申请号"`
}

// TrademarkAuthor 商标-作者关联模型
// 记录商标与作者的关联关系及作者角色
type TrademarkAuthor struct {
	TrademarkID   int  `json:"trademark_id" gorm:"primaryKey;type:bigint;comment:商标ID"`
	UserID        int  `json:"user_id" gorm:"primaryKey;type:bigint;comment:用户ID"`
	IsFirstAuthor bool `json:"is_first_author" gorm:"comment:是否第一作者"`
}

// TrademarkDB 全局数据库连接实例
// 使用utils包中初始化的数据库连接，用于执行数据库操作
var TrademarkDB *gorm.DB = utils.DB

// NewTrademark 创建商标实例
// 参数：
//   - trademarkType: 商标类型代码
//   - title: 商标标题
//   - authorIDs: 所有作者ID列表
//   - firstAuthorID: 第一作者ID（必须包含在authorIDs中）
//   - abstract: 摘要内容
func NewTrademark(
	trademarkType int,
	title string,
	authorIDs []int,
	firstAuthorID int,
	abstract string,
) (*Trademark, *TrademarkFee, error) {

	// 验证第一作者合法性
	if !contains(authorIDs, firstAuthorID) {
		return nil, nil, errors.New("第一作者必须包含在作者列表中")
	}

	// 初始化作者关联记录
	var authors []TrademarkAuthor
	for _, uid := range authorIDs {
		authors = append(authors, TrademarkAuthor{
			UserID:        uid,
			IsFirstAuthor: uid == firstAuthorID,
		})
	}
	// 初始化年费对象
	var reviewFee float64
	switch trademarkType {
	case PracticalInvention:
		reviewFee = 1
	case PatentInvention:
		reviewFee = 2
	case AppearanceDesign:
		reviewFee = 1
	}

	patentFee := &TrademarkFee{
		TrademarkID:     0, // 后续在创建专利时更新
		ReviewFee:       reviewFee,
		IsPaid:          false,
		CreatedAt:       time.Now(),
		PaymentDeadline: time.Now().AddDate(0, 1, 0), //一个月
		Status:          0,
	}
	return &Trademark{
		TrademarkType:  trademarkType,
		Title:          title,
		Abstract:       abstract,
		FirstAuthorID:  firstAuthorID,
		Authors:        authors,
		ApplyDate:      time.Now(),
		CurrentStep:    1,
		ApprovalStatus: 0,
	}, patentFee, nil
}

// CreateTrademarkService 创建商标服务方法
// 使用事务保证数据一致性：
// 1. 创建主记录
// 2. 创建作者关联记录
// 3. 更新第一作者外键
func (trademark *Trademark) CreateTrademarkService(fee *TrademarkFee) error {
	return TrademarkDB.Transaction(func(tx *gorm.DB) error {
		// 创建主记录
		if err := tx.Create(trademark).Error; err != nil {
			return err
		}
		// 更新商标年费对象的 trademark_id
		fee.TrademarkID = trademark.ID

		// 更新商标年费对象的
		if err := tx.Create(fee).Error; err != nil {
			return err
		}
		// 清理旧关联记录
		if err := tx.Where("trademark_id = ?", trademark.ID).Delete(&TrademarkAuthor{}).Error; err != nil {
			return err
		}

		// 准备关联记录
		authors := make([]TrademarkAuthor, len(trademark.Authors))
		for i, a := range trademark.Authors {
			authors[i] = TrademarkAuthor{
				TrademarkID:   trademark.ID,
				UserID:        a.UserID,
				IsFirstAuthor: a.UserID == trademark.FirstAuthorID,
			}
		}

		// 批量创建新关联记录
		if err := tx.Create(authors).Error; err != nil {
			return err
		}

		return nil
	})
}

// UpdateTrademarkUrl 更新资源url
func (trademark *Trademark) UpdateTrademarkUrl() error {
	return TrademarkDB.Model(&Trademark{}).Where("id =?", trademark.ID).Select("attachment_url").Updates(trademark).Error
}

// GetAllTrademarks 获取所有商标及其关联信息
// 支持分页和预加载关联数据
// 参数：
//   - keyword: 关键词
//   - approvalStatus: 审核状态
//   - trademarkType: 商标类型
//   - page: 页码（从1开始）
//   - pageSize: 每页数量
//   - preloadAuthors: 是否预加载作者详细信息
//
// 返回：
//   - []Trademark 商标列表
//   - int 总记录数
//   - error 错误信息
func GetAllTrademarks(
	keyword string,
	approvalStatus int,
	trademarkType int,
	page int,
	pageSize int,
	preloadAuthors bool,
) ([]Trademark, int64, error) {
	var trademarks []Trademark
	var total int64

	query := TrademarkDB.Model(&Trademark{})

	// 关键词模糊查询（修复后）
	if keyword != "" {
		keyword = "%" + keyword + "%"
		query = query.Where(
			"title LIKE ? OR EXISTS ("+
				"SELECT 1 FROM trademark_authors "+
				"JOIN users ON trademark_authors.user_id = users.id "+
				"WHERE trademark_authors.trademark_id = trademarks.id "+
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

	// 商标类型过滤
	if trademarkType > 0 {
		query = query.Where("trademark_type = ?", trademarkType)
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
	if err := query.Order("created_at DESC").Find(&trademarks).Error; err != nil {
		return nil, 0, err
	}

	return trademarks, total, nil
}

// GetTrademarkFile 根据商标id来查询所有文件地址
func GetTrademarkFile(id int) ([]string, error) {
	//var ans []string
	//var directory string
	//directory = strconv.Itoa(id)
	//baseurl := utils.NgX.LocationTrademark + "/" + directory
	//dir, err := os.ReadDir(baseurl)
	//if err != nil {
	//	return nil, err
	//}
	//for _, file := range dir {
	//	ans = append(ans, "/trademark/"+directory+"/"+file.Name())
	//}
	var ans []string
	var trademark Trademark
	PatentDB.Model(&Trademark{}).Select("application_number,attachment_url").Where("id =?", id).Find(&trademark)

	baseurl := utils.NgX.LocationTrademark + "/" + trademark.ApplicationNumber
	dir, err := os.ReadDir(baseurl)
	if err != nil {
		return nil, err
	}
	for _, file := range dir {
		ans = append(ans, "/trademark/"+trademark.ApplicationNumber+"/"+file.Name())
	}
	return ans, nil
}

// DeleteTrademark 删除商标及其关联数据
// 参数：
//   - id: 商标ID
//
// 返回：
//   - error 错误信息
func DeleteTrademark(id int) error {
	// 开启数据库事务
	err := TrademarkDB.Transaction(func(tx *gorm.DB) error {
		// 1. 删除作者关联记录
		if err := tx.Where("trademark_id = ?", id).Delete(&TrademarkAuthor{}).Error; err != nil {
			return err
		}

		// 2. 删除主记录（修正后的版本）
		if err := tx.Where("id = ?", id).Delete(&Trademark{}).Error; err != nil {
			return err
		}

		return nil
	})

	// 3. 清理文件系统（在事务成功后执行）
	if err == nil {
		if cleanErr := cleanTrademarkFiles(id); cleanErr != nil {
			log.Printf("文件清理失败: %v", cleanErr)
		}
	}

	return err
}

// cleanTrademarkFiles 清理商标相关文件
func cleanTrademarkFiles(id int) error {
	// 构建安全路径
	safePath := filepath.Join(utils.NgX.LocationTrademark, "/", strconv.Itoa(id))

	// 验证路径是否在允许的根目录下（防止路径遍历攻击）
	if !strings.HasPrefix(filepath.Clean(safePath)+string(filepath.Separator),
		filepath.Clean(utils.NgX.LocationTrademark)+string(filepath.Separator)) {
		return errors.New("非法路径")
	}

	// 删除目录及其内容
	return os.RemoveAll(safePath)
}

// GetTrademarkById 根据id查询Trademark
func GetTrademarkById(id int) (Trademark, error) {
	var trademark Trademark
	if err := TrademarkDB.Model(&Trademark{}).Select("current_step", "id").Where("id = ?", id).Find(&trademark).Error; err != nil {
		return Trademark{}, err
	}
	return trademark, nil
}

// UpdateTrademarkStatus 更新状态
func UpdateTrademarkStatus(
	trademarkId int,
	ReviewerID int,
	Comment string,
	Status int, // 0=驳回，1=通过
) error {
	// 查询完整商标信息
	trademark, err := GetTrademarkById(trademarkId)
	if err != nil {
		return err
	}

	now := time.Now()
	updates := make(map[string]interface{})

	if trademark.CurrentStep == ApprovalStepInitial {
		// 初审逻辑
		trademark.InitialReviewerID = ReviewerID
		trademark.InitialComment = Comment
		trademark.InitialSubmitTime = now
		trademark.InitialStatus = Status == 1

		if Status == 1 {
			// 初审通过，进入终审
			trademark.CurrentStep = ApprovalStepFinal
		} else {
			// 初审驳回
			trademark.ApprovalStatus = ApprovalStatusRejected
		}

		updates = map[string]interface{}{
			"current_step":        trademark.CurrentStep,
			"approval_status":     trademark.ApprovalStatus,
			"initial_reviewer_id": trademark.InitialReviewerID,
			"initial_comment":     trademark.InitialComment,
			"initial_submit_time": trademark.InitialSubmitTime,
			"initial_status":      trademark.InitialStatus,
		}
	} else if trademark.CurrentStep == ApprovalStepFinal {
		// 终审逻辑
		trademark.FinalReviewerID = ReviewerID
		trademark.FinalComment = Comment
		trademark.FinalSubmitTime = now
		trademark.FinalStatus = (Status == 1)

		if Status == 1 {
			// 终审通过
			trademark.ApprovalStatus = ApprovalStatusApproved
		} else {
			// 终审驳回
			trademark.ApprovalStatus = ApprovalStatusRejected
		}

		updates = map[string]interface{}{
			"approval_status":   trademark.ApprovalStatus,
			"final_reviewer_id": trademark.FinalReviewerID,
			"final_comment":     trademark.FinalComment,
			"final_submit_time": trademark.FinalSubmitTime,
			"final_status":      trademark.FinalStatus,
		}
	}

	// 执行数据库更新
	if err := TrademarkDB.Model(&Trademark{}).
		Where("id = ?", trademarkId).
		Updates(updates).Error; err != nil {
		return err
	}

	return nil
}
