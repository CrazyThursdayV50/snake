package strategy

import (
	"snake/internal/strategy"

	gmap "github.com/CrazyThursdayV50/pkgo/builtin/map"
	"github.com/gin-gonic/gin"
)

type ListStrategiesStrategy struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	Balance          string `json:"balance"`
	Position         string `json:"position"`
	Profit           string `json:"profit"`
	ProfitPercentage string `json:"profit_percentage"`
}

type ListStrategiesData struct {
	List []*ListStrategiesStrategy `json:"list"`
}

func (s *Service) ListStrategies(ctx *gin.Context) {
	s.strategyLock.RLock()
	defer s.strategyLock.RUnlock()

	var data ListStrategiesData
	gmap.From(s.strategies).Iter(func(k int64, v strategy.Strategy) (bool, error) {
		absolute, percentage := v.Profit()
		data.List = append(data.List, &ListStrategiesStrategy{
			ID:               k,
			Name:             v.Name(),
			Balance:          v.Balance().Amount.String(),
			Position:         v.Position().Amount.String(),
			Profit:           absolute.String(),
			ProfitPercentage: percentage.String(),
		})
		return true, nil
	})

	ctx.JSON(200, successResponse(&data))
}
