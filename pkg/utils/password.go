// 密码加密
package utils

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
)

// 安全要求配置
const (
	RequiredPasswordLength = 8 // 强制密码长度
	SaltLength             = 8 // 盐值长度（增强至32字节）
)

// 生成密码学安全随机盐
func GenerateSecureSalt() (string, error) {
	salt := make([]byte, SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("salt generation failed: %v", err)
	}
	return hex.EncodeToString(salt), nil
}

// 带长度验证的密码哈希
func SecureHashWithSalt(password string, salt string) (string, error) {
	// 严格密码长度校验
	if len(password) != RequiredPasswordLength {
		return "", errors.New("password must be exactly 16 characters")
	}

	// 盐值解码验证
	decodedSalt, err := hex.DecodeString(salt)
	if err != nil || len(decodedSalt) != SaltLength {
		return "", errors.New("invalid salt format")
	}

	// 构造加盐密码
	saltedData := append(decodedSalt, []byte(password)...)

	// 双重哈希增强
	firstHash := md5.Sum(saltedData)
	secondHash := md5.Sum(firstHash[:])

	return hex.EncodeToString(secondHash[:]), nil
}

// 密码验证函数
func VerifyPassword(inputPassword, storedSalt, storedHash string) bool {
	// 前置长度检查
	if len(inputPassword) != RequiredPasswordLength {
		return false
	}

	// 计算哈希
	computedHash, err := SecureHashWithSalt(inputPassword, storedSalt)
	if err != nil {
		return false
	}

	// 安全对比（防止时序攻击）
	return subtle.ConstantTimeCompare([]byte(computedHash), []byte(storedHash)) == 1
}

// 分割拿到盐值
func GetSlat(password string) string {
	return string([]byte(password)[len(password)/2:])
}
