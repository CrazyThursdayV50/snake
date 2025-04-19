package kline

import (
	"github.com/gin-gonic/gin"
)

func (s *Service) Subscribe(ctx *gin.Context) {
	s.ws.Run(ctx, ctx.Writer, ctx.Request, nil)
}
