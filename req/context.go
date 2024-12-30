package req

import "github.com/gin-gonic/gin"

func GetInt64(c *gin.Context, key string) int64 {
	value, ok := c.Get(key)
	if ok {
		return int64(value.(float64))
	}
	return 0
}
