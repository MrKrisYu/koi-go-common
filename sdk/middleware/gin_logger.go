package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/MrKrisYu/koi-go-common/sdk/api"
	"github.com/gin-gonic/gin"
	"io"
	"reflect"
	"runtime"
	"strings"
	"time"
)

func GinLogger(trafficKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestLogger := api.GetRequestLogger(c)
		// 获取请求地址和方法
		requestPath := c.Request.URL.String()
		requestMethod := c.Request.Method
		responseStatus := c.Writer.Status()
		var reqParams interface{}

		// 读取请求体
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
		}
		// 将读取过的请求体重新放入请求中
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// 创建响应记录器
		recorder := NewResponseRecorder(c.Writer)
		c.Writer = recorder
		// 处理请求前记录时间
		start := time.Now()
		// 处理请求
		c.Next()
		// 处理请求后记录时间
		cost := time.Since(start)

		// 请求的后置处理
		switch c.ContentType() { // 根据 Content-Type 选择合适的方式来获取请求参数
		case "application/json":
			var jsonMap map[string]interface{}
			if err := json.Unmarshal(bodyBytes, &jsonMap); err == nil {
				reqParams = jsonMap
			}
		case "application/x-www-form-urlencoded", "multipart/form-data":
			reqParams = c.Request.Form
		case "application/octet-stream":
			reqParams = "[BINARY DATA]"
		default:
			reqParams = fmt.Sprintf("[UNSUPPORTED CONTENT TYPE: %s]", c.ContentType()) //其他类型数据暂不做处理
		}

		// 获取响应数据
		var responseData string
		respContentType := recorder.ResponseWriter.Header().Get("Content-Type")
		// 解析 Content-Type，移除可选参数
		contentTypeWithoutParams := strings.SplitN(respContentType, ";", 2)[0]
		// 请求的后置处理
		switch contentTypeWithoutParams {
		case "application/json":
			// 处理所有JSON类型的响应，忽略字符集
			responseData = recorder.Body.String()
		case "application/x-www-form-urlencoded", "multipart/form-data":
			// 处理form类型的响应
			responseData = recorder.Body.String()
		case "audio/mpeg", "text/html", "application/octet-stream":
			// 忽略记录静态资源或二进制数据的响应
			responseData = "[BINARY DATA]"
		default:
			// 对于未受支持的 Content-Type，不记录详细的请求参数或响应数据
			responseData = fmt.Sprintf("[UNSUPPORTED CONTENT TYPE: %s]", contentTypeWithoutParams)
		}

		clientIP := c.ClientIP()
		// 获取 handler 名称
		handlerName := runtime.FuncForPC(reflect.ValueOf(c.Handler()).Pointer()).Name()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()
		requestId := c.Request.Header.Get(trafficKey)

		requestLogger.Info("---------------------------请求开始-----------------------------")
		requestLogger.Infof("X-Request-Id: %s", requestId)
		requestLogger.Infof("CLASS METHOD: %s", handlerName)
		requestLogger.Infof("请求地址: %s", requestPath)
		requestLogger.Infof("请求参数: %+v", reqParams)
		requestLogger.Infof("HTTP METHOD: %s", requestMethod)
		requestLogger.Infof("IP: %s", clientIP)
		requestLogger.Infof("HTTP STATUS: %d", responseStatus)
		requestLogger.Infof("Error Message: %s", errorMessage)
		requestLogger.Infof("响应数据: %s", responseData)
		requestLogger.Infof("耗时: %s", cost.String())
		requestLogger.Info("---------------------------请求结束-----------------------------\n")
	}
}
