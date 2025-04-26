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

// 专利类型编码
const (
	PatentInvention    = iota // 发明专利
	PracticalInvention        // 实用新型
	AppearanceDesign          // 外观设计
)

// Patent 专利核心模型
// 包含专利基本信息和审批流程管理字段
type Patent struct {
	ID                int       `json:"id" gorm:"primaryKey;autoIncrement;type:bigint"`
	PatentType        int       `json:"patent_type" gorm:"type:int;comment:专利类型代码"`
	Title             string    `json:"title" gorm:"type:varchar(255);comment:专利全称"`
	Abstract          string    `json:"abstract" gorm:"type:text;comment:详细摘要"`
	ApplyDate         time.Time `json:"apply_date" gorm:"type:date;comment:申请日期"`
	AttachmentUrl     string    `json:"attachment_url" gorm:"type:varchar(255);comment:附件存储路径"`
	ApplicationNumber string    `json:"application_number" gorm:"type:varchar(255);comment:专利申请号"` // 新增申请号字段

	// 作者管理字段
	FirstAuthorID int            `json:"first_author_id" gorm:"type:bigint;comment:第一作者ID"`
	FirstAuthor   User           `json:"first_author" gorm:"foreignKey:FirstAuthorID;comment:第一作者详细信息"`
	Authors       []PatentAuthor `json:"authors" gorm:"foreignKey:PatentID;comment:所有作者关联记录"`

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

	CreatedAt time.Time
	UpdatedAt time.Time
}

// PatentAuthor 专利-作者关联模型
// 记录专利与作者的关联关系及作者角色
type PatentAuthor struct {
	PatentID      int  `json:"patent_id" gorm:"primaryKey;type:bigint;comment:专利ID"`
	UserID        int  `json:"user_id" gorm:"primaryKey;type:bigint;comment:用户ID"`
	IsFirstAuthor bool `json:"is_first_author" gorm:"comment:是否第一作者"`
}

// PatentDB 全局数据库连接实例
// 使用utils包中初始化的数据库连接，用于执行数据库操作
var PatentDB *gorm.DB = utils.DB

// NewPatent 创建专利实例
// 参数：
//   - patentType: 专利类型代码
//   - title: 专利标题
//   - authorIDs: 所有作者ID列表
//   - firstAuthorID: 第一作者ID（必须包含在authorIDs中）
//   - abstract: 摘要内容
//   - attachmentUrl: 附件地址
func NewPatent(
	patentType int,
	title string,
	authorIDs []int,
	firstAuthorID int,
	abstract string,
	applicationNumber string, // 新增申请号参数
) (*Patent, *PatentFee, error) {

	// 验证第一作者合法性
	if !contains(authorIDs, firstAuthorID) {
		return nil, nil, errors.New("第一作者必须包含在作者列表中")
	}

	// 初始化作者关联记录
	var authors []PatentAuthor
	for _, uid := range authorIDs {
		authors = append(authors, PatentAuthor{
			UserID:        uid,
			IsFirstAuthor: uid == firstAuthorID,
		})
	}

	// 初始化年费对象
	var reviewFee float64
	switch patentType {
	case PracticalInvention:
		reviewFee = 1
	case PatentInvention:
		reviewFee = 2
	case AppearanceDesign:
		reviewFee = 1
	}

	patentFee := &PatentFee{
		PatentID:        0, // 后续在创建专利时更新
		ReviewFee:       reviewFee,
		IsPaid:          false,
		CreatedAt:       time.Now(),
		PaymentDeadline: time.Now().AddDate(0, 1, 0), //一个月
		Status:          0,
	}

	return &Patent{
		PatentType:        patentType,
		Title:             title,
		Abstract:          abstract,
		FirstAuthorID:     firstAuthorID,
		Authors:           authors,
		ApplyDate:         time.Now(),
		CurrentStep:       1,
		ApprovalStatus:    0,
		ApplicationNumber: applicationNumber, // 初始化申请号字段
	}, patentFee, nil
}

