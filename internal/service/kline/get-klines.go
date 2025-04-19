package kline

import (
	"net/http"
	"snake/internal/kline"
	"snake/internal/kline/acl"
	"snake/internal/kline/interval"

	"github.com/CrazyThursdayV50/pkgo/builtin/collector"
	"github.com/gin-gonic/gin"
)

type GetKlinesParams struct {
	Interval string `form:"interval"`
	From     int64  `form:"from"`
	To       int64  `form:"to"`
}

type GetKlineData struct {
	List []*kline.Kline
}

type GetKlineReponse = Response[GetKlineData]

func (s *Service) GetKlines(ctx *gin.Context) {
	var params GetKlinesParams
	err := ctx.ShouldBind(&params)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, failResponse[GetKlineData](err.Error(), "invalid params"))
		return
	}

	klines, err := s.repoKline.List(ctx, interval.Interval(params.Interval), params.From, params.To)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, failResponse[GetKlineData](err.Error(), "list kline failed"))
		return
	}

	klineList := collector.Slice(klines, acl.DB2Service)
	var data GetKlineData
	data.List = klineList

	ctx.JSON(http.StatusOK, successResponse[GetKlineData](&data))
}
