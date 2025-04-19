package strategy

import (
	"github.com/CrazyThursdayV50/pkgo/request/resty"
)

type Clients struct {
	resty *resty.Client
}

func (s *Server) initClients() {
	s.clients.resty = resty.New(
		resty.WithConfig(s.cfg.Clients.Resty),
		resty.WithLogger(s.logger),
	)
}
