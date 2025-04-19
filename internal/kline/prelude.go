package kline

import (
	"snake/internal/kline/interval"
	"snake/internal/kline/repository"
	"snake/internal/kline/workers"
)

type Interval = interval.Interval

var StoreKline = workers.StoreKline
var UptodateKline = workers.UptodateKline
var CheckKline = workers.Checker
var NewRepository = repository.New
