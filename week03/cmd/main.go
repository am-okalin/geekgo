package main

import (
	"context"
	"fmt"
	"geekgo/week03"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"
)

func handelApp() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintln(w, "app: 8080")
	})
	return mux
}

func main() {
	app := week03.NewApp()
	group := errgroup.Group{}
	ctx, cancel := context.WithCancel(context.Background())

	group.Go(func() error {
		return app.Server(ctx, ":8080", handelApp())
	})

	group.Go(func() error {
		return app.Server(ctx, ":8001", http.DefaultServeMux)
	})

	//测试:: 键入ctrl+c 是否能通知到所有server关闭
	go func() {
		<-app.Sign
		app.Close()
		cancel()
	}()

	//测试:: 10秒后某个server异常退出 是否能通知到其他server
	go func() {
		<-time.After(10 * time.Second)
		app.StopServer(":8080")
	}()

	//pending...
	err := group.Wait()
	log.Println("所有服务已退出", err)
}
