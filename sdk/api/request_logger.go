package api

import (
	"github.com/MrKrisYu/koi-go-common/logger"
	"github.com/MrKrisYu/koi-go-common/sdk"
	"github.com/gin-gonic/gin"
)

const (
	TrafficKey = "X-Request-Id"
	//LoggerKey  = "logger-request"
)

// GetRequestLogger 获取上下文提供的日志
func GetRequestLogger(c *gin.Context) *logger.Helper {
	return sdk.RuntimeContext.GetLogger()
}
