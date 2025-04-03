package utils

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Mysql struct {
	User   string `json:"user"`
	Pass   string `json:"pass"`
	Ip     string `json:"ip"`
	Port   string `json:"port"`
	DbName string `json:"dbname"`
}

// 读取mysql 配置文件
func getMysqlConfig() Mysql {
	m := Mysql{}
	viper.SetConfigName("mysql")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		Logger.Error("读取配置错误")
	}
	m.User = viper.GetString("mysql.user")
	m.Pass = viper.GetString("mysql.pass")
	m.Port = viper.GetString("mysql.port")
	m.Ip = viper.GetString("mysql.ip")
	m.DbName = viper.GetString("mysql.dbname")
	return m
}

// 获取dsn
func getDns() string {
	m := getMysqlConfig()
	dsn := m.User + ":" + m.Pass + "@tcp(" + m.Ip + ":" + m.Port + ")/" + m.DbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	return dsn
}

// ContactMysql 连接到mysql
func ContactMysql() error {
	dsn := getDns()
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		Logger.Error("连接Mysql错误")
		return err
	}
	DB = db
	return err
}
