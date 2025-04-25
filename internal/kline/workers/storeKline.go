package workers

import (
	"context"
	"fmt"
	"snake/internal/kline"
	"snake/internal/kline/interval"
	"snake/internal/kline/storage/mysql/models"
	"time"

	"github.com/CrazyThursdayV50/pkgo/goo"
	"github.com/CrazyThursdayV50/pkgo/log"
	"github.com/CrazyThursdayV50/pkgo/worker"
)

func tryFunc(f func() error, handler func(error), count int) {
	for range count {
		err := f()
		if err == nil {
			return
		}

		if handler != nil {
			handler(err)
		}
	}

	panic("max tries reached")
}

func StoreKline(ctx context.Context, logger log.Logger, interval interval.Interval, repoKline kline.Repository) func(*models.Kline) {
	var klinePipe = make(chan *models.Kline)
	goo.Go(func() {
		var klinesCache = make([]*models.Kline, 0, 1000)
		var ticker = time.NewTicker(time.Second)

		for {
			select {
			case <-ticker.C:
				if len(klinesCache) == 0 {
					break
				}

				err := repoKline.Insert(ctx, interval, klinesCache)
				if err == nil {
					klinesCache = make([]*models.Kline, 0, 1000)
				}

			case kline := <-klinePipe:
				klinesCache = append(klinesCache, kline)
				if len(klinesCache) < 1000 {
					break
				}

				err := repoKline.Insert(ctx, interval, klinesCache)
				if err == nil {
					klinesCache = make([]*models.Kline, 0, 1000)
				}
			}
		}
	})

	worker, trigger := worker.New(fmt.Sprintf("StoreKline-%s", interval.String()), func(k *models.Kline) {
		klinePipe <- k
	})

	worker.WithLogger(logger)
	worker.WithContext(ctx)
	worker.WithGraceful(true)
	worker.Run()
	return trigger
}
