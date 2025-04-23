package acl

import (
	"snake/internal/kline/storage/mysql/models"

	binance_connector "github.com/binance/binance-connector-go"
	"github.com/shopspring/decimal"
)

func ApiToDB(src *binance_connector.KlinesResponse) *models.Kline {
	volume, _ := decimal.NewFromString(src.Volume)
	amount, _ := decimal.NewFromString(src.QuoteAssetVolume)

	var m models.Kline
	m.Volume = src.Volume
	m.Amount = src.QuoteAssetVolume
	m.Close = src.Close
	m.CloseTs = int64(src.CloseTime)
	m.High = src.High
	m.Low = src.Low
	m.Open = src.Open
	m.OpenTs = int64(src.OpenTime)
	m.TakerBuyAmount = src.TakerBuyQuoteAssetVolume
	m.TakerBuyVolume = src.TakerBuyBaseAssetVolume
	m.TradeCount = int64(src.NumberOfTrades)
	if !volume.IsZero() {
		m.Average = amount.Div(volume).String()
	}
	return &m
}

func WsToDB(src *binance_connector.WsKlineEvent) *models.Kline {
	volume, _ := decimal.NewFromString(src.Kline.Volume)
	amount, _ := decimal.NewFromString(src.Kline.QuoteVolume)

	var m models.Kline
	m.Volume = src.Kline.Volume
	m.Amount = src.Kline.QuoteVolume
	m.Close = src.Kline.Close
	m.CloseTs = src.Kline.EndTime
	m.High = src.Kline.High
	m.Low = src.Kline.Low
	m.Open = src.Kline.Open
	m.OpenTs = src.Kline.StartTime
	m.TakerBuyAmount = src.Kline.ActiveBuyQuoteVolume
	m.TakerBuyVolume = src.Kline.ActiveBuyVolume
	m.TradeCount = src.Kline.TradeNum
	m.Average = amount.Div(volume).String()
	return &m
}