package utils

import (
	"errors"
	"strconv"
	"time"
)

/*
第一部分：中国国家代码为CN；
第二部分：前4位数字为专利申请年份；
第三部分：第5位数字代表类型：数字1表示发明专利申请；数字2表示实用新型专利申请；数字3表示外观设计专利申请；数字8表示进入中国国家阶段的PCT发明专利申请；数字9表示进入中国国家阶段的PCT 实用新型专利申请。
第四部分：第6~12位是流水号；
第五部分：校验位；校验位是指以专利申请号中使用的数字组合作为源数据经过计算得出的1位阿拉伯数字（0至9）或大写英文字母X。
*/

// calculateCheckDigit 计算校验位
func calculateCheckDigit(patentID string) string {
	sum := 0
	for i, char := range patentID {
		num, _ := strconv.Atoi(string(char))
		weight := i + 1
		sum += num * weight
	}

	remainder := sum % 11
	if remainder == 10 {
		return "X"
	}
	return strconv.Itoa(remainder)
}

// generateSerialNumber 生成基于时间戳的7位流水号
func generateSerialNumber() string {
	timestamp := time.Now().UnixNano() // 获取纳秒级时间戳
	serialNumber := strconv.FormatInt(timestamp, 10)[len(strconv.FormatInt(timestamp, 10))-7:]
	return serialNumber
}

// GenerateApplyNo 生成专利申请号
func GenerateApplyNo(countryCode string, year string, patentType string) (string, error) {
	serialNumber := generateSerialNumber()
	if len(year) != 4 || len(patentType) != 1 || len(serialNumber) != 7 {
		return "", errors.New("输入的参数长度不正确")
	}

	patentID := countryCode + year + patentType + serialNumber
	checkDigit := calculateCheckDigit(patentID)
	return patentID + checkDigit, nil
}
