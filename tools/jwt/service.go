package jwt

import "time"

var _jwt *JWT

// Init 使用之前必须先调用Init初始化
func Init(secretKey string) {
	_jwt = NewJWT(secretKey)
}

func GenerateToken(claims CustomClaims, expirationTime time.Duration) (string, error) {
	return _jwt.GenerateToken(claims, expirationTime)
}

func GenToken(userId string, username string, expirationTime time.Duration) (string, error) {
	claims := CustomClaims{
		UserId:   userId,
		Username: username,
	}
	return _jwt.GenerateToken(claims, expirationTime)
}

func ParseToken(tokenString string) (*CustomClaims, error) {
	return _jwt.ParseToken(tokenString)
}
