package utils

import (
	"snake/internal/kline/interval"
	"snake/internal/kline/storage/mysql/models"
	m "snake/internal/models"
	"time"

	"github.com/shopspring/decimal"
)

func GetNextTime(currentTime uint64, interval interval.Interval) uint64 {
	return currentTime + uint64((interval.Duration()).Milliseconds())
}

func GetEndTimeByStartTime(startTime uint64, interval interval.Interval, count int64) uint64 {
	return startTime + uint64(interval.Duration().Milliseconds()*count) - uint64(time.Millisecond.Milliseconds())
}

func GetLastTime(currentTime uint64, interval interval.Interval) uint64 {
	return currentTime - uint64((interval.Duration()).Milliseconds())
}

func GetStartTimeByEndTime(endTime uint64, interval interval.Interval, count int64) uint64 {
	return endTime - uint64(interval.Duration().Milliseconds()*count) + uint64(time.Millisecond.Milliseconds())
}

func GenNextTimeToN(currentTime uint64, to uint64, interval interval.Interval, n int) []int64 {
	var sli = make([]int64, 0, n)
	sli = append(sli, int64(currentTime))
	for range n - 1 {
		next := GetNextTime(currentTime, interval)
		if next > to {
			return sli
		}

		sli = append(sli, int64(next))
		currentTime = next
	}

	return sli
}

func FillKlines(klines []*m.Kline, interval interval.Interval, to int64) []*m.Kline {
	if len(klines) == 0 {
		return klines
	}

	first := klines[0].S
	n := (to-first)/interval.Duration().Milliseconds() + 1
	tsSli := GenNextTimeToN(uint64(first), uint64(to), interval, int(n))
	
	// 创建时间到 Kline 的映射
	klineMap := make(map[int64]*m.Kline)
	for _, kline := range klines {
		klineMap[kline.S] = kline
	}

	// 填充缺失的 Kline
	var result []*m.Kline
	var lastKline *m.Kline

	for _, ts := range tsSli {
		if kline, exists := klineMap[ts]; exists {
			result = append(result, kline)
			lastKline = kline
		} else if lastKline != nil {
			// 使用上一个 Kline 的收盘价创建新的 Kline
			newKline := &m.Kline{
				O: lastKline.C, // 开盘价 = 上一个 Kline 的收盘价
				C: lastKline.C, // 收盘价 = 上一个 Kline 的收盘价
				H: lastKline.C, // 最高价 = 上一个 Kline 的收盘价
				L: lastKline.C, // 最低价 = 上一个 Kline 的收盘价
				V: decimal.Zero, // 成交量 = 0
				A: decimal.Zero, // 成交额 = 0
				S: ts,          // 开盘时间
				E: ts + interval.Duration().Milliseconds() - 1, // 收盘时间
			}
			result = append(result, newKline)
		}
	}

	return result
}


func FillKlinesDB(klines []*models.Kline, interval interval.Interval, to int64) []*models.Kline {
	if len(klines) == 0 {
		return klines
	}

	first := klines[0].OpenTs
	n := (to-first)/interval.Duration().Milliseconds() + 1
	tsSli := GenNextTimeToN(uint64(first), uint64(to), interval, int(n))
	
	// 创建时间到 Kline 的映射
	klineMap := make(map[int64]*models.Kline)
	for _, kline := range klines {
		klineMap[kline.OpenTs] = kline
	}

	// 填充缺失的 Kline
	var result []*models.Kline
	var lastKline *models.Kline

	for _, ts := range tsSli {
		if kline, exists := klineMap[ts]; exists {
			result = append(result, kline)
			lastKline = kline
		} else if lastKline != nil {
			// 使用上一个 Kline 的收盘价创建新的 Kline
			newKline := &models.Kline{
				Open: lastKline.Close, // 开盘价 = 上一个 Kline 的收盘价
				Close: lastKline.Close, // 收盘价 = 上一个 Kline 的收盘价
				High: lastKline.Close, // 最高价 = 上一个 Kline 的收盘价
				Low: lastKline.Close, // 最低价 = 上一个 Kline 的收盘价
				Volume: "0", // 成交量 = 0
				Amount: "0", // 成交额 = 0
				OpenTs: ts,          // 开盘时间
				CloseTs: ts + interval.Duration().Milliseconds() - 1, // 收盘时间
			}
			result = append(result, newKline)
		}
	}

	return result
}
