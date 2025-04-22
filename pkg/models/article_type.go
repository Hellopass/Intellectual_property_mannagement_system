// models/article_type.go
package models

import (
	"errors"
	"strconv"
)

// ArticleType 文献类型枚举
// 定义系统支持的六种文献类型代码
type ArticleType int

const (
	Books               ArticleType = iota // 书籍（代码0）
	JournalPapers                          // 期刊论文（代码1）
	ConferencePapers                       // 会议论文（代码2）
	DegreePapers                           // 学位论文（代码3）
	TechnologyStandards                    // 技术标准（代码4）
	WebResources                           // 网页资源（代码5）
)

// articleTypeMap 类型映射字典
// 实现中文类型名称到枚举值的映射，用于前端展示和输入解析
var articleTypeMap = map[string]ArticleType{
	"书籍":   Books,
	"期刊论文": JournalPapers,
	"会议论文": ConferencePapers,
	"学位论文": DegreePapers,
	"技术标准": TechnologyStandards,
	"网页资源": WebResources,
}

// ParseArticleType 类型解析函数
// 参数：keyword 可以是中文类型名称或数字字符串（0-5）
// 返回值：(类型代码, error)
// 成功时返回对应类型代码，失败返回-1和错误信息
// 使用场景：处理用户输入的类型数据转换
func ParseArticleType(keyword string) (int, error) {
	// 尝试中文名称匹配
	if t, exists := articleTypeMap[keyword]; exists {
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
	return -1, errors.New("无法识别的文献类型：" + keyword)
}
