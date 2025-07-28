package strategy

import (
	"context"
	"encoding/json"
	"net/http"
	"snake/internal/kline"
	"snake/internal/kline/interval"
	"snake/internal/strategy"
	"snake/internal/strategy/strategies/ma_cross"

	gchan "github.com/CrazyThursdayV50/pkgo/builtin/chan"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type BackTestParams struct {
	StrategyName   string          `json:"strategy_name"`
	StrategyParams json.RawMessage `json:"strategy_params"`
	// kline interval
	Interval string `json:"interval"`
	From     int64  `json:"from"`
	// 余额
	Balance string `json:"balance"`
	// 仓位
	Position string `json:"position"`
	// 仓位总成本
	Cost *string `json:"cost"`
}

type StrategyParams struct{}

type BackTestData struct {
	Name             string `json:"name"`
	Balance          string `json:"balance"`
	Position         string `json:"position"`
	Profit           string `json:"profit"`
	ProfitPercentage string `json:"profit_percentage"`
}

func (s *Service) BackTest(ctx *gin.Context) {
	var params BackTestParams
	err := ctx.ShouldBind(&params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, failResponse[TestData](err.Error(), "invalid params"))
		return
	}

	interval, err := interval.ParseString(params.Interval)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, failResponse[TestData](err.Error(), "invalid params"))
		return
	}

	position, err := decimal.NewFromString(params.Position)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, failResponse[TestData](err.Error(), "invalid position"))
		return
	}

	balance, err := decimal.NewFromString(params.Balance)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, failResponse[TestData](err.Error(), "invalid balance"))
		return
	}

	var costs []decimal.Decimal
	if params.Cost != nil {
		cost, err := decimal.NewFromString(*params.Cost)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, failResponse[TestData](err.Error(), "invalid balance"))
			return
		}
		costs = append(costs, cost)
	}

	strategyCtx, cancel := context.WithCancel(s.ctx)

	var st strategy.Strategy
	switch params.StrategyName {
	case strategy.STRATEGY_MA_CROSS:
		st = ma_cross.New(strategyCtx, cancel)

	default:
		st = ma_cross.New(strategyCtx, cancel)
	}

	err = st.Init(position, balance, costs...)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, failResponse[TestData](err.Error(), "init strategy failed"))
		return
	}

	defer st.Stop()

	ch := s.klineRepo.GetKlines(strategyCtx, interval, params.From)

	_, err = gchan.FromRead(ch).Iter(func(k int, v *kline.Kline) (bool, error) {
		signal, err := st.Update(v)
		if err != nil {
			s.logger.Errorf("udpate strategy failed: %v", err)
			return false, err
		}

		s.logger.Infof("SIGNAL: %+#v", signal)
		return true, nil
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, failResponse[*BackTestData](err.Error(), "backtest failed"))
		return
	}

	absolute, percentage := st.Profit()
	data := &BackTestData{
		Name:             st.Name(),
		Balance:          st.Balance().Amount.String(),
		Position:         st.Position().Amount.String(),
		Profit:           absolute.String(),
		ProfitPercentage: percentage.String(),
	}

	ctx.JSON(http.StatusOK, successResponse(data))
}
