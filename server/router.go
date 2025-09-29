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

func UseCustomMidd(conf *config.Config, engin *gin.Engine) {
	if len(conf.Server.Cros) > 0 {
		engin.Use(midd.Cors(conf.Server))
	}
	if conf.Auth.IsAuth {
		engin.Use(midd.Auth(conf.Auth))
	}
	if len(conf.Cache.NeedCache) > 0 {
		engin.Use(midd.Cache(conf.Cache))
	}
}
