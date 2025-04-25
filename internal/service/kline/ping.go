package kline

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping 健康检查接口
func (s *Service) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "K线服务正常运行",
	})
}
