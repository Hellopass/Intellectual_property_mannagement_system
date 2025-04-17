package utils

import "github.com/spf13/viper"

var NgX Nginx

type Nginx struct {
	LocationAvatar string `json:"location_avatar"`
	LocationDocs   string `json:"location_docs"`
	Url            string `json:"url"`
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
	m.LocationDocs = viper.GetString("nginx.location_docs")
	return m
}
