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
		requestId := c.Request.Header.Get(trafficKey)
		// 获取 handler 名称
		handlerName := runtime.FuncForPC(reflect.ValueOf(c.Handler()).Pointer()).Name()
		// 请求来源的IP地址
		clientIP := c.ClientIP()
		// 获取请求地址和方法
		requestPath := c.Request.URL.String()
		requestMethod := c.Request.Method
		var reqParams interface{}
		// 处理请求体
		switch c.ContentType() { // 根据 Content-Type 选择合适的方式来获取请求参数
		case "application/json":
			// 读取请求体
			var bodyBytes []byte
			if c.Request.Body != nil {
				bodyBytes, _ = io.ReadAll(c.Request.Body)
			}
			// 将读取过的请求体重新放入请求中
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			// 反序列化请求体
			var jsonMap map[string]interface{}
			if err := json.Unmarshal(bodyBytes, &jsonMap); err == nil {
				if marshal, err := json.Marshal(jsonMap); err == nil {
					reqParams = string(marshal)
				}
			}
		case "application/x-www-form-urlencoded", "multipart/form-data":
			reqParams = c.Request.Form
		case "application/octet-stream":
			reqParams = "[BINARY DATA]"
		default:
			reqParams = fmt.Sprintf("[UNSUPPORTED CONTENT TYPE: %s]", c.ContentType()) //其他类型数据暂不做处理
		}
		// 创建响应记录器
		recorder := NewResponseRecorder(c.Writer)
		c.Writer = recorder

		// 处理请求前记录时间
		start := time.Now()

		logTemplate := "X-Request-Id:%s" +
			"\n---------------------------请求开始-----------------------------\n" +
			"CLASS METHOD: %s\n" +
			"请求地址: %s\n" +
			"请求参数: %+v\n" +
			"HTTP METHOD: %s\n" +
			"IP: %s\n"
		logStr := fmt.Sprintf(logTemplate,
			requestId,
			handlerName,
			requestPath,
			reqParams,
			requestMethod,
			clientIP,
		)
		requestLogger.Info(logStr)

		// 处理请求
		c.Next()
		// 处理请求后记录时间
		cost := time.Since(start)

		// 响应状态
		responseStatus := c.Writer.Status()
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

		// 响应错误信息
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		responseTemplate := "X-Request-Id:%s\n" +
			"HTTP STATUS: %d\n" +
			"Error Messages: %s\n" +
			"响应数据: %s\n" +
			"响应大小: %d\n" +
			"耗时: %s\n" +
			"---------------------------请求结束-----------------------------\n"
		respLogStr := fmt.Sprintf(responseTemplate,
			requestId,
			responseStatus,
			errorMessage,
			responseData,
			recorder.Body.Len(),
			cost.String(),
		)
		requestLogger.Info(respLogStr)
	}
}
