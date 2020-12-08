package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// 首先要搞一个http服务
type Server struct {
	srv *http.Server
}

type GetHandler struct {
}

func (h *GetHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	//模拟耗时请求
	fmt.Println("进入请求")
	time.Sleep(100 * time.Second) //若是在请求过程中被停止就会报错
	_, _ = w.Write([]byte("进入请求"))
}

// 创建http服务，并模拟一个耗时接口
func NewServer() *Server {
	mux := http.NewServeMux()
	mux.Handle("/get", &GetHandler{})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	return &Server{srv: srv}
}

func (s *Server) Start() error {
	fmt.Printf("[HTTP] 监听端口: %s\n", s.srv.Addr)
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func main() {
	stop := make(chan struct{})
	//使用一个携带上下文信息的errgroup
	g, ctx := errgroup.WithContext(context.Background())
	svr := NewServer()
	//启动服务，当任何errorgroup中的goroutine产生error时，关闭httpServer
	g.Go(func() error {
		fmt.Println("开启http服务")
		go func() {
			<-ctx.Done()
			fmt.Println("http 上下文取消咯")
			ctx2, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := svr.Shutdown(ctx2); err != nil {
				fmt.Println("哦吼，炸了！", err)
			}
			stop <- struct{}{}
			fmt.Println("Http服务退出了")
		}()
		return svr.Start()
	})

	//监听signal信号，当接收到退出相关信号退出
	g.Go(func() error {
		quit := make(chan os.Signal)
		//监听到指定信号就给quit
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)
		for {
			fmt.Println("等待退出信号")
			select {
			case <-ctx.Done():
				fmt.Println("上下文关闭")
				return ctx.Err()
			case <-quit:
				return errors.New("收到退出信号") //这个人造错误用来退出errgroup
			}
		}
	})

	//其他后台任务
	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("background Context 被取消啦")
				return ctx.Err()
			default:
				fmt.Println("搞点事情？")
				time.Sleep(1 * time.Second)
			}
		}
	})
	err := g.Wait()
	fmt.Println(err)
	<-stop
	fmt.Println("都停下来了")
}
