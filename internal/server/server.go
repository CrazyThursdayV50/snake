package server

import (
	"context"
	"os"
	"os/signal"
	"snake/internal/kline"
	"snake/internal/kline/handler"
	"snake/internal/kline/interval"
	"snake/internal/kline/storage/mysql/migrate"
	"snake/internal/kline/storage/mysql/models"
	"snake/internal/kline/workers"
	"snake/pkg/binance"

	"github.com/CrazyThursdayV50/pkgo/log"
	defaultlogger "github.com/CrazyThursdayV50/pkgo/log/default"
	"github.com/CrazyThursdayV50/pkgo/store/db/gorm"
	jaeger "github.com/CrazyThursdayV50/pkgo/trace/jaeger"
	binance_connector "github.com/binance/binance-connector-go"
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
	clients  *Clients
	repos    *Repositories
	Workers  *Workers
	Handlers *Handlers
}

func New(cfg *Config) *Server {
	return &Server{cfg: cfg, clients: &Clients{}, repos: &Repositories{}, Workers: &Workers{}, Handlers: &Handlers{}}
}

func (s *Server) initClients() {
	loggerCfg := defaultlogger.DefaultConfig()
	logger := defaultlogger.New(loggerCfg)
	logger.Init()

	jaegerCfg := jaeger.DefaultConfig()
	tracer, err := jaeger.New(context.Background(), jaegerCfg, logger)
	if err != nil {
		panic(err)
	}

	s.logger = logger
	s.clients.db = gorm.NewDB(logger, tracer.NewTracer("mysql"), s.cfg.Mysql)
	s.clients.binanceMarket = binance.New(s.cfg.Binance)
}

func (s *Server) initWorkers(ctx context.Context) {
	s.Workers.StoreKlineTrigger = workers.NewIntervalTrigger[*models.Kline]()
	s.Workers.UptodateKlineTrigger = workers.NewIntervalTrigger[uint64]()
	s.Workers.CheckerTrigger = workers.NewIntervalTrigger[uint64]()
	s.Handlers.WsKlineEvent = handler.NewIntervalHandler[*handler.WsKlineHandler]()

	// var symbolIntervalMap = make(map[string]string)
	for _, in := range interval.All() {
		storeTrigger := kline.StoreKline(ctx, s.logger, in, s.repos.repoKline)
		s.Workers.StoreKlineTrigger.Add(in, storeTrigger)

		checkTrigger := kline.CheckKline(ctx, s.logger, s.cfg.Service.Symbol, in, s.repos.repoKline, s.clients.binanceMarket, storeTrigger)
		s.Workers.CheckerTrigger.Add(in, checkTrigger)

		uptodateTrigger := kline.UptodateKline(ctx, s.logger, s.cfg.Service.Symbol, in, s.repos.repoKline, s.clients.binanceMarket, storeTrigger, checkTrigger)
		s.Workers.UptodateKlineTrigger.Add(in, uptodateTrigger)

		// symbolIntervalMap[s.cfg.Service.Symbol] = in.String()
		// s.Handlers.WsKlineEvent.Add(
		// 	in.String(),
		// 	handler.NewWsKline(uptodateTrigger, storeTrigger))

		if in == interval.Min1() {
			handler := handler.NewWsKline(uptodateTrigger, storeTrigger)
			s.clients.binanceMarket.Stream.WsKlineServe(s.cfg.Service.Symbol, interval.Min1().String(), func(event *binance_connector.WsKlineEvent) {
				s.logger.Infof("kline event: %+#v", event)
				handler.Handle(event)
			}, func(err error) {
				s.logger.Error("get kline error: %v", err)
			})
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

	ctx := context.Background()
	migrate.AutoMigrate(ctx, s.clients.db)

	s.initWorkers(ctx)

	var signalChan = make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	s.logger.Warn("SERVER EXIT ...")
}
