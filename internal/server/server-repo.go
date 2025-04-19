package server

import (
	"snake/internal/kline"
	"snake/internal/repository"
)

type Repositories struct {
	repoKline repository.KlineRepository
}

func (s *Server) initRepositories() {
	s.repos.repoKline = kline.NewRepository(s.clients.db)
}
