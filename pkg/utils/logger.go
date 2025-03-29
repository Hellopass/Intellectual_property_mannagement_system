/*
                  日志配置包
使用zap包 ---地址 go.uber.org/zap

*/

package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

type LogConfig struct {
	Level        string `json:"level"`          //最低日志等级
	FileName     string `json:"file_name"`      //日志文件
	MaxSize      int    `json:"max_size"`       //切割之前日志文件最大大小
	MaxAge       int    `json:"max_age"`        //日志文件保存时间
	MaxBackups   int    `json:"max_backups"`    //保留旧日志文件的最大数量
	IsStdout     bool   `json:"is_stdout"`      //是否输出控制台
	IsStackTrace bool   `json:"is_stack_trace"` //是否输出堆栈信息
}

// InitLogger 初始化Logger
func InitLogger(lCfg LogConfig) (err error) {
	writeSyncer := getLogWriter(lCfg.FileName, lCfg.MaxSize, lCfg.MaxBackups, lCfg.MaxAge, lCfg.IsStdout)
	encoder := getEncoder()
	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(lCfg.Level))
	if err != nil {
		return
	}
	core := zapcore.NewCore(encoder, writeSyncer, l)
	var logger *zap.Logger
	if lCfg.IsStackTrace {
		logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
	} else {
		logger = zap.New(core, zap.AddCaller())
	}
	zap.ReplaceGlobals(logger)
	return
}

// 负责设置encoding 的日志文件
func getEncoder() zapcore.Encoder {
	encodeConfig := zap.NewProductionEncoderConfig()
	encodeConfig.EncodeTime = TimeEncoder
	encodeConfig.TimeKey = "time"
	encodeConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encodeConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encodeConfig)
}
func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

// 负责日志写入的位置
func getLogWriter(filename string, maxsize, maxBackup, maxAge int, isStdout bool) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,  // 文件位置
		MaxSize:    maxsize,   // 进行切割之前,日志文件的最大大小(MB为单位)
		MaxAge:     maxAge,    // 保留旧文件的最大天数
		MaxBackups: maxBackup, // 保留旧文件的最大个数
		Compress:   false,     // 是否压缩/归档旧文件
	}
	if isStdout {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(lumberJackLogger), zapcore.AddSync(os.Stdout))
	} else {
		return zapcore.AddSync(lumberJackLogger)
	}
}
