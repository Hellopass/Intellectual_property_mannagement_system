package utils

import (
	"crypto/rand"
	"fmt"
	"github.com/jordan-wright/email"
	"github.com/spf13/viper"
	"math/big"
	"net/smtp"
)

// 发送邮件

type Email struct {
	MyEmail  string `json:"my_email"`
	Password string `json:"password"`
	Url      string `json:"url"`
}

// getEmailConfig 读取Email配置文件
func getEmailConfig() Email {
	m := Email{}
	viper.SetConfigName("email")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		Logger.Error("读取配置错误")
	}
	m.MyEmail = viper.GetString("email.my_email")
	m.Password = viper.GetString("email.password")
	m.Url = viper.GetString("email.url")
	return m
}

/*   邮件发送模板
【知识产权】您此次验证码为1234，5分钟内有效，请您尽快验证！
*/

// generateSecureCode 并发安全的验证码生成
func generateSecureCode() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(10000))
	return fmt.Sprintf("%04d", n)
}

// SendAddUserEmailCode 发送邮件
func SendAddUserEmailCode(toEmail string) error {
	config := getEmailConfig()
	e := email.NewEmail()
	//设置发送方的邮箱
	e.From = "zh <" + config.MyEmail + ">"
	// 设置接收方的邮箱
	e.To = []string{toEmail}
	//设置主题
	e.Subject = "注册验证码"
	//设置code
	code := generateSecureCode()
	err := RedisSet(toEmail, code)
	if err != nil {
		Logger.Error(err.Error())
		return err
	}
	//设置文件发送的内容
	e.Text = []byte("【知识产权】您此次验证码为" + code + "，5分钟内有效，请您尽快验证！")
	//设置服务器相关的配置

	return e.Send(config.Url+":25", smtp.PlainAuth("", config.MyEmail, config.Password, config.Url))
}