// CreatePatentService 创建专利服务方法
// 使用事务保证数据一致性：
// 1. 创建主记录
// 2. 创建作者关联记录
// 3. 更新第一作者外键
func (patent *Patent) CreatePatentService(patentFee *PatentFee) error {
	return PatentDB.Transaction(func(tx *gorm.DB) error {
		// 创建主记录
		if err := tx.Create(patent).Error; err != nil {
			return err
		}
		// 更新专利年费对象的 PatentID
		patentFee.PatentID = patent.ID

		// 保存专利年费记录
		if err := tx.Create(patentFee).Error; err != nil {
			return err
		}
		// 清理旧关联记录
		if err := tx.Where("patent_id = ?", patent.ID).Delete(&PatentAuthor{}).Error; err != nil {
			return err
		}

		// 准备关联记录
		authors := make([]PatentAuthor, len(patent.Authors))
		for i, a := range patent.Authors {
			authors[i] = PatentAuthor{
				PatentID:      patent.ID,
				UserID:        a.UserID,
				IsFirstAuthor: a.UserID == patent.FirstAuthorID,
			}
		}

		// 批量创建新关联记录
		if err := tx.Create(authors).Error; err != nil {
			return err
		}

		return nil
	})
}

// UpdatePatentUrl 更新资源url
func (patent *Patent) UpdatePatentUrl() error {
	return PatentDB.Model(&Patent{}).Where("id =?", patent.ID).Select("attachment_url").Updates(patent).Error
}

