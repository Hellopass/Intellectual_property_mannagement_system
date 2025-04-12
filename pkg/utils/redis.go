package utils

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
	"time"
)

var RDB *redis.Client

type Redis struct {
	Addr         string `json:"addr"`
	Pwd          string `json:"password"`
	DB           int    `json:"db"`
	PoolSize     int    `json:"pool_size"`
	MinIdleConns int    `json:"min_idle_conns"`
}

// redis数据库
var ctx = context.Background()

// 拿到redis配置文件
func getRedisConfig() Redis {
	m := Redis{}
	viper.SetConfigName("redis")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		Logger.Error("读取配置错误")
	}
	m.Addr = viper.GetString("redis.addr")
	m.Pwd = viper.GetString("redis.password")
	m.DB = viper.GetInt("redis.db")
	m.PoolSize = viper.GetInt("redis.pool_size")
	m.MinIdleConns = viper.GetInt("redis.min_idle_conns")

	return m
}

// NewRedisClient 创建一个新的 Redis 客户端
func NewRedisClient() *redis.Client {
	rd := getRedisConfig()
	return redis.NewClient(&redis.Options{
		Addr:         rd.Addr,           // Redis 服务器地址
		Password:     rd.Pwd,            // Redis 服务器密码
		DB:           rd.DB,             // Redis 数据库索引
		PoolSize:     rd.PoolSize,       // 连接池大小
		MinIdleConns: rd.MinIdleConns,   // 最小空闲连接数
		IdleTimeout:  240 * time.Second, // 空闲连接超时时间
	})
}

// RedisSet 添加redis --设置过期时间为5min key为邮箱，value为验证码
func RedisSet(key, value string) error {
	return RDB.Set(ctx, key, value, 5*time.Minute).Err()
}

// RedisGet 获取key -value
func RedisGet(key string) (string, error) {
	value, err := RDB.Get(ctx, key).Result()
	if err != nil {
		Logger.Error(err.Error())
		return "", err
	}
	return value, nil
}

//删除

func RedisDel(key string) error {
	return RDB.Del(ctx, key).Err()
}
