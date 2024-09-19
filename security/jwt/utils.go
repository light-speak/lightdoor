package jwt

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/light-speak/lightdoor/security/kitex_gen/token"
	"time"
)

func GetToken(user *token.UserIdRequest) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, buildClaim(user.UserId)).SignedString(key)
}

func buildClaim(userId int64) *claim {
	now := time.Now()
	return &claim{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * time.Duration(24) * time.Duration(30))), // 30天过期
			Issuer:    "com.staticoft.cos",
		},
	}
}

func GetUserId(token string) (*int64, error) {
	// 解析 token
	tokenType, err := jwt.ParseWithClaims(token, &claim{},
		func(token *jwt.Token) (interface{}, error) {
			return key, nil
		})
	if err != nil {
		return nil, err
	}
	// 验证 token 的有效性
	if claims, ok := tokenType.Claims.(*claim); ok && tokenType.Valid {
		return &claims.UserId, nil
	}
	return nil, errors.New("token解析失败或无效")
}
