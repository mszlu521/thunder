package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mszlu521/thunder/config"
	"github.com/mszlu521/thunder/event"
	"github.com/mszlu521/thunder/pay/wxPay"
)

// Server 是我们应用的核心结构体
type Server struct {
	Engine     *gin.Engine
	httpServer *http.Server
	conf       *config.Config
}

// NewServer 创建一个新的 Server 实例
func NewServer(conf *config.Config) *Server {
	if conf.Server == nil {
		panic("Server config is required ")
	}
	// 根据配置设置 Gin 模式
	gin.SetMode(conf.Server.GetMode())

	// 初始化微信支付
	if conf.Pay != nil {
		wxPayConfig := conf.Pay.WxPay
		if wxPayConfig != nil {
			wxPay.Init(wxPayConfig)
		}
	}

	engine := gin.Default() // 使用默认的中间件 (Logger, Recovery)
	//自定义的一些中间件，可通过配置文件开启，减少代码重复书写
	UseCustomMidd(conf, engine)
	return &Server{
		Engine: engine,
		conf:   conf,
	}
}

// RegisterRouters 批量注册路由
// 参数是实现了 IRouter 接口的实例
func (s *Server) RegisterRouters(event event.IEvent, routers ...IRouter) {
	//注册 事件
	event.Register()
	for _, r := range routers {
		r.Register(s.Engine)
	}
	log.Println("Routers registered successfully.")
}

// Start 启动服务并实现优雅启停
func (s *Server) Start() {
	// 从配置中获取服务器地址和超时设置
	address := fmt.Sprintf("%s:%d", s.conf.Server.GetHost(), s.conf.Server.GetPort())
	readTimeout := s.conf.Server.GetReadTimeout()
	writeTimeout := s.conf.Server.GetWriteTimeout()

	s.httpServer = &http.Server{
		Addr:         address,
		Handler:      s.Engine,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	// 启动 http server (goroutine)
	go func() {
		log.Printf("Server starting on http://%s", address)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// ---- 优雅启停逻辑 ----
	// 创建一个 channel 用于接收系统信号
	quit := make(chan os.Signal, 1)
	// 监听 SIGINT (Ctrl+C) 和 SIGTERM 信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 阻塞在此，直到接收到信号
	<-quit
	log.Println("Shutting down server...")

	// 创建一个有5秒超时的 context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 调用 Shutdown() 优雅地关闭服务器
	// 这会等待正在处理的请求完成，但不会接受新请求
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server exited gracefully.")
}