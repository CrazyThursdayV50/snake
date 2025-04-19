package main

import (
	"flag"
	"snake/internal/server/kline"

	"github.com/CrazyThursdayV50/pkgo/config"
)

var cfgDir string
var cfgName string

func init() {
	flag.StringVar(&cfgDir, "d", ".", "配置所在目录")
	flag.StringVar(&cfgName, "c", "config", "配置文件名（没有扩展名）")
}

func main() {
	flag.Parse()
	cfg, err := config.GetConfig[kline.Config](cfgDir, cfgName, "yml")
	if err != nil {
		panic(err)
	}

	cfg.Binance.UpdateKeys()
	server := kline.New(cfg)
	server.Run()
}
