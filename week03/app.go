package week03

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
)

const (
	OPENED = iota
	CLOSED
)

type App struct {
	closed  int32          //0开启 1关闭
	Sign    chan os.Signal //信号量
	Servers []*http.Server //存储所有的http.Server
}

var mu sync.Mutex

func NewApp() *App {
	sign := make(chan os.Signal)
	signal.Notify(sign, syscall.SIGINT, syscall.SIGTERM)
	return &App{Sign: sign, Servers: make([]*http.Server, 0, 10)}
}

//Close 关闭
func (a *App) Close() {
	atomic.StoreInt32(&a.closed, CLOSED)
}

//isClosed 是否已关闭
func (a *App) isClosed() bool {
	return atomic.LoadInt32(&a.closed) == CLOSED
}

//addServers 往a.Servers中添加server
func (a *App) addServers(s *http.Server) {
	mu.Lock()
	defer mu.Unlock()
	a.Servers = append(a.Servers, s)
}

//Server 添加Server
func (a *App) Server(c context.Context, addr string, handler http.Handler) error {
	s := http.Server{
		Addr:    addr,
		Handler: handler,
	}

	a.addServers(&s)

	//ctx.Done时，发送信号至s.doneChan
	go func() {
		ctx, _ := context.WithCancel(c)
		<-ctx.Done()
		_ = s.Shutdown(ctx)
	}()

	//接收信号s.doneChan
	err := s.ListenAndServe()

	//如果当前server异常关闭则通知主协程关闭其他server
	if !a.isClosed() {
		log.Println(addr, "出现异常，通知其他server......")
		a.Sign <- syscall.SIGINT
	}

	log.Printf("%s 服务已退出 %v\n", addr, err)
	return err
}

//StopServer 关闭某个Server
func (a *App) StopServer(addr string) {
	mu.Lock()
	defer mu.Unlock()
	for _, server := range a.Servers {
		if server.Addr == addr {
			_ = server.Shutdown(context.Background())
		}
	}
	return
}
