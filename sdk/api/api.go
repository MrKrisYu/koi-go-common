package api

import (
	"errors"
	"fmt"
	"github.com/MrKrisYu/koi-go-common/logger"
	"github.com/MrKrisYu/koi-go-common/sdk/api/response"
	"github.com/MrKrisYu/koi-go-common/sdk/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io"
	"net/http"
	"reflect"
	"strings"
)

var (
	EmptyData = struct{}{}
)

type Api struct {
	//Context *gin.Context
	Logger *logger.Helper
	//Errors error
}

//func (e *Api) AddError(err error) {
//	if e.Errors == nil {
//		e.Errors = err
//	} else if err != nil {
//		e.Logger.Error(err)
//		e.Errors = fmt.Errorf("%v; %w", e.Errors, err)
//	}
//}

// MakeContext 设置http上下文
func (e *Api) MakeContext(c *gin.Context) *Api {
	//e.Context = c
	e.Logger = GetRequestLogger(c)
	return e
}

func (e *Api) MakeService(s *service.Service) *Api {
	s.Logger = e.Logger
	return e
}

// GetLogger 获取上下文提供的日志
//func (e *Api) GetLogger() *logger.Helper {
//	return GetRequestLogger(e.Context)
//}

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
			mergedErr = fmt.Errorf("%v; %w", mergedErr, err)
			continue
		}
	}
	return mergedErr
}

func (e *Api) OK(ctx *gin.Context, data any) {
	if data == nil {
		ctx.JSON(http.StatusOK, response.Response[any]{
			Code:    response.OK.Code,
			Message: response.OK.Message,
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
			Message: response.OK.Message,
			Data:    ListDataWrapper[any]{Items: data},
		}
		break
	case reflect.Struct:
		// 结构体，直接返回
		ret = response.Response[any]{
			Code:    response.OK.Code,
			Message: response.OK.Message,
			Data:    data,
		}
		break
	default:
		// 非切片且非结构体，使用ValueWrapper封装
		ret = response.Response[ValueWrapper[any]]{
			Code:    response.OK.Code,
			Message: response.OK.Message,
			Data:    ValueWrapper[any]{Value: data},
		}
	}
	ctx.JSON(http.StatusOK, ret)
	// 一个请求事务完结后，把错误清空，避免错误过度传递，影响下个请求事务
	//e.Errors = nil
}

func (e *Api) Error(ctx *gin.Context, businessStatus response.Status, errMsg ...string) {
	errMessage := businessStatus.Message
	if len(errMsg) > 0 {
		errMessage = strings.Join([]string{errMessage, errMsg[0]}, ": ")
	}
	ctx.JSON(http.StatusOK, response.Response[struct{}]{
		Code:    businessStatus.Code,
		Message: errMessage,
		Data:    EmptyData,
	})
	// 一个请求事务完结后，把错误清空，避免错误过度传递，影响下个请求事务
	//e.Errors = nil
}

func (e *Api) Translate(from, to any) {
	CopyProperties(from, to)
}

// ListDataWrapper 用于包装切片
type ListDataWrapper[T any] struct {
	Items T `json:"items"` // 列表数据
}

// ValueWrapper 用于包装非切片且非结构体的数据类型
type ValueWrapper[T any] struct {
	Value T `json:"value"` // 结构体数据
}
