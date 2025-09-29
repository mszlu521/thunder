package wx

import (
	"github.com/gin-gonic/gin"
	"github.com/mszlu521/thunder/req"
	"github.com/mszlu521/thunder/res"
	"github.com/mszlu521/thunder/tools/crypro"
	"sort"
	"strings"
)

type TokenVerifyReq struct {
	Signature string `json:"signature" form:"signature"`
	Timestamp string `json:"timestamp" form:"timestamp"`
	EchoStr   string `json:"echostr" form:"echostr"`
	Nonce     string `json:"nonce" form:"nonce"`
}

func (r *TokenVerifyReq) verify(token string) bool {
	arr := []string{token, r.Timestamp, r.Nonce}
	sort.Slice(arr, func(i, j int) bool {
		return strings.ToLower(arr[i]) < strings.ToLower(arr[j])
	})
	sign := strings.Join(arr, "")
	sign = crypro.Sha1(sign)
	if sign == r.Signature {
		return true
	}
	return false
}

func Verify(token string, method string, urls []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !req.PathInArray(ctx, method, urls) {
			ctx.Abort()
			return
		}
		var tokenVerify TokenVerifyReq
		if err := req.QueryParam(ctx, &tokenVerify); err != nil {
			res.Error(ctx, err)
			return
		}
		if !tokenVerify.verify(token) {
			ctx.Abort()
		}
		ctx.Next()
	}
}
