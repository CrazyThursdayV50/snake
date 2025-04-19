package utils

import (
	"snake/internal/kline/interval"
	"snake/internal/kline/storage/mysql/models"
	"time"

	"github.com/CrazyThursdayV50/pkgo/builtin/collector"
	"github.com/CrazyThursdayV50/pkgo/builtin/slice"
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

func FillKlines(klines []*models.Kline, interval interval.Interval, to int64) []*models.Kline {
	first := klines[0].OpenTs
	n := (to-first)/interval.Duration().Milliseconds() + 1
	tsSli := GenNextTimeToN(uint64(first), uint64(to), interval, int(n))
	klineGroup := collector.Map(klines, func(_ int, v *models.Kline) (bool, int64, *models.Kline) {
		return true, v.OpenTs, v
	})

	slice.From(tsSli...).Iter(func(k int, v int64) (bool, error) {

	})
}
