package acl

import (
	"snake/internal/kline/storage/mysql/models"

	"github.com/CrazyThursdayV50/goex/binance/websocket-api/models/klines"
	"github.com/shopspring/decimal"
)

func ApiToDB(src klines.Kline) *models.Kline {
	volume, _ := decimal.NewFromString(src.Volume)
	amount, _ := decimal.NewFromString(src.Amount)

	var m models.Kline
	m.Volume = src.Volume
	m.Amount = src.Amount
	m.Close = src.Close
	m.CloseTs = src.CloseTs
	m.High = src.High
	m.Low = src.Low
	m.Open = src.Open
	m.OpenTs = src.OpenTs
	m.TakerBuyAmount = src.AmountBuy
	m.TakerBuyVolume = src.VolumeBuy
	m.TradeCount = src.TradeCount
	if !volume.IsZero() {
		m.Average = amount.Div(volume).String()
	}
	return &m
}

// func ApiToDB(src *binance_connector.KlinesResponse) *models.Kline {
// 	volume, _ := decimal.NewFromString(src.Volume)
// 	amount, _ := decimal.NewFromString(src.QuoteAssetVolume)

// 	var m models.Kline
// 	m.Volume = src.Volume
// 	m.Amount = src.QuoteAssetVolume
// 	m.Close = src.Close
// 	m.CloseTs = int64(src.CloseTime)
// 	m.High = src.High
// 	m.Low = src.Low
// 	m.Open = src.Open
// 	m.OpenTs = int64(src.OpenTime)
// 	m.TakerBuyAmount = src.TakerBuyQuoteAssetVolume
// 	m.TakerBuyVolume = src.TakerBuyBaseAssetVolume
// 	m.TradeCount = int64(src.NumberOfTrades)
// 	if !volume.IsZero() {
// 		m.Average = amount.Div(volume).String()
// 	}
// 	return &m
// }
