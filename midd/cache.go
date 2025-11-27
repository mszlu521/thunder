package midd

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mszlu521/thunder/cache"
	"github.com/mszlu521/thunder/config"
	"github.com/mszlu521/thunder/logs"
	"github.com/mszlu521/thunder/res"
	"github.com/mszlu521/thunder/tools/crypro"
	"github.com/mszlu521/thunder/tools/gptr"

	"io"
	"net/http"
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
func Cache(cacheConfig *config.Cache) gin.HandlerFunc {
	return func(c *gin.Context) {
		//打印超时时间
		start := time.Now()
		for _, pattern := range cacheConfig.GetNeedCache() {
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
					logs.Infof("cache -------start----- time: %s", time.Since(start))
					if redisCache.Exist(cacheKey) {
						logs.Infof("cache Exist time: %s", time.Since(start))
						cacheData, err := redisCache.Get(cacheKey)
						if err == nil {
							c.Data(http.StatusOK, "application/json", []byte(cacheData))
							c.Abort()
							logs.Infof("cache time: %s", time.Since(start))
							return
						}
						logs.Errorf("get cache err: %v", err)
						c.Next()
					} else {
						c.Next()
					}
					if c.Writer.Status() == 200 {
						responseBody := writer.body
						var result res.Result
						err := json.Unmarshal(responseBody.Bytes(), &result)
						if err != nil {
							logs.Errorf("cache json Unmarshal err: %v", err)
						} else {
							if result.Code == res.OK {
								if cacheConfig.Expire == nil {
									cacheConfig.Expire = gptr.Of(int64(5 * 60)) //默认5分钟
								}
								err := redisCache.Set(cacheKey, string(responseBody.Bytes()), cacheConfig.GetExpire())
								if err != nil {
									logs.Errorf("cache redisCache.Set err: %v", err)
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
