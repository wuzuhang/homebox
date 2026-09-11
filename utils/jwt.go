package utils

import (
	"errors"
	"time"

	"home-box/config"

	"github.com/golang-jwt/jwt/v5"
)

// MyCustomClaims 自定义声明结构体，附加我们需要存储的信息
type MyCustomClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT Token
func GenerateToken(userID uint, username string) (string, error) {
	secret := []byte(config.AppConfig.JWT.Secret)
	expireHours := config.AppConfig.JWT.ExpireHours

	claims := MyCustomClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireHours) * time.Hour)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                                             // 签发时间
			Issuer:    "homebox",                                                                  // 签发人
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// ParseToken 解析并验证 JWT Token
func ParseToken(tokenString string) (*MyCustomClaims, error) {
	secret := []byte(config.AppConfig.JWT.Secret)

	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
