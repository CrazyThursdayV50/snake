package acl

import (
	"snake/internal/kline"
	"snake/internal/kline/storage/mysql/models"

	binance_connector "github.com/binance/binance-connector-go"
	"github.com/shopspring/decimal"
)

func DB2Service(_ int, src *models.Kline) (bool, *kline.Kline) {
	if src == nil {
		return false, nil
	}

	var dst kline.Kline
	amount, err := decimal.NewFromString(src.Amount)
	if err != nil {
		return false, nil
	}
	dst.A = amount

	close, err := decimal.NewFromString(src.Close)
	if err != nil {
		return false, nil
	}
	dst.C = close

	high, err := decimal.NewFromString(src.High)
	if err != nil {
		return false, nil
	}
	dst.H = high

	low, err := decimal.NewFromString(src.Low)
	if err != nil {
		return false, nil
	}
	dst.L = low

	open, err := decimal.NewFromString(src.Open)
	if err != nil {
		return false, nil
	}
	dst.O = open

	volume, err := decimal.NewFromString(src.Volume)
	if err != nil {
		return false, nil
	}
	dst.V = volume

	dst.S = src.OpenTs
	dst.E = src.CloseTs
	return true, &dst
}

func Ws2Service(_ int, src *binance_connector.WsKlineEvent) (bool, *kline.Kline) {
	volume, _ := decimal.NewFromString(src.Kline.Volume)
	amount, _ := decimal.NewFromString(src.Kline.QuoteVolume)
	close, _ := decimal.NewFromString(src.Kline.Close)
	high, _ := decimal.NewFromString(src.Kline.High)
	low, _ := decimal.NewFromString(src.Kline.Low)
	open, _ := decimal.NewFromString(src.Kline.Open)

	var m kline.Kline
	m.V = volume
	m.A = amount
	m.C = close
	m.E = src.Kline.EndTime
	m.H = high
	m.L = low
	m.O = open
	m.S = src.Kline.StartTime
	return true, &m
}
