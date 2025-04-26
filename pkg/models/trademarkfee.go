package models

import (
	"time"
)

// TrademarkFee 商标年费模型
type TrademarkFee struct {
	ID          int       `json:"id" gorm:"primaryKey;autoIncrement;type:bigint"`
	TrademarkID int       `json:"trademark_id" gorm:"type:bigint;comment:商标ID"`
	Trademark   Trademark `json:"trademark" gorm:"foreignKey:TrademarkID"`
	ReviewFee   float64   `json:"review_fee" gorm:"type:float;comment:审核费用"`
	IsPaid      bool      `json:"is_paid" gorm:"comment:是否已支付"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:datetime;comment:创建时间"`
	PaymentDate time.Time `json:"payment_date" gorm:"type:datetime;comment:支付时间"`
	// 新增截至缴费日期字段
	PaymentDeadline time.Time `json:"payment_deadline" gorm:"type:datetime;comment:截至缴费日期"`
	// 新增状态字段，可根据实际情况修改类型和注释
	Status int `json:"status" gorm:"type:int;comment:费用状态"`
}

// GetAllTrademarkFees 获取所有商标年费
func GetAllTrademarkFees(keyword string, status int, page int, pageSize int) ([]TrademarkFee, int64, error) {
	var fees []TrademarkFee
	var total int64

	query := TrademarkDB.Preload("Trademark")

	// 关键词模糊查询
	if keyword != "" {
		keyword = "%" + keyword + "%"
		query = query.Where(
			"EXISTS ("+
				"SELECT 1 FROM trademarks "+
				"WHERE trademarks.id = trademark_fees.trademark_id "+
				"AND (trademarks.title LIKE ? OR trademarks.application_number LIKE ?)"+
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
	if err := query.Model(&TrademarkFee{}).Count(&total).Error; err != nil {
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

// GetMonthlyTrademarkFeeStats 获取本月商标年费统计信息
func GetMonthlyTrademarkFeeStats() (int, float64, int, float64, error) {

	var fees []TrademarkFee
	if err := TrademarkDB.Debug().Find(&fees).Error; err != nil {
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