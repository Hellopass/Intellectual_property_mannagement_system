package utils

import (
	"fmt"
	"go.uber.org/zap"
)

//初始化需要的工具类

var Logger *zap.Logger //最后记得 	defer Logger.Sync()

func init() {
	lc := LogConfig{
		Level:        "debug",
		FileName:     "./log/test.log",
		MaxSize:      1,
		MaxBackups:   5,
		MaxAge:       30,
		IsStackTrace: true,
		IsStdout:     true,
	}
	err := InitLogger(lc)
	if err != nil {
		fmt.Println(err)
	}
	// L()：获取全局logger
	Logger = zap.L()

	//初始化mysql
	err = ContactMysql()
	if err != nil {
		Logger.Error("初始化mysql失败")
	}

}
