package utils

import (
	"fmt"
	"github.com/smartwalle/alipay/v3"
	"github.com/spf13/viper"
)

type AliPay struct {
	AppID       string `json:"app_id"`
	PrivateKey  string `json:"private_key"`
	PublicKey   string `json:"public_key"`
	NotifyUrl   string `json:"notify_url"`
	ReturnUrl   string `json:"return_url"`
	ProductCode string `json:"product_code"`
}

// getAlipayConfig 读取Email配置文件
func getAlipayConfig() AliPay {
	m := AliPay{}
	viper.SetConfigName("alipay")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	if err := viper.ReadInConfig(); err != nil {
		Logger.Error("读取配置错误")
	}
	m.AppID = viper.GetString("alipay.app_id")
	m.PrivateKey = viper.GetString("alipay.private_key")
	m.PublicKey = viper.GetString("alipay.public_key")
	m.NotifyUrl = viper.GetString("alipay.notify_url")
	m.ReturnUrl = viper.GetString("alipay.return_url")
	m.ProductCode = viper.GetString("alipay.product_code")
	return m
}

// PaymentOrderCreation 支付订单创建
func PaymentOrderCreation(Subject, OutTradeNo, TotalAmount string) string {
	alipayConfig := getAlipayConfig()
	fmt.Println("alipayConfig:", alipayConfig)
	client, err := alipay.New(alipayConfig.AppID, alipayConfig.PrivateKey, false)
	if err != nil {
		Logger.Error(err.Error())
		return ""
	}
	//加载阿里云公钥证书
	if err := client.LoadAliPayPublicKey(alipayConfig.PublicKey); err != nil {
		Logger.Error(err.Error())
		return ""
	}
	//创建交易支付订单
	var pay = alipay.TradeWapPay{}
	pay.NotifyURL = alipayConfig.NotifyUrl
	pay.ReturnURL = alipayConfig.ReturnUrl
	pay.Subject = Subject
	pay.OutTradeNo = OutTradeNo
	pay.TotalAmount = TotalAmount
	pay.ProductCode = alipayConfig.ProductCode

	url, err := client.TradeWapPay(pay)
	if err != nil {
		Logger.Error(err.Error())
		return ""
	}
	return url.String()
}

// PaymentOrderCreationPatent 专利支付订单创建
func PaymentOrderCreationPatent(ApplicationNumber, TotalAmount string) string {
	return PaymentOrderCreation("专利年费", ApplicationNumber, TotalAmount)
}

// PaymentOrderCreationArticle 著作支付订单创建
func PaymentOrderCreationArticle(ApplicationNumber, TotalAmount string) string {
	return PaymentOrderCreation("著作年费", ApplicationNumber, TotalAmount)
}

// PaymentOrderCreationTrademark 商标支付订单创建
func PaymentOrderCreationTrademark(ApplicationNumber, TotalAmount string) string {
	return PaymentOrderCreation("商标年费", ApplicationNumber, TotalAmount)
}
