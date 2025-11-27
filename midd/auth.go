package midd

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mszlu521/thunder/config"
	"github.com/mszlu521/thunder/tools/jwt"
)

func Auth(authConf *config.Auth) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, pattern := range authConf.GetIgnores() {
			if isMatch(c.Request.URL.Path, pattern) {
				c.Next()
				return
			}
		}
		// 从请求头中获取 token
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			reject(c, "Authorization header is missing", authConf.NeedLogins)
			return
		}
		// 删除 "Bearer " 前缀，只保留 token 部分
		if len(tokenString) > 7 && strings.ToLower(tokenString[:7]) == "bearer " {
			tokenString = tokenString[7:]
		}
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			reject(c, "Invalid token", authConf.NeedLogins)
			return
		}
		c.Set("userId", claims.UserId)
		c.Next()
	}
}
func reject(ctx *gin.Context, errMsg string, needLoginUrls []string) {
	for _, v := range needLoginUrls {
		if isMatch(ctx.Request.URL.Path, v) {
			ctx.Next()
			return
		}
	}
	ctx.JSON(http.StatusUnauthorized, gin.H{"error": errMsg})
	ctx.Abort()
}
func isMatch(path string, pattern string) bool {
	// 将 `**` 替换为 `.*`，以适应正则表达式中的任意字符匹配
	pattern = strings.ReplaceAll(pattern, "**", ".*")
	// 将 `pattern` 包装成正则表达式
	regexPattern := fmt.Sprintf("^%s$", pattern)
	fmt.Println(regexPattern)
	matched, _ := regexp.MatchString(regexPattern, path)
	return matched
}
