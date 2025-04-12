package api

import (
	"intellectual_property/pkg/utils"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger = utils.Logger

// NewEngine 创建预配置的 Gin 引擎
// envType: "debug" 或 "release"
func NewEngine(envType string) *gin.Engine {
	// 初始化日志
	logger := initZapLogger(envType)
	defer logger.Sync()

	// 配置 Gin 模式
	if envType == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 创建引擎
	engine := gin.New()

	// 添加日志中间件
	engine.Use(ginZapLogger(logger))   // 替换默认 Logger
	engine.Use(ginZapRecovery(logger)) // 替换默认 Recovery

	// 添加通用中间件
	engine.Use(gin.Recovery())
	//中间
	engine.Use(utils.Cors()) //解决跨域问题
	return engine
}

// 初始化 Zap 日志（根据环境配置）
func initZapLogger(env string) *zap.Logger {
	isDev := env != "release"

	// 日志编码配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 多环境输出配置
	var cores []zapcore.Core

	// 开发环境：控制台输出
	if isDev {
		cores = append(cores, zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.Lock(os.Stdout),
			zapcore.DebugLevel,
		))
	}

	// 生产环境：文件输出（带切割）
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "log/gin.log",
		MaxSize:    100, // MB
		MaxBackups: 30,
		MaxAge:     7, // days
		Compress:   true,
	})

	cores = append(cores, zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		fileWriter,
		zapcore.InfoLevel,
	))

	// 创建核心组合
	core := zapcore.NewTee(cores...)

	return zap.New(core, zap.AddCaller())
}

// Gin 日志中间件（使用 Zap）
func ginZapLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 记录日志
		logger.Info("HTTP Request",
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.Duration("latency", time.Since(start)),
			zap.Strings("errors", c.Errors.Errors()),
		)
	}
}

// Gin 错误恢复中间件（使用 Zap）
func ginZapRecovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("HTTP Panic",
					zap.Any("error", err),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.String("client_ip", c.ClientIP()),
					zap.Stack("stack"),
				)

				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
