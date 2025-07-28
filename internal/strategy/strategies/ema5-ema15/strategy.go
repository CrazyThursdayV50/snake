package strategy

import (
	"snake/internal/indicates/ema"
	"snake/internal/kline"
	"snake/internal/strategy"
)

type Stategy struct {
	*strategy.BaseStrategy
	ema5  *ema.EMA
	ema15 *ema.EMA
}

func (s *Stategy) Update(kline *kline.Kline) (*strategy.Signal, error) {
	if s.ema5.Timestamp == kline.S {
		if s.ema5.Update(kline) && s.ema15.Update(kline) {
			return s.Hold(), nil
		}
	}

}
