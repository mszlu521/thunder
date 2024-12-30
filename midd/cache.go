package midd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"thunder/cache"
	"thunder/res"
	"thunder/tools/crypro"
	"time"
)

// CustomResponseWriter 自定义 ResponseWriter 来捕获响应内容
type CustomResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *CustomResponseWriter) Write(p []byte) (n int, err error) {
	// 将数据写入 body
	w.body.Write(p)
	// 使用 gin 原有的 ResponseWriter 将数据写回客户端
	return w.ResponseWriter.Write(p)
}
func Cache(needCache []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		//打印超时时间
		start := time.Now()
		for _, pattern := range needCache {
			if isMatch(c.Request.URL.Path, pattern) {
				//对数据进行缓存
				if c.Request.Method == http.MethodPost {
					body, err := io.ReadAll(c.Request.Body)
					if err != nil {
						c.AbortWithStatus(http.StatusInternalServerError)
						return
					}
					c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
					cacheKey := fmt.Sprintf("CACHE:%s:%s", c.Request.RequestURI, crypro.Md5(body))
					redisCache := cache.NewRedisCache()
					writer := &CustomResponseWriter{body: bytes.NewBuffer([]byte{}), ResponseWriter: c.Writer}
					c.Writer = writer
					log.Println("cache -------start----- time:", time.Since(start))
					if redisCache.Exist(cacheKey) {
						log.Println("cache Exist time:", time.Since(start))
						cacheData, err := redisCache.Get(cacheKey)
						if err == nil {
							c.Data(http.StatusOK, "application/json", []byte(cacheData))
							c.Abort()
							log.Println("cache time:", time.Since(start))
							return
						}
						log.Println("get cache err:", err)
						c.Next()
					} else {
						c.Next()
					}
					if c.Writer.Status() == 200 {
						responseBody := writer.body
						var result res.Result
						err := json.Unmarshal(responseBody.Bytes(), &result)
						if err != nil {
							log.Println("cache json Unmarshal err:", err)
						} else {
							if result.Code == res.OK {
								err := redisCache.Set(cacheKey, string(responseBody.Bytes()), 60*5)
								if err != nil {
									log.Println("cache redisCache.Set err:", err)
								}
							}
						}
					}
				}
			}
		}
		c.Next()
	}
}
