package models

import (
	"time"
)

// PatentFee 专利年费模型
type PatentFee struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement;type:bigint"`
	PatentID    int       `json:"patent_id" gorm:"type:bigint;comment:专利ID"`
	Patent      Patent    `json:"patent" gorm:"foreignKey:PatentID"`
	ReviewFee   float64   `json:"review_fee" gorm:"type:float;comment:审核费用"`
	IsPaid      bool      `json:"is_paid" gorm:"comment:是否已支付"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:datetime;comment:创建时间"`
	PaymentDate time.Time `json:"payment_date" gorm:"type:datetime;comment:支付时间"`
	// 新增截至缴费日期字段
	PaymentDeadline time.Time `json:"payment_deadline" gorm:"type:datetime;comment:截至缴费日期"`
	// 新增状态字段，可根据实际情况修改类型和注释
	Status int `json:"status" gorm:"type:int;comment:费用状态"`
}

// GetAllPatentFees 获取所有专利年费
func GetAllPatentFees(keyword string, status int, page int, pageSize int) ([]PatentFee, int64, error) {
	var fees []PatentFee
	var total int64

	query := PatentDB.Preload("Patent")

	// 关键词模糊查询
	if keyword != "" {
		keyword = "%" + keyword + "%"
		query = query.Where(
			"EXISTS ("+
				"SELECT 1 FROM patents "+
				"WHERE patents.id = patent_fees.patent_id "+
				"AND (patents.title LIKE ? OR patents.application_number LIKE ?)"+
				")",
			keyword,
			keyword,
		)
	}

	// 状态过滤
	if status >= 0 {
		query = query.Where("status = ?", status)
	}

	// 计算总数
	if err := query.Model(&PatentFee{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页处理
	if page > 0 && pageSize > 0 {
		query = query.Offset((page - 1) * pageSize).Limit(pageSize)
	}

	// 获取结果
	if err := query.Find(&fees).Error; err != nil {
		return nil, 0, err
	}

	return fees, total, nil
}

// GetMonthlyPatentFeeStats 获取本月专利年费统计信息
func GetMonthlyPatentFeeStats() (int, float64, int, float64, error) {

	var fees []PatentFee
	if err := PatentDB.Find(&fees).Error; err != nil {
		return 0, 0, 0, 0, err
	}

	var pendingCount, overdueCount int
	var paidAmount, totalAmount float64

	for _, fee := range fees {
		totalAmount += fee.ReviewFee
		if fee.Status == 0 {
			pendingCount++
		} else if fee.Status == 1 {
			paidAmount += fee.ReviewFee
		} else if fee.Status == 2 {
			overdueCount++
		}
	}

	return pendingCount, paidAmount, overdueCount, totalAmount, nil
}
