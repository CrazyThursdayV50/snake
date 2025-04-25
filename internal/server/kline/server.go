package kline

import (
	"context"
	"os"
	"os/signal"
	"snake/internal/kline/acl"
	"snake/internal/kline/handler"
	"snake/internal/kline/interval"
	"snake/internal/kline/storage/mysql/migrate"
	"snake/internal/kline/storage/mysql/models"
	"snake/internal/kline/workers"
	"snake/pkg/binance"
	"sync"

	"github.com/CrazyThursdayV50/pkgo/goo"
	"github.com/CrazyThursdayV50/pkgo/json"
	"github.com/CrazyThursdayV50/pkgo/log"
	defaultlogger "github.com/CrazyThursdayV50/pkgo/log/default"
	"github.com/CrazyThursdayV50/pkgo/store/db/gorm"
	"github.com/CrazyThursdayV50/pkgo/trace"
	jaeger "github.com/CrazyThursdayV50/pkgo/trace/jaeger"
	"github.com/CrazyThursdayV50/pkgo/websocket/server"
	binance_connector "github.com/binance/binance-connector-go"
	"github.com/gorilla/websocket"
	jsoniter "github.com/json-iterator/go"
)

type Clients struct {
	db            *gorm.DB
	binanceMarket *binance.MarketClient
}

type Workers struct {
	StoreKlineTrigger    *workers.IntervalTrigger[*models.Kline]
	UptodateKlineTrigger *workers.IntervalTrigger[uint64]
	CheckerTrigger       *workers.IntervalTrigger[uint64]
}

type Handlers struct {
	WsKlineEvent *handler.IntervalHandler[*handler.WsKlineHandler]
}

type Server struct {
	cfg      *Config
	logger   log.Logger
	tracer   trace.TracerCreator
	clients  *Clients
	repos    *Repositories
	Workers  *Workers
	Handlers *Handlers
	Services *Services
	Wsserver *server.Server

	min1Uptodater func(uint64)
	mint1Storer   func(*models.Kline)
}

func New(cfg *Config) *Server {
	return &Server{cfg: cfg, clients: &Clients{}, repos: &Repositories{}, Workers: &Workers{}, Handlers: &Handlers{}, Services: &Services{}}
}

func (s *Server) initClients() {
	json.Init(&jsoniter.Config{
		EscapeHTML:    true,
		UseNumber:     true,
		SortMapKeys:   true,
		CaseSensitive: true,
	})

	logger := defaultlogger.New(s.cfg.Log)
	logger.Init()

	jaegerCfg := jaeger.DefaultConfig()
	tracer, err := jaeger.New(context.Background(), jaegerCfg, logger)
	if err != nil {
		panic(err)
	}

	s.tracer = tracer
	s.logger = logger
	s.clients.db = gorm.NewDB(logger, tracer.NewTracer("mysql"), s.cfg.Mysql)
	s.clients.binanceMarket = binance.New(s.cfg.Binance)
}

func (s *Server) initWorkers(ctx context.Context) {
	s.Workers.StoreKlineTrigger = workers.NewIntervalTrigger[*models.Kline]()
	s.Workers.UptodateKlineTrigger = workers.NewIntervalTrigger[uint64]()
	s.Workers.CheckerTrigger = workers.NewIntervalTrigger[uint64]()
	s.Handlers.WsKlineEvent = handler.NewIntervalHandler[*handler.WsKlineHandler]()

	symbol := s.cfg.Binance.Symbol

	// var symbolIntervalMap = make(map[string]string)
	for _, in := range interval.All() {
		storeTrigger := workers.StoreKline(ctx, s.logger, in, s.repos.repoKline)
		s.Workers.StoreKlineTrigger.Add(in, storeTrigger)

		checkTrigger := workers.Checker(ctx, s.logger, symbol, in, s.repos.repoKline, s.clients.binanceMarket, storeTrigger)
		s.Workers.CheckerTrigger.Add(in, checkTrigger)

		uptodateTrigger := workers.UptodateKline(ctx, s.logger, symbol, in, s.repos.repoKline, s.clients.binanceMarket, storeTrigger, checkTrigger)
		s.Workers.UptodateKlineTrigger.Add(in, uptodateTrigger)

		// symbolIntervalMap[s.cfg.Service.Symbol] = in.String()
		// s.Handlers.WsKlineEvent.Add(
		// 	in.String(),
		// 	handler.NewWsKline(uptodateTrigger, storeTrigger))

		if in == interval.Min1() {
			s.min1Uptodater = uptodateTrigger
			s.mint1Storer = storeTrigger
		}
	}

	// goo.Goo(func() {
	// 	s.clients.binanceMarket.Stream.WsCombinedKlineServe(symbolIntervalMap, func(event *binance_connector.WsKlineEvent) {
	// 		s.logger.Infof("event: %+v", event)
	// 		h, ok := s.Handlers.WsKlineEvent.Get(event.Kline.Interval)
	// 		if ok {
	// 			h.Handle(event)
	// 		}
	// 	}, func(err error) {
	// 		s.logger.Error("get kline error: %v", err)
	// 	})
	// }, func(err error) {
	// 	s.logger.Error("ws server failed: %v", err)
	// })
}

func (s *Server) Run() {
	s.initClients()
	s.initRepositories()
	s.initServices()

	ctx, cancel := context.WithCancel(context.Background())
	migrate.AutoMigrate(ctx, s.clients.db)

	s.initWorkers(ctx)

	handler := handler.NewWsKline(s.min1Uptodater, s.mint1Storer)

	var done = new(chan struct{})
	var stop = new(chan struct{})
	var err error
	var startKline = func() {
		*done, *stop, err = s.clients.binanceMarket.Stream.WsKlineServe(s.cfg.Binance.Symbol, interval.Min1().String(), func(event *binance_connector.WsKlineEvent) {
			s.logger.Infof("kline event: %+#v", event)
			handler.Handle(event)
			_, kline := acl.Ws2Service(0, event)
			data, _ := kline.MarshalBinary()
			s.Wsserver.Broadcast(ctx, websocket.TextMessage, data)
		}, func(err error) {
			s.logger.Error("get kline error: %v", err)
		})
	}

	startKline()
	if err != nil {
		s.logger.Errorf("connect binance failed")
		panic(err)
	}

	goo.Go(func() {
		for {
			select {
			case <-*done:
				startKline()

			case <-ctx.Done():
				return
			}
		}

	})

	goo.Go(func() {
		<-ctx.Done()
		close(*stop)
	})

	var wg sync.WaitGroup
	s.Services.Run(ctx, s.cfg.Service, &wg)

	var signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	s.logger.Warn("SERVER EXIT ...")
	cancel()
	wg.Wait()
}
