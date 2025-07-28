package kline

import (
	"context"
	"fmt"
	k "snake/internal/kline"
	"snake/internal/kline/interval"
	"snake/internal/service/kline"
	"time"

	"github.com/CrazyThursdayV50/pkgo/builtin/slice"
	"github.com/CrazyThursdayV50/pkgo/goo"
)

func (r *Repository) ListKlines(ctx context.Context, interval interval.Interval, from int64) <-chan *Kline {
	var result kline.GetKlineReponse
	_, err := r.client.Request(ctx).SetQueryParams(map[string]string{
		"interval": interval.String(),
		"from":     fmt.Sprintf("%d", 0),
		"to":       fmt.Sprintf("%d", time.Now().Unix()*1000),
	}).
		SetResult(&result).
		SetError(&result).
		Get(r.endpoint)

	if err != nil {
		r.logger.Errorf("get klines failed: %v", err)
		return nil
	}

	if result.Error != "" {
		r.logger.Errorf("get klines failed: %v: %v", result.Message, result.Error)
		return nil
	}

	list := slice.From(result.Data.List...)
	var ch = make(chan *k.Kline, list.Len())
	goo.Go(func() {
		defer close(ch)
		list.Iter(func(k int, v *k.Kline) (bool, error) {
			ch <- v
			return true, nil
		})
	})

	return ch
}
