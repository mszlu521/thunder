package midd

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mszlu521/thunder/config"
)

func Cors(conf *config.Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowedOrigin, ok := isOriginAllowed(origin, conf.Cros)

		if ok {
			c.Header("Access-Control-Allow-Origin", allowedOrigin)
			// 当允许的源不是 '*' 时，才设置 credentials
			if allowedOrigin != "*" {
				c.Header("Access-Control-Allow-Credentials", "true")
			}
		}
		method := c.Request.Method
		c.Header("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Content-SessionType, Token")
		c.Header("Access-Control-Expose-Headers", "Access-Control-Allow-Headers, Token")
		c.Header("Access-Control-Max-Age", "172800")
		c.Header("Access-Control-Allow-Credentials", "true")
		if method == "OPTIONS" {
			//c.JSON(200, Controller.R(200, nil, "Options Request"))
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func isOriginAllowed(origin string, allowedOrigins []string) (string, bool) {
	// 完全匹配 "*"
	for _, o := range allowedOrigins {
		if o == "*" {
			return "*", true
		}
	}
	// 完全匹配域名
	for _, o := range allowedOrigins {
		if o == origin {
			return origin, true
		}
	}
	// 匹配子域通配符, e.g., *.example.com
	for _, o := range allowedOrigins {
		if strings.HasPrefix(o, "*.") {
			domainSuffix := strings.TrimPrefix(o, "*.")
			if strings.HasSuffix(origin, "."+domainSuffix) {
				// 找到了匹配的子域，返回具体的 origin
				return origin, true
			}
		}
	}
	return "", false
}
