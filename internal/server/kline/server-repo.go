package kline

import (
	"snake/internal/kline"
	"snake/internal/kline/repository"
)

type Repositories struct {
	repoKline kline.Repository
}

func (s *Server) initRepositories() {
	s.repos.repoKline = repository.New(s.clients.db)
}
