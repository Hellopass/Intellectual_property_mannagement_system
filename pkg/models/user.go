package models

import (
	"intellectual_property/pkg/utils"
)

// User 用户信息表
type User struct {
	ID         int    `json:"id" gorm:"primaryKey;autoIncrement;type:bigint"` //编号
	UserName   string `json:"user_name" gorm:"user_name"`                     //姓名
	DepID      int    `json:"dep_id" gorm:"dep_id"`                           //单位ID
	Password   string `json:"password" gorm:"password"`                       //登录密码
	Email      string `json:"email" gorm:"email"`                             //邮箱
	Authority  string `json:"authority" gorm:"authority"`                     //权限控制 admin/user--管理员和普通用户
	Sex        string `json:"sex" gorm:"sex"`                                 //性别
	Birth      string `json:"birth" gorm:"birth"`                             //出生日期
	IDCard     string `json:"id_card" gorm:"id_card"`                         //身份证号码
	Political  string `json:"political" gorm:"political"`                     //政治面貌
	Unit       string `json:"unit" gorm:"unit"`                               //所属学院
	LastDegree string `json:"last_degree" gorm:"last_degree"`                 //最高学历
	TechIP     string `json:"tech_ip" gorm:"tech_ip"`                         //职称
	Cour       string `json:"cour" gorm:"cour"`                               //一级学科
	Research   string `json:"research" gorm:"research"`                       //研究方向
	Status     string `json:"status" gorm:"status"`                           //用户状态   0-表示用户以注销，1-表示用户正常
	AvatarUrl  string `json:"avatar_url" gorm:"avatar_url"`                   //用户头像
}

// SimpleUser 如果只需要特定字段，可以定义轻量结构体
type SimpleUser struct {
	ID        int    `json:"id"`
	UserName  string `json:"user_name"`
	Authority string `json:"authority"`
}

// GetAllSimpleUsers 获取部分user
func GetAllSimpleUsers() ([]SimpleUser, error) {
	var users []SimpleUser
	err := utils.DB.Model(&User{}).Select("id, user_name,authority").Find(&users).Error
	return users, err
}

// GetSimpleUserByID 根据用户ID获取部分用户信息
func GetSimpleUserByID(userID int) (SimpleUser, error) {
	var user SimpleUser
	if err := utils.DB.Model(&User{}).Select("id, user_name,authority").Find(&user, userID).Error; err != nil {
		return SimpleUser{}, err
	}
	return user, nil
}

// CreateUser 创建新用户
func CreateUser(user *User) error {
	return utils.DB.Create(user).Error
}

// GetUserByID 根据用户ID获取用户信息
func GetUserByID(userID int) (User, error) {
	var user User
	if err := utils.DB.Find(&user, userID).Error; err != nil {
		return User{}, err
	}
	return user, nil
}

// UpdateUser 更新用户信息
func UpdateUser(user *User) error {
	return utils.DB.Model(&User{}).Where("id=?", user.ID).Select("dep_id", "political", "unit", "last_degree", "tech_ip", "cour", "research").Updates(&user).Error
}

// DeleteUser 删除用户
func DeleteUser(user *User) error {
	return utils.DB.Delete(user).Error
}

// GetUserByEmail 根据电子邮件获取用户信息
func GetUserByEmail(email string) (*User, error) {
	var user User
	if err := utils.DB.Where("email = ?", email).Find(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UploadAvatar 上传头像地址
func UploadAvatar(id int, url string) error {
	return utils.DB.Model(&User{}).Where("id = ?", id).Update("avatar_url", url).Error
}
