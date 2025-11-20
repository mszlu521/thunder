package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims 自定义声明类型，内嵌 jwt.RegisteredClaims
// RegisteredClaims 包含官方定义的标准字段：(iss, sub, aud, exp, nbf, iat, jti)
// 我们这里增加了自定义字段 UserID 和 Username
type CustomClaims struct {
	UserId string `json:"userId"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// JWT 是我们的 JWT 工具结构体
type JWT struct {
	// SecretKey 用于签名的密钥，应该是保密的
	SecretKey []byte
}

// NewJWT 创建一个 JWT 工具实例
func NewJWT(secretKey string) *JWT {
	return &JWT{
		SecretKey: []byte(secretKey),
	}
}

// GenerateToken 生成一个 token
// 它接收自定义的声明和过期时间
func (j *JWT) GenerateToken(claims CustomClaims, expirationTime time.Duration) (string, error) {
	// 设置标准声明的过期时间
	// jwt.NewNumericDate 是一个辅助函数，用于将 time.Time 转换为 JWT 使用的 Unix 时间戳
	claims.RegisteredClaims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(expirationTime))
	// 设置签发时间
	claims.RegisteredClaims.IssuedAt = jwt.NewNumericDate(time.Now())
	// 设置签发人 (可选)
	// claims.RegisteredClaims.Issuer = "my-project-name"

	// 使用指定的签名方法和声明创建一个新的 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用秘钥进行签名，并获取完整的编码后的字符串 token
	return token.SignedString(j.SecretKey)
}

// ParseToken 解析和验证一个 token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	// jwt.ParseWithClaims 是核心解析函数
	// 1. tokenString: 要解析的 token 字符串
	// 2. &CustomClaims{}: 一个空的声明对象指针，用于告诉库如何解码 payload
	// 3. Keyfunc: 一个回调函数，用于提供验证签名所需的密钥
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 重要的安全校验：检查 token 使用的签名算法是否是我们期望的 HMAC 算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return j.SecretKey, nil
	})

	if err != nil {
		// 这里处理各种解析错误
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token has expired")
		}
		// 其他错误，例如签名无效、格式错误等
		return nil, errors.New("invalid token")
	}

	// 检查 token 是否有效，并且声明是否是我们定义的类型
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
