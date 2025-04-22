package models

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"intellectual_property/pkg/utils"
	"math"
	"sort"
	"strings"
	"time"
)

// 科技数据集
var domainKeywords = map[string][]string{
	// 基础软件技术
	"基础软件": {
		"操作系统", "数据库", "中间件", "编译器", "虚拟机",
		"文件系统", "分布式存储", "容错机制", "事务处理",
		"查询优化", "内存管理", "安全沙箱",
	},

	// 云计算与分布式系统
	"云计算": {
		"容器编排", "微服务", "服务网格", "无服务器架构",
		"云原生", "Kubernetes", "服务发现", "弹性伸缩",
		"多租户隔离", "混合云", "边缘计算",
	},

	// 前端开发技术
	"前端开发": {
		"React", "Vue", "TypeScript", "WebAssembly",
		"响应式设计", "SPA", "PWA", "前端性能优化",
		"跨平台框架", "微前端", "低代码可视化",
	},

	// 后端系统架构
	"后端架构": {
		"高并发", "分布式锁", "消息队列", "API网关",
		"服务熔断", "负载均衡", "缓存穿透", "读写分离",
		"分库分表", "CQRS", "事件溯源",
	},

	// 开发工具链
	"开发工具": {
		"IDE插件", "CI/CD", "静态分析", "代码审查",
		"版本控制", "调试工具", "性能剖析", "依赖管理",
		"自动化测试", "混沌工程", "监控告警",
	},

	// 数据科学与大数据
	"数据处理": {
		"Spark", "Flink", "实时计算", "数据湖",
		"特征工程", "OLAP", "数据治理", "隐私计算",
		"流批一体", "数据血缘", "质量监控",
	},

	// 网络安全技术
	"软件安全": {
		"漏洞扫描", "渗透测试", "加密算法", "零信任",
		"代码审计", "WAF", "沙箱隔离", "证书管理",
		"密钥轮换", "访问控制", "日志溯源",
	},

	// 人工智能工程化
	"AI工程化": {
		"模型部署", "MLOps", "特征存储", "在线推理",
		"A/B测试", "模型监控", "自动化标注", "联邦学习",
		"模型压缩", "知识蒸馏", "边缘推理",
	},

	// 新兴技术方向
	"前沿技术": {
		"Web3.0", "低代码平台", "无代码开发", "RPA",
		"Serverless", "量子编程", "AI代码生成",
		"数字孪生", "元宇宙开发", "区块链智能合约",
	},

	// 行业解决方案
	"行业软件": {
		"ERP", "CRM", "SCM", "医疗信息化",
		"金融核心系统", "工业软件", "GIS系统",
		"CAD/CAE", "EDA工具链", "数字孪生平台",
	},
	// 半导体与芯片技术
	"半导体": {
		"芯片", "半导体", "集成电路", "硅基", "晶圆",
		"光刻胶", "封装测试", "第三代半导体", "功率器件",
		"存储芯片", "射频芯片", "MEMS",
	},

	// 新能源与储能技术
	"新能源": {
		"电池", "太阳能", "光伏", "储能", "锂电池",
		"钠离子电池", "氢能源", "燃料电池", "超级电容",
		"风能", "智能电网", "能源互联网",
	},

	// 人工智能与算法
	"人工智能": {
		"AI", "人工智能", "机器学习", "深度学习", "神经网络",
		"自然语言处理", "计算机视觉", "强化学习", "知识图谱",
		"边缘AI", "联邦学习", "大模型",
	},

	// 生物医药与基因技术
	"生物医药": {
		"基因", "蛋白", "疫苗", "试剂", "药物",
		"细胞治疗", "基因编辑", "抗体药物", "生物标记物",
		"mRNA", "合成生物学", "器官芯片",
	},

	// 区块链与分布式技术
	"区块链": {
		"区块链", "分布式账本", "智能合约", "加密货币",
		"DeFi", "NFT", "跨链协议", "零知识证明",
		"共识算法", "去中心化身份", "链上治理",
	},

	// 5G/6G通信技术
	"通信技术": {
		"5G", "毫米波", "Massive MIMO", "网络切片",
		"边缘计算", "太赫兹通信", "卫星互联网", "空天地一体化",
		"URLLC", "O-RAN", "6G",
	},

	// 量子信息科技
	"量子科技": {
		"量子计算", "量子比特", "量子纠缠", "量子加密",
		"量子算法", "量子模拟", "量子传感", "量子通信",
		"超导量子", "离子阱", "拓扑量子",
	},

	// 机器人技术
	"机器人": {
		"工业机器人", "协作机器人", "SLAM", "运动控制",
		"柔性抓取", "人机交互", "仿生机器人", "无人机",
		"自主导航", "力觉反馈", "群体智能",
	},

	// 先进制造技术
	"先进制造": {
		"3D打印", "数控机床", "数字孪生", "工业互联网",
		"精密加工", "智能工厂", "预测性维护", "柔性生产",
		"复合材料", "无损检测", "工艺优化",
	},

	// 元宇宙与虚拟现实
	"元宇宙": {
		"虚拟现实", "增强现实", "数字孪生", "虚拟化身",
		"空间计算", "脑机接口", "NFT", "沉浸式交互",
		"光场显示", "虚实融合", "Web3",
	},

	// 更多领域...
	"自动驾驶": {
		"激光雷达", "多模态融合", "高精地图", "V2X",
		"路径规划", "仿真测试", "车路协同", "影子模式",
		"端到端学习", "传感器标定",
	},

	"绿色科技": {
		"碳捕捉", "生物降解", "循环经济", "清洁能源",
		"碳足迹", "生态修复", "可持续材料", "零碳建筑",
		"蓝碳", "气候模型",
	},

	"航空航天": {
		"可重复火箭", "卫星星座", "高温合金", "电推进",
		"空间站", "高超声速", "复合材料", "在轨服务",
		"月球基地", "深空探测",
	},

	"网络安全": {
		"零信任", "威胁检测", "数据加密", "APT防御",
		"隐私计算", "安全多方计算", "漏洞挖掘", "攻防演练",
		"安全运营", "区块链审计",
	},

	"智慧城市": {
		"城市大脑", "智能交通", "智慧灯杆", "数字政务",
		"地下管网", "应急指挥", "社区治理", "智慧园区",
		"环境监测", "一网统管",
	},
}

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

