package res

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/mszlu521/thunder/errs"
	"net/http"
)

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}
type Page struct {
	Total       int64 `json:"total"`
	List        any   `json:"list"`
	PageSize    int64 `json:"pageSize"`
	CurrentPage int64 `json:"currentPage"`
}

func Success(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, Result{
		Code: OK,
		Data: data,
	})
}
func Error400(ctx *gin.Context) {
	ctx.Writer.WriteHeader(http.StatusBadRequest)
}

func Error500(ctx *gin.Context) {
	ctx.Writer.WriteHeader(http.StatusInternalServerError)
}

func Error(ctx *gin.Context, err error) {
	var er *errs.Errors
	switch {
	case errors.As(err, &er):
		if errors.Is(er, errs.ErrParam) {
			ctx.Writer.WriteHeader(http.StatusBadRequest)
		} else if errors.Is(er, errs.ErrUnauthorized) {
			ctx.Writer.WriteHeader(http.StatusUnauthorized)
		} else {
			Fail(ctx, er)
		}
	default:
		Error500(ctx)
	}
}

func Fail(ctx *gin.Context, e *errs.Errors) {
	ctx.JSON(http.StatusOK, Result{
		Code: e.Code,
		Msg:  e.Msg,
	})
}

func SetCookie(ctx *gin.Context, key string, value string, expire int64) {
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:   key,
		Value:  value,
		MaxAge: int(expire),
		Path:   "/",
	})
}
