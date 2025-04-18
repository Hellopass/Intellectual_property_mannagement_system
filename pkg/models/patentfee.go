package models

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"intellectual_property/pkg/utils"
	"strings"
	"time"
)

// PaymentStatus  缴费状态枚举
type PaymentStatus int

const (
	StatusUnpaid  PaymentStatus = iota //未支付
	StatusPaid                         //已支付
	StatusOverdue                      //已逾期
)

// PatentFee 专利年费表
type PatentFee struct {
	// 基础字段
	ID            int           `json:"id" gorm:"primary_key;type:bigint;AUTO_INCREMENT"` //编号
	PatentID      int           `json:"patent_id" gorm:"type:bigint;"`                    //专利ID
	Patent        Patent        //专利信息
	FeeYear       int           `json:"fee_year" gorm:"fee_year;not null;check:fee_year > 0;comment:缴费年度"` //缴费年度
	Amount        float64       `json:"amount" gorm:"type:decimal(10,2);not null"`                         //年费
	PaymentStatus PaymentStatus ` json:"payment_status" gorm:"payment_status;type:int;default:0"`          // 状态
	// 必要时间字段
	DeadlineDate  time.Time  `gorm:"type:date;not null"` // 截止日
	ActualPayDate *time.Time `gorm:"type:date"`          // 实际缴费日

	// 自动维护时间戳
	CreatedAt time.Time
	UpdatedAt time.Time
}

// FeeStatistics 统计数据结构体
type FeeStatistics struct {
	YearPendingCount int       `json:"year_pending_count"` // 年度待缴费数量
	YearPaidAmount   float64   `json:"year_paid_amount"`   // 年度已缴金额
	OverdueCount     int       `json:"overdue_count"`      // 逾期未缴数量
	TotalAnnualFee   float64   `json:"total_annual_fee"`   // 年度总年费
	CurrentYear      int       `json:"current_year"`       // 统计年度
	LastUpdated      time.Time `json:"last_updated"`       // 最后更新时间
}

// Pagination 新增分页结构体
type Pagination struct {
	Page       int   `json:"page"`        // 当前页码
	PageSize   int   `json:"page_size"`   // 每页数量
	Total      int64 `json:"total"`       // 总记录数
	TotalPages int   `json:"total_pages"` // 总页数
}

