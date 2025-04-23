package backtest

import (
	"snake/internal/kline/storage/mysql/models"
	internalModels "snake/internal/models"

	"github.com/shopspring/decimal"
)

// convertKline 将 MySQL 的 Kline 模型转换为内部 Kline 模型
func convertKline(k *models.Kline) *internalModels.Kline {
	o, _ := decimal.NewFromString(k.Open)
	c, _ := decimal.NewFromString(k.Close)
	h, _ := decimal.NewFromString(k.High)
	l, _ := decimal.NewFromString(k.Low)
	v, _ := decimal.NewFromString(k.Volume)
	a, _ := decimal.NewFromString(k.Amount)

	return &internalModels.Kline{
		O: o,
		C: c,
		H: h,
		L: l,
		V: v,
		A: a,
		S: k.OpenTs,
		E: k.CloseTs,
	}
}

// convertKlines 批量转换 Kline 模型
func convertKlines(klines []*models.Kline) []*internalModels.Kline {
	result := make([]*internalModels.Kline, len(klines))
	for i, k := range klines {
		result[i] = convertKline(k)
	}
	return result
} 