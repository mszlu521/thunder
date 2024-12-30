package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	Name           string
	Port           int
	Host           string
	RouterRegister RegisterHandler
	OtherHandle    func(*gin.Engine)
	StopHandle     func(*gin.Engine)
	cancel         context.CancelFunc
	g              *gin.Engine
	Cros           bool
	Auth           bool
	Ignores        []string
	NeedLogins     []string
	NeedCache      []string
}

func NewServer(name string, port int, host string) *Server {
	return &Server{
		Name: name,
		Port: port,
		Host: host,
	}
}

func Default() *Server {
	return NewServer("default", 8080, "localhost")
}

func (s *Server) Run(ctx context.Context) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		s.cancel = cancel
		defer cancel()
	}
	go func() {
		//注册路由
		r := RegisterRouter(s)
		s.g = r
		//启动之前 还有别的事情去处理 比如启动监控等
		if s.OtherHandle != nil {
			s.OtherHandle(r)
		}
		//完成前置操作 取消超时context
		if s.cancel != nil {
			s.cancel()
		}
		//http接口
		if err := r.Run(fmt.Sprintf("%s:%d", s.Host, s.Port)); err != nil {
			log.Fatalf("gate gin run err:%v", err)
		}
	}()
	//期望有一个优雅启停 遇到中断 退出 终止 挂断
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGHUP)
	for {
		select {
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				s.stop()
				//time out
				return errors.New("server run time out")
			}
		case signal := <-c:
			switch signal {
			case syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT:
				s.stop()
				return errors.New("termination")
			case syscall.SIGHUP:
				s.stop()
				return errors.New("hang up")
			default:
				return nil
			}
		}
	}
}

func (s *Server) stop() {
	if s.StopHandle != nil {
		s.StopHandle(s.g)
	}
}
