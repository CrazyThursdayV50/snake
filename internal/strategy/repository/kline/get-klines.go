package kline

import (
	"context"
	"fmt"
	k "snake/internal/kline"
	"snake/internal/kline/interval"
	"snake/internal/kline/utils"
	"snake/internal/service/kline"
	"time"

	"github.com/CrazyThursdayV50/pkgo/builtin/slice"
	"github.com/CrazyThursdayV50/pkgo/goo"
	"github.com/CrazyThursdayV50/pkgo/log"
	"github.com/CrazyThursdayV50/pkgo/websocket/client"
	"github.com/gorilla/websocket"
)

type Kline = k.Kline

func (r *Repository) GetKlines(ctx context.Context, interval interval.Interval, from int64) <-chan *Kline {
	var ch = make(chan *k.Kline, 100)

	var klineInited bool
	client := client.New(
		client.WithURL(r.wsEndpoint),
		client.WithContext(ctx),
		client.WithPingLoop(func(done <-chan struct{}, conn *websocket.Conn) {
			var t = time.NewTicker(time.Second * 10)
			for {
				select {
				case <-done:
					return
				case <-t.C:
					conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(time.Second*10))
				}
			}
		}),
		client.WithLogger(r.logger),
		client.WithMessageHandler(func(ctx context.Context, logger log.Logger, _ int, data []byte, f func(error)) (int, []byte) {
			logger.Infof("kline: %s", data)
			var line Kline
			err := line.UnmarshalBinary(data)
			if err != nil {
				f(err)
				return client.BinaryMessage, nil
			}

			if !klineInited {
				to := utils.GetLastTime(uint64(line.S), interval)
				var result kline.GetKlineReponse
				_, err := r.client.Request(ctx).SetQueryParams(map[string]string{
					"interval": interval.String(),
					"from":     fmt.Sprintf("%d", from),
					"to":       fmt.Sprintf("%d", to),
				}).
					SetResult(&result).
					SetError(&result).
					Get(r.endpoint)

				if err != nil {
					logger.Errorf("get klines failed: %v", err)
					return client.BinaryMessage, nil
				}

				if result.Error != "" {
					logger.Errorf("get klines failed: %v: %v", result.Message, result.Error)
					return client.BinaryMessage, nil
				}

				slice.From(result.Data.List...).Iter(func(k int, v *k.Kline) (bool, error) {
					ch <- v
					return true, nil
				})

				klineInited = true
			}

			ch <- &line
			return client.BinaryMessage, nil
		}),
	)

	goo.Go(func() {
		<-ctx.Done()
		fmt.Printf("receive exit\n")
		client.Stop()
		close(ch)
	})

	client.Run()
	return ch
}
