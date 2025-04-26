package models

import (
	"intellectual_property/pkg/utils"
	"time"
)

// StatisticsResult 统计结果表
type StatisticsResult struct {
	PatentCount          int64     `json:"patent_count"`
	ArticleCount         int64     `json:"article_count"`
	TrademarkCount       int64     `json:"trademark_count"`
	PendingApproval      int64     `json:"pending_approval"`
	PatentTrendLastSixMonths []int64 `json:"patent_trend_last_six_months"`
}

var StatisticsDB = utils.DB

// GetStatistics 执行统计操作
func GetStatistics() (*StatisticsResult, error) {
	var result StatisticsResult

	// 统计专利总数
	if err := StatisticsDB.Model(&Patent{}).Count(&result.PatentCount).Error; err != nil {
		return nil, err
	}

	// 统计著作总数
	if err := StatisticsDB.Model(&Article{}).Count(&result.ArticleCount).Error; err != nil {
		return nil, err
	}

	// 统计商标总数
	if err := StatisticsDB.Model(&Trademark{}).Count(&result.TrademarkCount).Error; err != nil {
		return nil, err
	}

	// 统计审批中状态的待审批数量
	// 假设专利、著作和商标都有 ApprovalStatus 字段，0 表示审批中
	var pendingPatent int64
	if err := StatisticsDB.Model(&Patent{}).Where("approval_status = ?", 0).Count(&pendingPatent).Error; err != nil {
		return nil, err
	}
	var pendingArticle int64
	if err := StatisticsDB.Model(&Article{}).Where("approval_status = ?", 0).Count(&pendingArticle).Error; err != nil {
		return nil, err
	}
	var pendingTrademark int64
	if err := StatisticsDB.Model(&Trademark{}).Where("approval_status = ?", 0).Count(&pendingTrademark).Error; err != nil {
		return nil, err
	}

	result.PendingApproval = pendingPatent + pendingArticle + pendingTrademark

	// 统计专利近半年申请趋势
	result.PatentTrendLastSixMonths = make([]int64, 6)
	now := time.Now()
	for i := 0; i < 6; i++ {
		start := time.Date(now.Year(), now.Month()-time.Month(i), 1, 0, 0, 0, 0, now.Location())
		end := start.AddDate(0, 1, 0)
		var count int64
		if err := StatisticsDB.Model(&Patent{}).Where("apply_date >= ? AND apply_date < ?", start, end).Count(&count).Error; err != nil {
			return nil, err
		}
		result.PatentTrendLastSixMonths[5-i] = count
	}

	return &result, nil
}