// AutoMigrate 数据库迁移
func AutoMigrate() {
	err := utils.DB.AutoMigrate(&PatentFee{})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

// GetFee 根据类型获取专利费用
func GetFee(t int) float64 {
	if t == 0 {
		return 900
	}
	if t == 1 {
		return 600
	}
	if t == 2 {
		return 600

	}
	return 600
}

// NewPatentAnnualFee 新建专利年费
func NewPatentAnnualFee(fee *PatentFee) error {
	return utils.DB.Create(fee).Error
}

// GetPatentFeeByID 根据ID获取专利年费记录
func GetPatentFeeByID(id int) (*PatentFee, error) {
	var fee PatentFee
	result := utils.DB.Preload("Patent").First(&fee, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &fee, nil
}

// GetAllPatentFees 分页获取专利年费记录
func GetAllPatentFees(page, pageSize int) (Pagination, []PatentFee, error) {
	var pagination Pagination
	var fees []PatentFee

	// 参数校验和默认值设置
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	} else if pageSize > 100 {
		pageSize = 100
	}

	// 计算分页偏移量
	offset := (page - 1) * pageSize

	// 获取总记录数
	if err := utils.DB.Model(&PatentFee{}).Count(&pagination.Total).Error; err != nil {
		return pagination, nil, fmt.Errorf("获取总记录数失败: %v", err)
	}

	// 计算总页数
	if pagination.Total > 0 {
		pagination.TotalPages = int((pagination.Total + int64(pageSize) - 1) / int64(pageSize))
	}

	// 执行分页查询
	result := utils.DB.
		Preload("Patent").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC"). // 添加排序保证分页稳定性
		Find(&fees)

	if result.Error != nil {
		return pagination, nil, fmt.Errorf("分页查询失败: %v", result.Error)
	}

	// 设置分页元数据
	pagination.Page = page
	pagination.PageSize = pageSize

	return pagination, fees, nil
}

// UpdatePatentFee 更新专利年费信息
func UpdatePatentFee(fee *PatentFee) error {
	// 使用Select明确更新字段，避免零值问题
	result := utils.DB.Model(fee).Select(
		"payment_status",
		"actual_pay_date",
		"deadline_date",
	).Updates(fee)
	return result.Error
}

// DeletePatentFee 根据ID删除专利年费记录
func DeletePatentFee(id int) error {
	result := utils.DB.Delete(&PatentFee{}, id)
	return result.Error
}

// GetPatentFeesByPatentID 根据专利ID获取相关年费记录
func GetPatentFeesByPatentID(patentID int) ([]PatentFee, error) {
	var fees []PatentFee
	result := utils.DB.Preload("Patent").Where("patent_id = ?", patentID).Find(&fees)
	return fees, result.Error
}

// GetPatentFeesByStatus 根据缴费状态获取记录
func GetPatentFeesByStatus(status PaymentStatus) ([]PatentFee, error) {
	var fees []PatentFee
	result := utils.DB.Preload("Patent").Where("payment_status = ?", status).Find(&fees)
	return fees, result.Error
}

// DeletePatentFeesByPatentID 根据专利ID删除所有相关年费
func DeletePatentFeesByPatentID(patentID int) error {
	result := utils.DB.Where("patent_id = ?", patentID).Delete(&PatentFee{})
	return result.Error
}

/*
1、并发查询优化
使用goroutine并发执行4个统计查询
通过channel统一收集结果
总耗时取决于最慢的查询
*/

// GetFeeStatistics 统计年度费用数据（并发查询）
func GetFeeStatistics() (*FeeStatistics, error) {
	now := time.Now().UTC()
	currentYear := now.Year()

	// 计算精确的年度时间范围（UTC时间）
	firstDayOfYear := time.Date(currentYear, 1, 1, 0, 0, 0, 0, time.UTC)
	lastDayOfYear := time.Date(currentYear, 12, 31, 23, 59, 59, 999999999, time.UTC)
	today := now.Format("2006-01-02")

	var stats FeeStatistics
	stats.CurrentYear = currentYear
	stats.LastUpdated = now

	type result struct {
		value interface{}
		err   error
		name  string
	}

	ch := make(chan result, 4)
	defer close(ch)

	// 并发执行所有统计查询
	go func() {
		var count int64
		err := utils.DB.Model(&PatentFee{}).
			Where("payment_status = ? AND deadline_date >= ?",
				StatusUnpaid, today).
			Count(&count).Error
		ch <- result{value: count, err: err, name: "YearPendingCount"}
	}()

	go func() {
		var amount float64
		err := utils.DB.Model(&PatentFee{}).
			Select("COALESCE(SUM(amount), 0)").
			Where("payment_status = ? AND actual_pay_date BETWEEN ? AND ?",
				StatusPaid, firstDayOfYear, lastDayOfYear).
			Scan(&amount).Error
		ch <- result{value: amount, err: err, name: "YearPaidAmount"}
	}()

	go func() {
		var count int64
		err := utils.DB.Model(&PatentFee{}).
			Where("payment_status = ? AND deadline_date < ?",
				StatusUnpaid, today).
			Count(&count).Error
		ch <- result{value: count, err: err, name: "OverdueCount"}
	}()

	go func() {
		var total float64
		err := utils.DB.Model(&PatentFee{}).
			Select("COALESCE(SUM(amount), 0)").
			Where("fee_year = ?", currentYear).
			Scan(&total).Error
		ch <- result{value: total, err: err, name: "TotalAnnualFee"}
	}()

	// 收集结果
	for i := 0; i < 4; i++ {
		res := <-ch
		if res.err != nil {
			return nil, fmt.Errorf("%s统计失败: %v", res.name, res.err)
		}

		switch res.name {
		case "YearPendingCount":
			stats.YearPendingCount = int(res.value.(int64))
		case "YearPaidAmount":
			stats.YearPaidAmount = res.value.(float64)
		case "OverdueCount":
			stats.OverdueCount = int(res.value.(int64))
		case "TotalAnnualFee":
			stats.TotalAnnualFee = res.value.(float64)
		}
	}

	return &stats, nil
}

// UpdatePatentFeeByApplyNo 根据申请号更新专利费用记录
func UpdatePatentFeeByApplyNo(applyNo string, updateFields map[string]interface{}) error {
	// 字段白名单验证
	allowedFields := map[string]bool{
		"amount":          true,
		"payment_status":  false,
		"deadline_date":   false,
		"actual_pay_date": false,
	}
	for field := range updateFields {
		if !allowedFields[field] {
			return fmt.Errorf("禁止更新字段: %s", field)
		}
	}

	// 开启事务
	tx := utils.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. 根据申请号获取专利信息
	var patent Patent
	if err := tx.Select("id").Where("apply_no = ?", applyNo).First(&patent).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("专利申请号不存在: %s", applyNo)
		}
		return fmt.Errorf("专利查询失败: %v", err)
	}

	// 2. 更新专利费用记录
	result := tx.Model(&PatentFee{}).
		Where("patent_id = ?", patent.Id).
		Updates(updateFields)

	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("费用更新失败: %v", result.Error)
	}

	// 3. 检查实际影响行数
	if result.RowsAffected == 0 {
		tx.Rollback()
		return fmt.Errorf("未找到可更新的费用记录")
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("事务提交失败: %v", err)
	}

	return nil
}

