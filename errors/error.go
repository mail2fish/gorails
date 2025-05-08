package errors

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// 格式为：错误类型-业务模块-错误编码
// ErrorType 定义错误类型
type ErrorType int

const (
	// HANDLER Handler层错误
	HANDLER ErrorType = 1
	// THIRD_PARTY 第三方库错误
	THIRD_PARTY ErrorType = 2
	// DAO 数据访问层错误
	DAO ErrorType = 3
)

// ErrorModule 定义业务模块
type ErrorModule int

// ErrorCode 定义错误代码
type ErrorCode string

// Error 自定义应用错误
type Error interface {
	error
	HTTPCode() int
	ErrorCode() ErrorCode
	Unwrap() error
	Render(*gin.Context)
}

// ErrorImp 自定义应用错误
type ErrorImp struct {
	httpCode int         // HTTP状态码
	typ      ErrorType   // 错误类型
	module   ErrorModule // 业务模块
	sn       int         // 错误序号
	message  string      // 错误消息
	err      error       // 原始错误
}

// Error 实现 error 接口
func (e *ErrorImp) Error() string {
	if e.err != nil {
		return fmt.Sprintf("[%s] %v", e.ErrorCode(), e.err)
	}
	return fmt.Sprintf("[%s]", e.ErrorCode())
}

// Unwrap 返回原始错误，支持 errors.Is 和 errors.As
func (e *ErrorImp) Unwrap() error {
	if e.err == nil {
		return fmt.Errorf("[%s] %s", e.ErrorCode(), e.message)
	}
	return e.err
}

// ErrorCode 获取错误代码
// 格式为：错误类型-业务模块-错误编码
func (e *ErrorImp) ErrorCode() ErrorCode {
	return ErrorCode(fmt.Sprintf("%d-%d-%d", e.typ, e.module, e.sn))
}

// HTTPCode 获取HTTP状态码
func (e *ErrorImp) HTTPCode() int {
	return e.httpCode
}

// Render 渲染错误
func (e *ErrorImp) Render(c *gin.Context) {
	c.JSON(e.HTTPCode(), gin.H{
		"code":    e.ErrorCode(),
		"message": e.message,
	})
}

// NewError 创建错误
func NewError(httpCode int, ty ErrorType, module ErrorModule, sn int, message string, err error) *ErrorImp {
	return &ErrorImp{
		httpCode: httpCode,
		typ:      ty,
		module:   module,
		sn:       sn,
		message:  message,
		err:      err,
	}
}
