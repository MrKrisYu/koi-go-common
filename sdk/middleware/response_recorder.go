package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

// ResponseRecorder 是一个包装了gin.ResponseWriter的结构体，用于记录响应状态码和响应体
type ResponseRecorder struct {
	gin.ResponseWriter
	Body *bytes.Buffer
	cnt  int64
}

func NewResponseRecorder(w gin.ResponseWriter) *ResponseRecorder {
	return &ResponseRecorder{w, &bytes.Buffer{}, 0}
}

// 写入方法`Write`同时将数据写入自己的`Body`缓存和原始的`ResponseWriter`中
func (r *ResponseRecorder) Write(b []byte) (int, error) {
	contentType := r.ResponseWriter.Header().Get("Content-Type")
	fmt.Println("Response Content-Type: ", contentType)
	if strings.Contains(contentType, "application/json") {
		r.Body.Write(b)
	}
	r.cnt += int64(len(b))
	return r.ResponseWriter.Write(b)
}