// GetPatentFeesByFilters 根据状态、关键字模糊查询（申请号/专利名称）分页查询
func GetPatentFeesByFilters(
	status *PaymentStatus,
	keyword string,
	page int,
	pageSize int,
) (Pagination, []PatentFee, error) {
	var pagination Pagination
	var fees []PatentFee

	// 参数校验和默认值设置
	if page < 1 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	} else if pageSize > 100 {
		pageSize = 100
	}

	// 创建基础查询
	query := utils.DB.Model(&PatentFee{}).
		Joins("INNER JOIN patents ON patent_fees.patent_id = patents.id")

	// 动态添加条件
	if status != nil {
		query = query.Where("patent_fees.payment_status = ?", *status)
	}
	if keyword != "" {
		escapedKeyword := escapeLike(keyword)
		// 同时查询申请号和专利名称
		query = query.Where(
			"patents.apply_no LIKE ? OR patents.patent_name LIKE ?",
			"%"+escapedKeyword+"%",
			"%"+escapedKeyword+"%",
		)
	}

	// 获取总记录数
	if err := query.Count(&pagination.Total).Error; err != nil {
		return pagination, nil, fmt.Errorf("获取总数失败: %v", err)
	}

	// 计算总页数
	if pagination.Total > 0 {
		pagination.TotalPages = (int(pagination.Total) + pageSize - 1) / pageSize
	}

	// 执行分页查询
	offset := (page - 1) * pageSize
	err := query.
		Preload("Patent").
		Offset(offset).
		Limit(pageSize).
		Order("patent_fees.deadline_date DESC").
		Find(&fees).Error

	if err != nil {
		return pagination, nil, fmt.Errorf("分页查询失败: %v", err)
	}

	// 设置分页元数据
	pagination.Page = page
	pagination.PageSize = pageSize

	return pagination, fees, nil
}

// escapeLike 转义LIKE查询中的特殊字符（保持原样）
func escapeLike(s string) string {
	s = strings.ReplaceAll(s, "%", "\\%")
	s = strings.ReplaceAll(s, "_", "\\_")
	return s
}
