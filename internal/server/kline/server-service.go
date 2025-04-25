package kline

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"snake/internal/service"
	"snake/internal/service/kline"
	"sync"
	"time"

	"github.com/CrazyThursdayV50/pkgo/goo"
	"github.com/CrazyThursdayV50/pkgo/websocket/server"
	"github.com/gin-gonic/gin"
)

type Services struct {
	kline *kline.Service
}

func (s *Server) initServices() {
	s.Wsserver = server.New(
		server.WithLogger(s.logger),
		server.WithTracer(s.tracer.NewTracer("websocket")),
		server.WithHandler(HandleKlineMessage),
	)
	s.Services.kline = kline.NewService(s.logger, s.Wsserver, s.repos.repoKline)
}

func (s *Services) Run(ctx context.Context, cfg *service.Config, wg *sync.WaitGroup) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	fmt.Printf("listen service at %s\n", l.Addr().String())

	handler := gin.Default()
	root := handler.Group("/")
	root.GET("ping", s.kline.Ping)
	root.GET("kline", s.kline.GetKlines)
	root.GET("ws", s.kline.Subscribe)

	srv := http.Server{Handler: handler}

	goo.Goo(func() {
		wg.Add(1)
		defer wg.Done()
		err := srv.Serve(l)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}, func(error) {})

	goo.Go(func() {
		wg.Add(1)
		defer wg.Done()
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second*10)
		defer cancel()
		_ = srv.Shutdown(ctx)
	})

}
