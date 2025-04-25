package service

import "github.com/gin-gonic/gin"

// Kline K线数据服务接口
type Kline interface {
	// 健康检查
	Ping(*gin.Context)

	// 获取K线数据
	GetKlines(*gin.Context)
	// 订阅K线数据
	Subscribe(*gin.Context)
}

// StrategyServiceRepository 策略服务接口
type StrategyServiceRepository interface {
	// 健康检查
	Ping(*gin.Context)

	// 策略管理
	ListStrategies(*gin.Context)
	CreateStrategy(*gin.Context)
	GetStrategy(*gin.Context)
	UpdateStrategy(*gin.Context)
	DeleteStrategy(*gin.Context)

	// 回测
	RunBacktest(*gin.Context)
}
