package api

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/MrKrisYu/koi-go-common/logger"
	"github.com/MrKrisYu/koi-go-common/sdk"
	"github.com/MrKrisYu/koi-go-common/sdk/api/header"
	"github.com/MrKrisYu/koi-go-common/sdk/api/response"
	"github.com/MrKrisYu/koi-go-common/sdk/i18n"
	"github.com/MrKrisYu/koi-go-common/sdk/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"golang.org/x/text/language"
	"io"
	"net/http"
	"reflect"
	"runtime"
)

var (
	EmptyData = struct{}{}
)

type Api struct {
	Logger *logger.Helper
}

// MakeContext 设置http上下文
func (e *Api) MakeContext(c *gin.Context) *Api {
	e.Logger = GetRequestLogger(c)
	return e
}

func (e *Api) MakeService(s *service.Service) *Api {
	s.Logger = e.Logger
	return e
}

// Bind 参数校验
func (e *Api) Bind(ctx *gin.Context, d any, bindings ...binding.Binding) error {
	var mergedErr error
	if len(bindings) == 0 {
		bindings = constructor.GetBindingForGin(d)
	}
	for i := range bindings {
		// 推荐每个请求体仅用一种绑定方式，否则会出现诸如结构体同时含有请求头和请求参数的必填字段，
		// 在绑定任一方式时，会导致另一个类型字段因未绑定而出现误判定为未填写的情况。
		var err error
		if bindings[i] == nil {
			err = ctx.ShouldBindUri(d)
		} else {
			err = ctx.ShouldBindWith(d, bindings[i])
		}
		if errors.Is(err, io.EOF) {
			e.Logger.Warn("request body is not present anymore. ")
			err = nil
			break
		}
		if err != nil {
			return err
		}
	}
	return mergedErr
}

func Translate(ctx *gin.Context, message i18n.Message) string {
	tag, exist := ctx.Get(header.AcceptLanguageFlag)
	if !exist {
		return message.DefaultMessage
	}
	lang := tag.(language.Tag)
	translator := sdk.RuntimeContext.GetTranslator()
	if translator == nil {
		return message.DefaultMessage
	}
	if len(message.ID) == 0 {
		return message.DefaultMessage
	}
	if message.Args == nil {
		return translator.Tr(lang, message)
	} else {
		return translator.TrWithData(lang, message)
	}
}

func (e *Api) OK(ctx *gin.Context, data any) {
	message := Translate(ctx, response.OK.Message)
	if data == nil {
		ctx.JSON(http.StatusOK, response.Response[any]{
			Code:    response.OK.Code,
			Message: message,
			Data:    EmptyData,
		})
		return
	}
	var ret any
	value := reflect.ValueOf(data)
	switch value.Kind() {
	case reflect.Slice:
		// 检查切片是否为空
		if value.IsNil() {
			// 创建一个非nil的空切片
			data = reflect.MakeSlice(value.Type(), 0, 0).Interface()
		}
		// 切片，使用ListDataWrapper封装
		ret = response.Response[ListDataWrapper[any]]{
			Code:    response.OK.Code,
			Message: message,
			Data:    ListDataWrapper[any]{Items: data},
		}
		break
	case reflect.Struct:
		// 结构体，直接返回
		ret = response.Response[any]{
			Code:    response.OK.Code,
			Message: message,
			Data:    data,
		}
		break
	default:
		// 非切片且非结构体，使用ValueWrapper封装
		ret = response.Response[ValueWrapper[any]]{
			Code:    response.OK.Code,
			Message: message,
			Data:    ValueWrapper[any]{Value: data},
		}
	}
	ctx.JSON(http.StatusOK, ret)

}

func (e *Api) Error(ctx *gin.Context, businessStatus response.Status, err error) {
	var msg string
	var myError i18n.MyError
	if !errors.As(err, &myError) { // 处理未知错误
		// 打印未知错误
		requestLogger := GetRequestLogger(ctx)
		requestLogger.Error(err)
		// 翻译响应状态码
		msg = Translate(ctx, businessStatus.Message)
		ctx.JSON(http.StatusOK, response.Response[struct{}]{
			Code:    businessStatus.Code,
			Message: msg,
			Data:    EmptyData,
		})
		return
	}
	// 若是内部错误类型，则获取翻译后的错误信息即可
	msg = Translate(ctx, myError.Message)
	ctx.JSON(http.StatusOK, response.Response[struct{}]{
		Code:    businessStatus.Code,
		Message: msg,
		Data:    EmptyData,
	})
}

// ListDataWrapper 用于包装切片
type ListDataWrapper[T any] struct {
	Items T `json:"items"` // 列表数据
}

// ValueWrapper 用于包装非切片且非结构体的数据类型
type ValueWrapper[T any] struct {
	Value T `json:"value"` // 结构体数据
}

// LogStack return call function stack info from start stack to end stack.
// if end is a positive number, return all call function stack.
func LogStack(start, end int) string {
	stack := bytes.Buffer{}
	for i := start; i < end || end <= 0; i++ {
		pc, str, line, _ := runtime.Caller(i)
		if line == 0 {
			break
		}
		stack.WriteString(fmt.Sprintf("%s:%d %s\n", str, line, runtime.FuncForPC(pc).Name()))
	}
	return stack.String()
}
