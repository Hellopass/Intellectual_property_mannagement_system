package models

import (
	"errors"
	"strconv"
)

// TrademarkType 商标类型枚举
// 定义系统支持的商标类型代码
type TrademarkType int

const (
	GoodsTrademark         TrademarkType = iota + 1 // 商品商标
	ServiceTrademark                                // 服务商标
	CollectiveTrademark                             // 集体商标
	CertificationTrademark                          // 证明商标
)

// trademarkTypeMap 类型映射字典
// 实现中文类型名称到枚举值的映射，用于前端展示和输入解析
var trademarkTypeMap = map[string]TrademarkType{
	"商品商标": GoodsTrademark,
	"服务商标": ServiceTrademark,
	"集体商标": CollectiveTrademark,
	"证明商标": CertificationTrademark,
}

// ParseTrademarkType 类型解析函数
// 参数：keyword 可以是中文类型名称或数字字符串（0-5）
// 返回值：(类型代码, error)
// 成功时返回对应类型代码，失败返回-1和错误信息
// 使用场景：处理用户输入的类型数据转换
func ParseTrademarkType(keyword string) (int, error) {
	// 尝试中文名称匹配
	if t, exists := trademarkTypeMap[keyword]; exists {
		return int(t), nil
	}

	// 尝试数字转换
	if n, err := strconv.Atoi(keyword); err == nil {
		if n >= 0 && n <= 5 {
			return n, nil
		}
		return -1, errors.New("类型代码超出有效范围（0-5）")
	}

	// 无效输入处理
	return -1, errors.New("无法识别的商标类型：" + keyword)
}
