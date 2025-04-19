package workers

import (
	"context"
	"fmt"
	"snake/internal/kline/acl"
	"snake/internal/kline/interval"
	"snake/internal/kline/storage/mysql/models"
	"snake/internal/kline/utils"
	"snake/internal/repository"
	"snake/pkg/binance"

	"github.com/CrazyThursdayV50/pkgo/builtin/slice"
	"github.com/CrazyThursdayV50/pkgo/log"
	"github.com/CrazyThursdayV50/pkgo/worker"
	binance_connector "github.com/binance/binance-connector-go"
)

const updateKlinesCount = 1000

func updateKlineFromStartTime(
	ctx context.Context,
	client *binance.MarketClient,
	logger log.Logger,
	symbol string,
	interval interval.Interval,
	startTime uint64,
	stopTime uint64,
	trigger func(*models.Kline),
) {
	for {
		nextStart := utils.GetNextTime(startTime, interval)
		if nextStart >= stopTime {
			return
		}

		endTime := utils.GetEndTimeByStartTime(nextStart, interval, updateKlinesCount)
		if endTime >= stopTime {
			endTime = utils.GetEndTimeByStartTime(stopTime, interval, 0)
		}

		if endTime <= startTime {
			return
		}

		resp, err := client.Restful.NewKlinesService().StartTime(nextStart).EndTime(endTime).Interval(interval.String()).Symbol(symbol).Limit(updateKlinesCount).Do(ctx)
		if err != nil {
			logger.Errorf("Failed to fetch klines: %v", err)
			return
		}

		slice.From(resp...).Iter(func(k int, v *binance_connector.KlinesResponse) (bool, error) {
			model := acl.ApiToDB(v)
			trigger(model)
			startTime = v.OpenTime
			return true, nil
		})
	}
}

func updateKlineToStartTime(
	ctx context.Context,
	client *binance.MarketClient,
	logger log.Logger,
	symbol string,
	interval interval.Interval,
	startTime uint64,
	trigger func(*models.Kline),
) {
	for {
		endTime := utils.GetEndTimeByStartTime(startTime, interval, 0)
		nextStartTime := utils.GetStartTimeByEndTime(endTime, interval, updateKlinesCount)

		if endTime <= nextStartTime {
			return
		}

		resp, err := client.Restful.NewKlinesService().
			StartTime(nextStartTime).
			EndTime(endTime).
			Interval(interval.String()).
			Symbol(symbol).
			Limit(updateKlinesCount).
			Do(ctx)
		if err != nil {
			logger.Errorf("Failed to fetch klines: %v", err)
			return
		}

		if len(resp) == 0 {
			return
		}

		slice.From(resp...).Iter(func(k int, v *binance_connector.KlinesResponse) (bool, error) {
			model := acl.ApiToDB(v)
			trigger(model)
			return true, nil
		})

		startTime = resp[0].OpenTime
	}
}

func UptodateKline(
	ctx context.Context,
	logger log.Logger,
	symbol string,
	interval interval.Interval,
	repoKline repository.KlineRepository,
	marketClient *binance.MarketClient,
	storeTrigger func(*models.Kline),
	checkTrigger func(uint64),
) func(uint64) {
	worker, trigger := worker.New(fmt.Sprintf("UptodateKline-%s", interval.String()), func(stopTimestamp uint64) {
		tryFunc(func() error {
			last, err := repoKline.Last(ctx, interval)
			if err != nil {
				return err
			}

			if last != nil {
				updateKlineFromStartTime(ctx, marketClient, logger, symbol, interval, uint64(last.OpenTs), stopTimestamp, storeTrigger)

				first, err := repoKline.First(ctx, interval)
				if err != nil {
					return err
				}

				updateKlineToStartTime(ctx, marketClient, logger, symbol, interval, uint64(first.OpenTs), storeTrigger)
				checkTrigger(uint64(stopTimestamp))
				return nil
			}

			updateKlineToStartTime(ctx, marketClient, logger, symbol, interval, uint64(stopTimestamp), storeTrigger)
			checkTrigger(uint64(stopTimestamp))
			return nil
		}, func(err error) {
			logger.Errorf("Uptodate kline error: %v", err)
		}, 3)
	})

	worker.WithContext(ctx)
	worker.WithLogger(logger)
	worker.WithGraceful(true)
	worker.Run()
	return trigger
}
