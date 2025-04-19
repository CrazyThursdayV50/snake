package strategy

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"snake/internal/service"
	"snake/internal/service/strategy"
	"sync"
	"time"

	"github.com/CrazyThursdayV50/pkgo/goo"
	"github.com/gin-gonic/gin"
)

type Services struct {
	strategy *strategy.Service
}

func (s *Server) initServices(ctx context.Context) {
	s.services.strategy = strategy.NewService(ctx, s.logger, s.repos.klineRepo)
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
	root.GET("strategy/list", s.strategy.ListStrategies)
	root.POST("strategy/test", s.strategy.Test)
	root.DELETE("strategy", s.strategy.StopStrategy)

	srv := http.Server{Handler: handler}

	goo.Goo(func() {
		wg.Add(1)
		defer wg.Done()
		err := srv.Serve(l)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}, func(err error) {
		if err != nil {
			fmt.Println("serve failed: ", err)
		}
	})

	goo.Go(func() {
		wg.Add(1)
		defer wg.Done()
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.TODO(), time.Second*10)
		defer cancel()
		_ = srv.Shutdown(ctx)
	})

}
