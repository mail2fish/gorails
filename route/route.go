package route

import (
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/mail2fish/gorails/errors"
)

type Params interface {
	Parse(*gin.Context) errors.Error
}

type Response interface {
	Render(*gin.Context)
}

type GoRailsHandlerFun func(c *gin.Context, params Params) (Response, errors.Error)

type ResponseImpl struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Wrap 将 GoRailsHandlerFun 转换为 gin.HandlerFunc
func Wrap[T Params, R Response](handler GoRailsHandlerFun) gin.HandlerFunc {

	// 获取参数类型
	t := reflect.TypeOf(*new(T)).Elem()

	return func(c *gin.Context) {
		// 创建参数实例
		params := reflect.New(t).Interface().(T)

		// 解析参数
		err := params.Parse(c)
		if err != nil {
			err.Render(c)
			return
		}

		// 调用处理函数
		response, err := handler(c, params)

		// 处理错误
		if err != nil {
			err.Render(c)
			return
		}

		// 渲染响应
		response.Render(c)
	}
}
