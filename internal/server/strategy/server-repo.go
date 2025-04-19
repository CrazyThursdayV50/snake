package strategy

import (
	"snake/internal/strategy"
	"snake/internal/strategy/repository/kline"
)

type Repositories struct {
	klineRepo strategy.KlineRepository
}

func (s *Server) initRepositories() {
	s.repos.klineRepo = kline.New(s.logger, s.cfg.Repository.Kline, s.clients.resty)
}
