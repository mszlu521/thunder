package server

import (
	"github.com/gin-gonic/gin"
	"thunder/midd"
)

type RegisterHandler func(*gin.Engine)

func RegisterRouter(s *Server) *gin.Engine {
	r := gin.Default()
	if s.Cros {
		r.Use(midd.Cors())
	}
	if s.Auth {
		r.Use(midd.Auth(s.Ignores, s.NeedLogins))
	}
	if len(s.NeedCache) > 0 {
		r.Use(midd.Cache(s.NeedCache))
	}
	if s.RouterRegister != nil {
		s.RouterRegister(r)
	}
	return r
}
