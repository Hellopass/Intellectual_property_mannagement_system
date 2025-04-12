package utils

import "github.com/spf13/viper"

var GinConfig App

type App struct {
	Env     string `mapstructure:"env" json:"env" yaml:"env"`
	Port    string `mapstructure:"port" json:"port" yaml:"port"`
	AppName string `mapstructure:"app_name" json:"app_name" yaml:"app_name"`
	AppUrl  string `mapstructure:"app_url" json:"app_url" yaml:"app_url"`
	JwtKey  string `mapstructure:"jwt_key" json:"jwt_key" yaml:"jwt_key"`
}

// GetGinConfig 读取Gin 配置文件
func GetGinConfig() {
	m := App{}
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		Logger.Error("读取配置错误")
	}
	m.Env = viper.GetString("app.env")
	m.AppUrl = viper.GetString("app.app_url")
	m.Port = viper.GetString("app.port")
	m.AppName = viper.GetString("app.app_name")
	m.JwtKey = viper.GetString("app.jwt_key")
	GinConfig = m
}
