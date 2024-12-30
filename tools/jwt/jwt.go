package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"thunder/tools/convert"
	"time"
)

// GenerateJWT 生成 JWT Token
func GenerateJWT(jwtSecret []byte, userId int64, exp int64) (string, error) {
	// 定义 token 的有效载荷
	claims := jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Second * time.Duration(exp)).Unix(),
	}
	// 创建 token 对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用密钥签名 token
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseJWT 解析并验证 JWT Token
func ParseJWT(tokenString string, jwtSecret []byte) (int64, error) {
	// 解析 token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 确保签名方法正确
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		return 0, err
	}
	// 验证 token 是否有效并提取声明信息
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId := convert.GetInterfaceToInt(claims["userId"])
		return userId, nil
	} else {
		return 0, fmt.Errorf("invalid token")
	}
}
func GenerateJWTCustom(jwtSecret []byte, userId any, exp int64) (string, error) {
	// 定义 token 的有效载荷
	claims := jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(time.Second * time.Duration(exp)).Unix(),
	}
	// 创建 token 对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用密钥签名 token
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseJWT 解析并验证 JWT Token
func ParseJWTCustom(tokenString string, jwtSecret []byte) (any, error) {
	// 解析 token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 确保签名方法正确
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})
	if err != nil {
		return 0, err
	}
	// 验证 token 是否有效并提取声明信息
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId := claims["userId"]
		return userId, nil
	} else {
		return 0, fmt.Errorf("invalid token")
	}
}
