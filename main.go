package main

import (
	"flag"
	"os"
	"snake/internal/server"

	"github.com/CrazyThursdayV50/pkgo/config"
)

var cfgDir string
var cfgName string
func init(){
	flag.StringVar(&cfgDir, "d", ".", "配置所在目录")
	flag.StringVar(&cfgName, "c", "config", "配置文件名（没有扩展名）")
}

func main() {
	flag.Parse()
	cfg,err :=config.GetConfig[server.Config](cfgDir	, cfgName, "yml")	
	if err!=nil{panic(err)}

	api, ok := os.LookupEnv("BN_APIKEY")
	if !ok {
		panic("invalid api")
	}
	
	secret, ok := os.LookupEnv("BN_SECRET")
	if !ok {
		panic("invalid secret")
	}
	
	cfg.Binance.APIKey = api
	cfg.Binance.SecretKey = secret

	server := server.New(cfg)
	server.Run()
}