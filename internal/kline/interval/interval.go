package interval

import (
	"errors"
	"time"
)

type Interval string

const (
	Interval1m  Interval = "1m"
	Interval3m  Interval = "3m"
	Interval5m  Interval = "5m"
	Interval15m Interval = "15m"
	Interval30m Interval = "30m"
	Interval1h  Interval = "1h"
	Interval2h  Interval = "2h"
	Interval4h  Interval = "4h"
	Interval6h  Interval = "6h"
	Interval8h  Interval = "8h"
	Interval12h Interval = "12h"
	Interval1d  Interval = "1d"
	Interval3d  Interval = "3d"
	Interval1w  Interval = "1w"
	Interval1M  Interval = "1M"
)

var dbNameMap = map[Interval]string{
	Interval1m:  "min_1",
	Interval3m:  "min_3",
	Interval5m:  "min_5",
	Interval15m: "min_15",
	Interval30m: "min_30",
	Interval1h:  "hour_1",
	Interval2h:  "hour_2",
	Interval4h:  "hour_4",
	Interval6h:  "hour_6",
	Interval8h:  "hour_8",
	Interval12h: "hour_12",
	Interval1d:  "day_1",
	Interval3d:  "day_3",
	Interval1w:  "week_1",
	Interval1M:  "month_1",
}

func Min1() Interval   { return Interval1m }
func Min3() Interval   { return Interval3m }
func Min5() Interval   { return Interval5m }
func Min15() Interval  { return Interval15m }
func Min30() Interval  { return Interval30m }
func Hour1() Interval  { return Interval1h }
func Hour2() Interval  { return Interval2h }
func Hour4() Interval  { return Interval4h }
func Hour6() Interval  { return Interval6h }
func Hour8() Interval  { return Interval8h }
func Hour12() Interval { return Interval12h }
func Day1() Interval   { return Interval1d }
func Day3() Interval   { return Interval3d }
func Week1() Interval  { return Interval1w }
func Month1() Interval { return Interval1M }

var all = []Interval{
	Interval1m,
	Interval3m,
	Interval5m,
	Interval15m,
	Interval30m,
	Interval1h,
	Interval2h,
	Interval4h,
	Interval6h,
	Interval8h,
	Interval12h,
	Interval1d,
	Interval3d,
	Interval1w,
	Interval1M,
}

func All() []Interval { return all }

func (i Interval) Duration() time.Duration {
	switch i {
	case Interval1m:
		return time.Minute
	case Interval3m:
		return time.Minute * 3
	case Interval5m:
		return time.Minute * 5
	case Interval15m:
		return time.Minute * 15
	case Interval30m:
		return time.Minute * 30
	case Interval1h:
		return time.Hour
	case Interval2h:
		return time.Hour * 2
	case Interval4h:
		return time.Hour * 4
	case Interval6h:
		return time.Hour * 6
	case Interval8h:
		return time.Hour * 8
	case Interval12h:
		return time.Hour * 12
	case Interval1d:
		return time.Hour * 24
	case Interval3d:
		return time.Hour * 72
	case Interval1w:
		return time.Hour * 168
	case Interval1M:
		return time.Hour * 720

	default:
		return time.Minute
	}
}

func (i Interval) String() string { return string(i) }
func (i Interval) DB() string     { return dbNameMap[i] }

// ParseString 将字符串解析为 Interval 类型
func ParseString(s string) (Interval, error) {
	for _, i := range all {
		if string(i) == s {
			return i, nil
		}
	}
	return Interval1m, errors.New("invalid interval string")
}

func ParseDuration(duration time.Duration) (Interval, error) {
	for _, i := range all {
		if duration == i.Duration() {
			return i, nil
		}
	}
	return Interval1m, errors.New("invalid duration")

}
