package server

import (
	"github.com/gin-gonic/gin"
	"github.com/mszlu521/thunder/config"
	"github.com/mszlu521/thunder/midd"
)

// IRouter 定义路由注册接口
type IRouter interface {
	Register(engine *gin.Engine)
}
type CloseIRouter interface {
	IRouter
	Close() error
}
func UseCustomMidd(conf *config.Config, engin *gin.Engine) {
	if conf.Server != nil {
		if len(conf.Server.GetCros()) > 0 {
			engin.Use(midd.Cors(conf.Server))
		}
	}
	if conf.Auth != nil {
		if conf.Auth.GetIsAuth() {
			engin.Use(midd.Auth(conf.Auth))
		}
	}
	if conf.Cache != nil {
		if len(conf.Cache.GetNeedCache()) > 0 {
			engin.Use(midd.Cache(conf.Cache))
		}
	}
}