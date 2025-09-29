package req

import (
	"github.com/gin-gonic/gin"
	"github.com/mszlu521/thunder/errs"
	"github.com/mszlu521/thunder/logs"
	"github.com/mszlu521/thunder/res"
	"log"
	"strconv"
)

func JsonParam(c *gin.Context, obj any) error {
	err := c.ShouldBindJSON(obj)
	if err != nil {
		logs.Errorf("req parse param err : %v", err)
		res.Error(c, errs.ErrParam)
		return errs.ErrParam
	}
	return nil
}

func QueryParam(c *gin.Context, obj any) error {
	err := c.ShouldBindQuery(obj)
	if err != nil {
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
		res.Error(ctx, errs.ErrParam)
		return 0, errs.ErrParam
	}
	return StringToInt64(param)
}
func Path(ctx *gin.Context, paramKey string, value any) error {
	param := PathParam(ctx, paramKey)
	if param == "" {
		res.Error(ctx, errs.ErrParam)
		return errs.ErrParam
	}

	switch v := value.(type) {
	case *string:
		*v = param
	case *int:
		i, err := strconv.Atoi(param)
		if err != nil {
			res.Error(ctx, errs.ErrParam)
			return errs.ErrParam
		}
		*v = i
	case *int64:
		i, err := strconv.ParseInt(param, 10, 64)
		if err != nil {
			res.Error(ctx, errs.ErrParam)
			return errs.ErrParam
		}
		*v = i
	default:
		res.Error(ctx, errs.ErrParam)
		return errs.ErrParam
	}
	
	return nil
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

func GetUserId(ctx *gin.Context) (int64, bool) {
	value, exists := ctx.Get("userId")
	if !exists {
		res.Error(ctx, errs.ErrUnauthorized)
		return 0, false
	}
	return value.(int64), true
}
