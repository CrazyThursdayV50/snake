package strategy

import (
	"context"
	"fmt"
	"net/http"
	"snake/internal/kline"
	"snake/internal/kline/interval"
	"snake/internal/strategy/strategies/ma_cross"
	"sync/atomic"
	"time"

	"github.com/CrazyThursdayV50/pkgo/worker"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type TestParams struct {
	// 余额
	Balance string `json:"balance"`
	// 仓位
	Position string `json:"position"`
	// 仓位总成本
	Cost string `json:"cost"`
}

type TestData struct {
	Strategy string `json:"strategy"`
	ID       int64  `json:"id"`
}

func (s *Service) Test(ctx *gin.Context) {
	var params TestParams
	err := ctx.ShouldBind(&params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, failResponse[TestData](err.Error(), "invalid params"))
		return
	}

	strategyCtx, cancel := context.WithCancel(s.ctx)
	strategy := ma_cross.New(strategyCtx, cancel)

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

	err = strategy.Init(position, balance)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, failResponse[TestData](err.Error(), "init strategy failed"))
		return
	}

	var interval = interval.Min1()
	from := time.Now().Add(-time.Hour)
	ch := s.klineRepo.GetKlines(strategyCtx, interval, from.Unix()*1000)

	id := atomic.AddInt64(&s.id, 1)
	s.strategyLock.Lock()
	s.strategies[id] = strategy
	s.strategyLock.Unlock()

	var data TestData
	data.ID = id

	worker, _ := worker.New(fmt.Sprintf("Turtle-%d", id), func(job *kline.Kline) {
		signal, err := strategy.Update(job)
		if err != nil {
			s.logger.Errorf("udpate strategy failed: %v", err)
			return
		}

		s.logger.Infof("signal: %#v", signal)
	})

	worker.WithLogger(s.logger)
	worker.WithContext(strategyCtx)
	worker.WithTrigger(ch)
	worker.Run()

	ctx.JSON(http.StatusOK, successResponse(&data))
}
