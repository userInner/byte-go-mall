package utils

import (
	"byte-go-mall/constant/config"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

// GenerateToken 生成JWT token
func GenerateToken(userID int64) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(24 * time.Hour) // token有效期24小时

	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			IssuedAt:  nowTime.Unix(),
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := tokenClaims.SignedString(config.AppConfig.App.Jwt)

	return token, err
}

// ParseToken 解析JWT token
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return config.AppConfig.App.Jwt, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
