package strategy

import (
	"context"
	"snake/internal/kline"
	"snake/internal/kline/interval"
	"time"

	gchan "github.com/CrazyThursdayV50/pkgo/builtin/chan"
	"github.com/CrazyThursdayV50/pkgo/log"
	defaultlogger "github.com/CrazyThursdayV50/pkgo/log/default"
)

type Server struct {
	cfg     *Config
	logger  log.Logger
	clients *Clients
	repos   *Repositories
}

func (s *Server) initLogger() {
	logger := defaultlogger.New(s.cfg.Log)
	logger.Init()
	s.logger = logger
}

func (s *Server) init() {
	s.initLogger()
	s.initClients()
	s.initRepositories()
}

func New(cfg *Config) *Server {
	return &Server{
		cfg:     cfg,
		clients: &Clients{},
		repos:   &Repositories{},
	}
}

func (s *Server) Run() {
	s.init()

	var ctx = context.Background()
	var now = time.Now()
	from := now.Add(-time.Hour)
	ch := s.repos.klineRepo.GetKlines(ctx, interval.Min1(), from.Unix()*1000)

	_, err := gchan.FromRead(ch).Iter(func(k int, v *kline.Kline) (bool, error) {
		s.logger.Infof("kline: %v", v)
		return true, nil
	})

	if err != nil {
		s.logger.Errorf("err: %v", err)
	}
}
