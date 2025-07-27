package workers

import (
	"context"
	"fmt"
	"snake/internal/kline"
	"snake/internal/kline/acl"
	"snake/internal/kline/interval"
	"snake/internal/kline/storage/mysql/models"
	"snake/internal/kline/utils"
	"snake/pkg/binance"
	"sort"

	"github.com/CrazyThursdayV50/goex/binance/websocket-api/models/klines"
	"github.com/CrazyThursdayV50/pkgo/builtin/collector"
	gmap "github.com/CrazyThursdayV50/pkgo/builtin/map"
	"github.com/CrazyThursdayV50/pkgo/builtin/slice"
	"github.com/CrazyThursdayV50/pkgo/log"
	"github.com/CrazyThursdayV50/pkgo/worker"
)

func gatherOpenTs(openTs []uint64, interval interval.Interval) map[uint64]int {
	if len(openTs) == 0 {
		return nil
	}

	sli := slice.From(openTs...)
	sli.WithLessFunc(func(a uint64, b uint64) bool { return a < b })
	sort.Sort(sli)

	n := (sli.Get(sli.Len()-1).Unwrap() - sli.Get(0).Unwrap()) / uint64(interval.Duration().Milliseconds())

	tsMap := collector.Map(openTs, func(_ int, v uint64) (bool, uint64, struct{}) { return true, v, struct{}{} })

	var k uint64
	var next = sli.Get(0).Unwrap()
	var paramsMap = make(map[uint64]int)
	for range n {
		_, ok := tsMap[next]
		if ok {
			if k == 0 {
				k = next
			}

			paramsMap[k]++
			next = utils.GetNextTime(next, interval)
			continue
		}

		k = 0
	}

	return paramsMap
}

func Checker(
	ctx context.Context,
	logger log.Logger,
	symbol string,
	interval interval.Interval,
	repoKline kline.Repository,
	marketClient *binance.MarketClient,
	storeTrigger func(*models.Kline),
) func(uint64) {
	worker, trigger := worker.New(fmt.Sprintf("KlineChecks-%s", interval.String()), func(stopTime uint64) {
		logger.Infof("check to %d", stopTime)
		tryFunc(func() error {
			first, err := repoKline.First(ctx, interval)
			if err != nil {
				logger.Errorf("failed to get first kline: %v", err)
				return err
			}

			updateKlineToStartTime(ctx, marketClient, logger, symbol, interval, uint64(first.OpenTs), storeTrigger)

			first, err = repoKline.First(ctx, interval)
			if err != nil {
				logger.Errorf("failed to get first kline: %v", err)
				return err
			}

			startTime := first.OpenTs

			for {
				tsRange := utils.GenNextTimeToN(uint64(startTime), stopTime, interval, 10000)

				missingTs, err := repoKline.CheckMissing(ctx, interval, tsRange)
				if err != nil {
					logger.Errorf("check missing klines failed: %v", err)
					return err
				}

				params := gatherOpenTs(missingTs, interval)
				_, err = gmap.From(params).Iter(func(k uint64, v int) (bool, error) {
					endTime := utils.GetEndTimeByStartTime(k, interval, int64(v))
					resp, err := marketClient.Restful.Klines().
						Symbol(symbol).
						StartTime(k).
						EndTime(endTime).
						Interval(interval.String()).Do(ctx)
					if err != nil {
						logger.Errorf("request for klines failed: %v", err)
						return false, err
					}

					klinesData := resp.Unwrap()
					slice.From(klinesData...).Iter(func(k int, v klines.Kline) (bool, error) {
						model := acl.ApiToDB(v)
						storeTrigger(model)
						return true, nil
					})

					return true, nil
				})

				if err != nil {
					return err
				}

				if len(tsRange) < 10000 {
					return nil
				}

				startTime = tsRange[len(tsRange)-1]
				startTime = int64(utils.GetNextTime(uint64(startTime), interval))
			}
		},
			func(err error) { logger.Errorf("check klines failed: %v", err) },
			3)
	})

	worker.WithContext(ctx)
	worker.WithLogger(logger)
	worker.Run()
	return trigger
}
