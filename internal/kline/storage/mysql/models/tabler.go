package models

import (
	"fmt"

	"gorm.io/gorm"
)

type s interface {
	DB() string
}

func KlineTable[S s](interval S) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		return tx.Table(fmt.Sprintf("kline_%s", interval.DB()))
	}
}

func (m *Kline) ResetDefault(openTs, closeTs uint64, open string) {
	m.Amount = "0"
	m.Average = "0"
	m.Open = open
	m.Close = open
	m.High = open
	m.Low = open
	m.Volume = "0"
	m.TakerBuyAmount = "0"
	m.TakerBuyVolume = "0"
	m.OpenTs = int64(openTs)
	m.CloseTs = int64(closeTs)
	m.TradeCount = 0
}
