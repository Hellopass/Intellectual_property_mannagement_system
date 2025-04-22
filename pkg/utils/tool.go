package utils

import (
	"strconv"
	"strings"
)

// ConvertStringSliceToInt 将字符串切片安全转换为整型切片
// 参数：
//   - strSlice：需要转换的字符串切片
//
// 返回值：
//   - []int：包含所有有效整数的切片
//
// 功能说明：
//  1. 过滤非数字字符串（如"abc"会跳过）
//  2. 忽略空字符串
//  3. 自动跳过前导/后置空格（如" 123 "转换为123）
//  4. 支持正负数（如"-123"转换为-123）
//  5. 预分配内存提升性能
func ConvertStringSliceToInt(strSlice []string) []int {
	// 预分配结果切片容量
	result := make([]int, 0, len(strSlice))

	for _, s := range strSlice {
		// 去除前后空格
		trimmed := strings.TrimSpace(s)
		if trimmed == "" {
			continue
		}

		// 尝试转换数字
		if num, err := strconv.Atoi(trimmed); err == nil {
			result = append(result, num)
		}
		// 注：这里会忽略所有转换失败的字符串
		// 如果需要严格处理，可以返回错误信息
	}

	return result
}