// GetAllPatents 获取所有专利及其关联信息
// 支持分页和预加载关联数据
// 参数：
//   - keyword: 关键词
//   - approvalStatus: 审核状态
//   - patentType: 专利类型
//   - page: 页码（从1开始）
//   - pageSize: 每页数量
//   - preloadAuthors: 是否预加载作者详细信息
//
// 返回：
//   - []Patent 专利列表
//   - int 总记录数
//   - error 错误信息
func GetAllPatents(
	keyword string,
	approvalStatus int,
	patentType int,
	page int,
	pageSize int,
	preloadAuthors bool,
) ([]Patent, int64, error) {
	var patents []Patent
	var total int64

	query := PatentDB.Model(&Patent{})

	// 关键词模糊查询（添加申请号的模糊查询）
	if keyword != "" {
		keyword = "%" + keyword + "%"
		query = query.Where(
			"title LIKE ? OR application_number LIKE ? OR EXISTS ("+
				"SELECT 1 FROM patent_authors "+
				"JOIN users ON patent_authors.user_id = users.id "+
				"WHERE patent_authors.patent_id = patents.id "+
				"AND (users.user_name LIKE ?)"+
				")",
			keyword,
			keyword,
			keyword,
		)
	}

	// 审核状态过滤
	if approvalStatus >= 0 {
		query = query.Where("approval_status = ?", approvalStatus)
	}

	// 专利类型过滤
	if patentType > 0 {
		query = query.Where("patent_type = ?", patentType)
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
	if err := query.Order("created_at DESC").Find(&patents).Error; err != nil {
		return nil, 0, err
	}

	return patents, total, nil
}

// GetPatentFile 根据专利id来查询所有文件地址
func GetPatentFile(id int) ([]string, error) {
	var ans []string
	var patent Patent
	PatentDB.Model(&Patent{}).Select("application_number").Where("id =?", id).Find(&patent)

	baseurl := utils.NgX.LocationPatent + "/" + patent.ApplicationNumber
	dir, err := os.ReadDir(baseurl)
	if err != nil {
		return nil, err
	}
	for _, file := range dir {
		ans = append(ans, "/patent/"+patent.ApplicationNumber+"/"+file.Name())
	}

	return ans, nil
}

// DeletePatent 删除专利及其关联数据
// 参数：
//   - id: 专利ID
//
// 返回：
//   - error 错误信息
func DeletePatent(id int) error {
	// 开启数据库事务
	err := PatentDB.Transaction(func(tx *gorm.DB) error {
		// 1. 删除作者关联记录
		if err := tx.Where("patent_id = ?", id).Delete(&PatentAuthor{}).Error; err != nil {
			return err
		}

		// 2. 删除主记录（修正后的版本）
		if err := tx.Where("id = ?", id).Delete(&Patent{}).Error; err != nil {
			return err
		}

		return nil
	})

	// 3. 清理文件系统（在事务成功后执行）
	if err == nil {
		if cleanErr := cleanPatentFiles(id); cleanErr != nil {
			log.Printf("文件清理失败: %v", cleanErr)
		}
	}

	return err
}

// cleanPatentFiles 清理专利相关文件
func cleanPatentFiles(id int) error {
	// 构建安全路径
	safePath := filepath.Join(utils.NgX.LocationPatent, "/", strconv.Itoa(id))

	// 验证路径是否在允许的根目录下（防止路径遍历攻击）
	if !strings.HasPrefix(filepath.Clean(safePath)+string(filepath.Separator),
		filepath.Clean(utils.NgX.LocationPatent)+string(filepath.Separator)) {
		return errors.New("非法路径")
	}

	// 删除目录及其内容
	return os.RemoveAll(safePath)
}

// GetPatentById 根据id查询Patent
func GetPatentById(id int) (Patent, error) {
	var patent Patent
	if err := PatentDB.Model(&Patent{}).Select("current_step", "id").Where("id = ?", id).Find(&patent).Error; err != nil {
		return Patent{}, err
	}
	return patent, nil
}

// UpdatePatentStatus 更新状态
func UpdatePatentStatus(
	patentId int,
	ReviewerID int,
	Comment string,
	Status int, // 0=驳回，1=通过
) error {
	// 查询完整专利信息
	patent, err := GetPatentById(patentId)
	if err != nil {
		return err
	}

	now := time.Now()
	updates := make(map[string]interface{})

	if patent.CurrentStep == ApprovalStepInitial {
		// 初审逻辑
		patent.InitialReviewerID = ReviewerID
		patent.InitialComment = Comment
		patent.InitialSubmitTime = now
		patent.InitialStatus = Status == 1

		if Status == 1 {
			// 初审通过，进入终审
			patent.CurrentStep = ApprovalStepFinal
		} else {
			// 初审驳回
			patent.ApprovalStatus = ApprovalStatusRejected
		}

		updates = map[string]interface{}{
			"current_step":        patent.CurrentStep,
			"approval_status":     patent.ApprovalStatus,
			"initial_reviewer_id": patent.InitialReviewerID,
			"initial_comment":     patent.InitialComment,
			"initial_submit_time": patent.InitialSubmitTime,
			"initial_status":      patent.InitialStatus,
		}
	} else if patent.CurrentStep == ApprovalStepFinal {
		// 终审逻辑
		patent.FinalReviewerID = ReviewerID
		patent.FinalComment = Comment
		patent.FinalSubmitTime = now
		patent.FinalStatus = Status == 1

		if Status == 1 {
			// 终审通过
			patent.ApprovalStatus = ApprovalStatusApproved
		} else {
			// 终审驳回
			patent.ApprovalStatus = ApprovalStatusRejected
		}

		updates = map[string]interface{}{
			"approval_status":   patent.ApprovalStatus,
			"final_reviewer_id": patent.FinalReviewerID,
			"final_comment":     patent.FinalComment,
			"final_submit_time": patent.FinalSubmitTime,
			"final_status":      patent.FinalStatus,
		}
	}

	// 执行数据库更新
	if err := PatentDB.Model(&Patent{}).
		Where("id = ?", patentId).
		Updates(updates).Error; err != nil {
		return err
	}

	return nil
}
