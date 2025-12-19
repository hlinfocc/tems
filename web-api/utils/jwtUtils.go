package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// 定义JWT Claims结构体
type CustomClaims struct {
	UserID   uint64 `json:"user_id"`
	ClassId  uint64 `json:"class_id"`
	RealName string `json:"realName"`
	Account  string `json:"account"`
	Status   int    `json:"status"`
	UserType int    `json:"userType"`
	Role     string `json:"role"`
	AppType  int    `json:"appType"` // 0 后台管理 1 小程序
	jwt.RegisteredClaims
}

// 密钥（实际应用中应从安全配置获取）
var jwtSecret = []byte("ylcxy.cn")

func GenJwtToken(id uint64, realName string, classId uint64, account string, status int, userType int, role string, appType int) (string, error) {
	// 生成JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		UserID:   id, // 用户ID
		ClassId:  classId,
		RealName: realName,
		Account:  account,
		Status:   status,
		UserType: userType,
		Role:     role,
		AppType:  appType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 过期时间
			Issuer:    "ylcxy.cn",                                         // 签发者
		},
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func JWTParse(tokenString string, appType int) (CustomClaims, error) {
	// 验证并解析token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return CustomClaims{}, fmt.Errorf("无效的令牌")
	}
	claims, ok := token.Claims.(*CustomClaims)
	if ok {
		if claims.AppType != appType {
			return CustomClaims{}, fmt.Errorf("无效的令牌,令牌不匹配")
		}
		return *claims, nil
	}
	return CustomClaims{}, fmt.Errorf("令牌解析失败")
}
