package req

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"thunder/errs"
)

func JsonParam(ctx *gin.Context, obj any) error {
	err := ctx.ShouldBindJSON(obj)
	if err != nil {
		log.Println(err)
		return errs.ErrParam
	}
	return nil
}

func QueryParam(ctx *gin.Context, obj any) error {
	err := ctx.ShouldBindQuery(obj)
	if err != nil {
		log.Println(err)
		return errs.ErrParam
	}
	return nil
}
func XMLParam(ctx *gin.Context, obj any) error {
	err := ctx.ShouldBindBodyWithXML(obj)
	if err != nil {
		log.Println(err)
		return errs.ErrParam
	}
	return nil
}

func PathParam(ctx *gin.Context, paramKey string) string {
	param := ctx.Param(paramKey)
	return param
}

func PathInt(ctx *gin.Context, paramKey string) (int64, error) {
	param := PathParam(ctx, paramKey)
	if param == "" {
		return 0, errs.ErrParam
	}
	return StringToInt64(param)
}

func StringToInt64(param string) (int64, error) {
	i, err := strconv.ParseInt(param, 10, 64)
	if err != nil {
		return 0, errs.ErrParam
	}
	return i, nil
}

func PathInArray(ctx *gin.Context, method string, urls []string) bool {
	path := ctx.Request.URL.Path
	for _, url := range urls {
		if path == url && method == ctx.Request.Method {
			return true
		}
	}
	return false
}

func GetUserId(ctx *gin.Context) int64 {
	value, exists := ctx.Get("userId")
	if !exists {
		return 0
	}
	return int64(value.(float64))
}
