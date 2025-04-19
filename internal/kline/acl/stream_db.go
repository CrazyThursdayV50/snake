package acl

import (
	"snake/internal/kline/storage/mysql/models"

	"github.com/CrazyThursdayV50/goex/binance/websocket-streams/models/klines"
	"github.com/shopspring/decimal"
)

func WsToDB(src *klines.Data) *models.Kline {
	volume, _ := decimal.NewFromString(src.Volume)
	amount, _ := decimal.NewFromString(src.Amount)

	var m models.Kline
	m.Volume = src.Volume
	m.Amount = src.Amount
	m.Close = src.Close
	m.CloseTs = src.CloseTime
	m.High = src.High
	m.Low = src.Low
	m.Open = src.Open
	m.OpenTs = src.OpenTime
	m.TakerBuyAmount = src.TakerBuyAmount
	m.TakerBuyVolume = src.TakerBuyVolume
	m.TradeCount = src.TradeCount
	m.Average = amount.Div(volume).String()
	return &m
}
