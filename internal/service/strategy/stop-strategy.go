package strategy

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type StopStrategyParams struct {
	ID int64 `json:"id"`
}

type StopStrategyData struct{}

func (s *Service) StopStrategy(ctx *gin.Context) {
	var params StopStrategyParams
	err := ctx.ShouldBind(&params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, failResponse[TestData](err.Error(), "invalid params"))
		return
	}

	s.strategyLock.Lock()
	defer s.strategyLock.Unlock()
	strategy := s.strategies[params.ID]

	if strategy != nil {
		strategy.Stop()
		delete(s.strategies, params.ID)
	}

	ctx.JSON(http.StatusOK, successResponse(new(StopStrategyData)))
}
