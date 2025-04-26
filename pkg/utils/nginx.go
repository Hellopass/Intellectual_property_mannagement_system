package utils

import "github.com/spf13/viper"

var NgX Nginx

type Nginx struct {
	LocationAvatar    string `json:"location_avatar"`
	LocationPatent    string `json:"location_patent"`
	LocationArticle   string `json:"location_article"`
	LocationTrademark string `json:"location_trademark"`
	Url               string `json:"url"`
}

// getNginxConfig 读取Nginx配置文件
func getNginxConfig() Nginx {
	m := Nginx{}
	viper.SetConfigName("nginx")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		Logger.Error("读取配置错误")
	}
	m.LocationAvatar = viper.GetString("nginx.location_avatar")
	m.Url = viper.GetString("nginx.url")
	m.LocationPatent = viper.GetString("nginx.location_patent")
	m.LocationTrademark = viper.GetString("nginx.location_trademark")
	m.LocationArticle = viper.GetString("nginx.location_article")
	return m
}
