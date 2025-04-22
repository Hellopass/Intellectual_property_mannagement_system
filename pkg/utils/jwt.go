package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

// 定义 JWT 密钥（生产环境应从安全位置获取）
var jwtKey = []byte(GinConfig.JwtKey)

// CustomClaims 自定义 Claims 结构体
type CustomClaims struct {
	UserID    int    `json:"user_id"`
	Authority string `json:"authority"`  //权限控制
	UserName  string `json:"user_name" ` //姓名
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT Token
func GenerateToken(userID int, authority string, userName string) (string, error) {

	claims := CustomClaims{
		UserID:    userID,
		Authority: authority,
		UserName:  userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour * 7)), // 过期时间-7天
			IssuedAt:  jwt.NewNumericDate(time.Now()),                         // 签发时间
			Issuer:    "Zhu",                                                  // 签发者
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// JWTMiddleware JWT 中间件
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 获取 token
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未提供认证令牌"})
			return
		}

		// 验证 token 格式
		if len(tokenString) < 7 || tokenString[:7] != "Bearer " {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "令牌格式错误"})
			return
		}
		tokenString = tokenString[7:]

		// 解析并验证 Token
		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("无效的签名方法: %v", token.Header["alg"])
			}
			return jwtKey, nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "令牌验证失败", "details": err.Error()})
			return
		}

		// 类型断言获取 Claims
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			c.Set("userID", claims.UserID)       // 将用户信息存入上下文 --设置需要的存入的信息
			c.Set("authority", claims.Authority) //用户权限
			c.Set("userName", claims.UserName)   //用户名字
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "无效令牌"})
		}
	}
}

// ParseToken 解析 JWT Token 并返回 Claims 信息
// 参数：tokenString - 去掉 "Bearer " 前缀的纯 Token 字符串
// 返回：解析后的 Claims 指针和错误信息
func ParseToken(tokenString string) (*CustomClaims, error) {
	// 解析 Token 并验证签名
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("非预期的签名方法: %v", token.Header["alg"])
		}
		return jwtKey, nil // 直接从配置获取最新密钥
	})

	if err != nil {
		return nil, err
	}

	// 类型断言获取 Claims
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("无效的 Token 声明")
}

//// 在需要解析的地方调用
//func SomeHandler(c *gin.Context) {
//	// 从请求头获取完整 Token
//	authHeader := c.GetHeader("Authorization")
//	if authHeader == "" {
//		// 处理无 Token 情况
//	}
//
//	// 提取纯 Token 部分（去掉 Bearer 前缀）
//	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
//
//	// 调用解析函数
//	claims, err := utils.ParseToken(tokenString)
//	if err != nil {
//		// 根据错误类型处理
//		switch {
//		case strings.Contains(err.Error(), "过期"):
//			c.JSON(http.StatusUnauthorized, gin.H{"error": "登录已过期"})
//		case strings.Contains(err.Error(), "签名验证失败"):
//			c.JSON(http.StatusUnauthorized, gin.H{"error": "非法 Token"})
//		default:
//			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证失败"})
//		}
//		return
//	}
//
//	// 使用解析后的字段
//	fmt.Printf("用户ID: %d\n权限: %s\n用户名: %s\n",
//		claims.UserID,
//		claims.Authority,
//		claims.UserName,
//	)
//
//	// 继续业务逻辑...
//}