// YearlyTrend 定义分析结果结构体
type YearlyTrend struct {
	Year  int `json:"year"`
	Count int `json:"count"`
}

type TypeDistribution struct {
	Type       string  `json:"type"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

type TopApplicant struct {
	Applicant string `json:"applicant"`
	Count     int    `json:"count"`
}

type TechDomain struct {
	Domain string `json:"domain"`
	Count  int    `json:"count"`
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

// GetYearlyTrends 获取近五年专利申请趋势（改进版）
func GetYearlyTrends() ([]YearlyTrend, error) {
	currentYear := time.Now().Year()
	startYear := currentYear - 4

	// 生成固定年份序列（包含最近五年）
	years := make([]int, 5)
	for i := 0; i < 5; i++ {
		years[i] = startYear + i
	}

	// 构建带年份生成的SQL查询
	query := `
	WITH year_series AS (
		SELECT ? AS year UNION ALL
		SELECT ? UNION ALL
		SELECT ? UNION ALL
		SELECT ? UNION ALL
		SELECT ?
	)
	SELECT 
		ys.year,
		COUNT(p.id) AS count
	FROM year_series ys
	LEFT JOIN patents p 
		ON YEAR(p.apply_data) = ys.year
	GROUP BY ys.year
	ORDER BY ys.year`

	var trends []struct {
		Year  int `gorm:"column:year"`
		Count int `gorm:"column:count"`
	}

	// 执行参数化查询（防止SQL注入）
	err := utils.DB.Raw(query, years[0], years[1], years[2], years[3], years[4]).Scan(&trends).Error
	if err != nil {
		return nil, fmt.Errorf("年度趋势查询失败: %v", err)
	}

	// 构建完整结果集（保证五年数据完整性）
	resultMap := make(map[int]int)
	for _, t := range trends {
		resultMap[t.Year] = t.Count
	}

	// 按顺序填充结果
	result := make([]YearlyTrend, 5)
	for i, y := range years {
		result[i] = YearlyTrend{
			Year:  y,
			Count: resultMap[y],
		}
	}

	return result, nil
}

// GetTypeDistribution 获取专利类型分布（修正字段和类型处理）
func GetTypeDistribution() ([]TypeDistribution, error) {
	var typeCounts []struct {
		Type  int `gorm:"column:patent_type"`
		Count int
	}

	// 修正字段名为 patent_type
	if err := utils.DB.Model(&Patent{}).
		Select("patent_type, COUNT(*) AS count").
		Group("patent_type").
		Scan(&typeCounts).Error; err != nil {
		return nil, fmt.Errorf("查询类型分布失败: %v", err)
	}

	total := 0
	for _, tc := range typeCounts {
		total += tc.Count
	}

	result := make([]TypeDistribution, 0, len(typeCounts))
	for _, tc := range typeCounts {
		percent := 0.0
		if total > 0 {
			percent = math.Round(float64(tc.Count)/float64(total)*100*100) / 100
		}

		result = append(result, TypeDistribution{
			Type:       getPatentTypeName(tc.Type),
			Count:      tc.Count,
			Percentage: percent,
		})
	}

	return result, nil
}

// GetTopApplicants 获取前10申请人（关联用户表查询）
func GetTopApplicants() ([]TopApplicant, error) {
	var applicants []TopApplicant

	// 关联用户表查询申请人姓名
	err := utils.DB.Model(&Patent{}).
		Select("users.user_name AS applicant, COUNT(*) AS count").
		Joins("LEFT JOIN users ON patents.user_id = users.id").
		Where("users.user_name IS NOT NULL AND users.user_name <> ''").
		Group("users.user_name").
		Order("count DESC").
		Limit(10).
		Scan(&applicants).Error

	if err != nil {
		return nil, fmt.Errorf("查询申请人统计失败: %v", err)
	}

	return applicants, nil
}

// GetTechDomains 改进的领域分析算法
func GetTechDomains() ([]TechDomain, error) {
	var techChan = make(chan map[string]int, 1)
	go func() {
		var patents []struct{ PatentName string }
		utils.DB.Model(&Patent{}).Select("patent_name").Find(&patents)

		countMap := make(map[string]int)

		for _, p := range patents {
			name := strings.ToLower(p.PatentName)
			matched := make(map[string]struct{})
			for domain, keys := range domainKeywords {
				for _, kw := range keys {
					if strings.Contains(name, strings.ToLower(kw)) {

						if _, exists := matched[domain]; !exists {
							countMap[domain]++
							matched[domain] = struct{}{}
						}
						break
					}
				}
			}
		}

		techChan <- countMap
	}()

	select {
	case result := <-techChan:
		domains := make([]TechDomain, 0, len(result))
		for k, v := range result {
			domains = append(domains, TechDomain{Domain: k, Count: v})
		}
		sort.Slice(domains, func(i, j int) bool {
			return domains[i].Count > domains[j].Count
		})
		return domains, nil
	case <-time.After(5 * time.Second):
		return nil, errors.New("技术领域分析超时")
	}
}

// getPatentTypeName 获取专利类型名称
func getPatentTypeName(t int) string {
	types := []string{"发明专利", "实用新型", "外观设计"}
	if t >= 0 && t < len(types) {
		return types[t]
	}
	return "其他"
}
