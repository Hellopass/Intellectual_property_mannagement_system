package models

import (
	"intellectual_property/pkg/utils"
	"time"
)

type User struct {
	ID           int       `json:"id" gorm:"id"`
	UserType     int       `json:"user_type" gorm:"user_type"`
	UserName     string    `json:"user_name" gorm:"user_name"`
	PasswordHash string    `json:"password_hash" gorm:"password_hash"`
	Register     time.Time `json:"register" gorm:"register"` //注册时间
	Email        string    `json:"email" gorm:"email"`
	Location     string    `json:"location" gorm:"location"` //地域统计
}

// CreateUserTable 迁移User
func CreateUserTable() {
	utils.DB.AutoMigrate(&User{})
}
