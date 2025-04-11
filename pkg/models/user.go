package models

import (
	"intellectual_property/pkg/utils"
)

// User 用户信息表
type User struct {
	ID         uint   `json:"id" gorm:"id"`                   //编号
	UserName   string `json:"user_name" gorm:"user_name"`     //姓名
	DepID      int    `json:"dep_id" gorm:"dep_id"`           //单位ID
	Password   string `json:"password" gorm:"password"`       //登录密码
	Authority  string `json:"authority" gorm:"authority"`     //权限控制
	Sex        string `json:"sex" gorm:"sex"`                 //性别
	Birth      string `json:"birth" gorm:"birth"`             //出生日期
	IDCard     string `json:"id_card" gorm:"id_card"`         //身份证号码
	Political  string `json:"political" gorm:"political"`     //政治面貌
	Unit       string `json:"unit" gorm:"unit"`               //所属学院
	LastDegree string `json:"last_degree" gorm:"last_degree"` //最高学历
	TechIP     string `json:"tech_ip" gorm:"teach_ip"`        //职称
	Cour       string `json:"cour" gorm:"cour"`               //一级学科
	Research   string `json:"research" gorm:"research"`       //研究方向
}

// CreateUser 创建新用户
func CreateUser(user *User) error {
	return utils.DB.Create(user).Error
}

// GetUserByID 根据用户ID获取用户信息
func GetUserByID(userID int) (*User, error) {
	var user User
	if err := utils.DB.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser 更新用户信息
func UpdateUser(user *User) error {
	return utils.DB.Save(user).Error
}

// DeleteUser 删除用户
func DeleteUser(user *User) error {
	return utils.DB.Delete(user).Error
}
