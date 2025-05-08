package strategy

import (
	"context"
	"os"
	"os/signal"
	"sync"

	"github.com/CrazyThursdayV50/pkgo/log"
	defaultlogger "github.com/CrazyThursdayV50/pkgo/log/default"
)

type Server struct {
	cfg      *Config
	logger   log.Logger
	clients  *Clients
	repos    *Repositories
	services *Services
}

func (s *Server) initLogger() {
	logger := defaultlogger.New(s.cfg.Log)
	logger.Init()
	s.logger = logger
}

func (s *Server) init(ctx context.Context) {
	s.initLogger()
	s.initClients()
	s.initRepositories()
	s.initServices(ctx)
}

func New(cfg *Config) *Server {
	return &Server{
		cfg:      cfg,
		clients:  &Clients{},
		repos:    &Repositories{},
		services: &Services{},
	}
}

func (s *Server) Run() {
	var ctx, cancel = context.WithCancel(context.Background())
	s.init(ctx)

	var wg sync.WaitGroup
	s.services.Run(ctx, s.cfg.Service, &wg)

	var signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	s.logger.Warn("SERVER EXIT ...")
	cancel()

	wg.Wait()
}
