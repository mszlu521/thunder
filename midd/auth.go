package midd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
	"strings"
	"thunder/config"
	"thunder/tools/jwt"
)

func Auth(ignores []string, needLogins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, pattern := range ignores {
			if isMatch(c.Request.URL.Path, pattern) {
				c.Next()
				return
			}
		}
		// 从请求头中获取 token
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			reject(c, "Authorization header is missing", needLogins)
			return
		}
		// 删除 "Bearer " 前缀，只保留 token 部分
		if len(tokenString) > 7 && strings.ToLower(tokenString[:7]) == "bearer " {
			tokenString = tokenString[7:]
		}
		userId, err := jwt.ParseJWTCustom(tokenString, []byte(config.Conf.Jwt.Secret))
		if err != nil {
			reject(c, "Invalid token", needLogins)
			return
		}
		c.Set("userId", userId)
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
