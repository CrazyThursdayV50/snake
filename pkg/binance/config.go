package binance

import "os"

type Config struct {
	APIKey    string
	SecretKey string
	Symbol    string
}

func (c *Config) UpdateKeys() {
	c.APIKey = os.Getenv("BN_APIKEY")
	c.SecretKey = os.Getenv("BN_SECRET")
}
