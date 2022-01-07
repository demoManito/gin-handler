package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Error Codes
const (
	CodeSuccess           = 0
	CodeErrSignature      = 100010002
	CodeLoginRequired     = 100010006
	CodeUnknownError      = 100010007
	CodeBindMobile        = 100010008

	CodeBannedUser                = 200020004
	CodeUserLimit                 = 200020005
	CodeSendMessageUserCountLimit = 200020006
	CodeSendMessageLimit          = 200020008

	CodeSoundNotFound   = 200110001
	CodeAlbumNotFound   = 200120001
	CodeCommentNotFound = 200310001
	CodeEventNotFound   = 200410001
	CodeTagNotFound     = 200430001
	CodeCardNotFound    = 200120001
	CodeDanmakuNotFound = 200710001
	CodeDramaNotFound   = 200130001

	CodeEmptyParam   = 201010001
	CodeInvalidParam = 201010002
)

// errors
var (
	ErrLoginRequired = NewActionError(http.StatusForbidden, CodeLoginRequired, "请先登录")
	ErrBadRequest    = NewActionError(http.StatusBadRequest, CodeUnknownError, "错误的请求")
	ErrEmptyParam    = NewActionError(http.StatusBadRequest, CodeEmptyParam, "参数不可为空")
	ErrInvalidParam  = NewActionError(http.StatusBadRequest, CodeInvalidParam, "参数不合法")

	ErrEmptyValue = errors.New("empty value")
)

// ErrInvalidDateRange is returned by GetParamDateRange when start_date > end_date
var ErrInvalidDateRange = errors.New("start_date > end_date")

// ErrRawResponse is the special error that should be returned if the result of a action handler is not in the format of basicResponse.
var ErrRawResponse = errors.New("response is managed by action handler")

// ResponseError response error interface
type ResponseError interface {
	// StatusCode 响应状态码
	StatusCode() int
	// ErrorCode 错误代码
	ErrorCode() int
	// ErrorInfo 响应错误信息
	ErrorInfo() interface{}

	error
}

// ActionError for action error
type ActionError struct {
	Status  int
	Code    int
	Message string
	Info    M
}

// APIError
type APIError struct {
	Status  int
	Message string
}

// Error return error message
func (e *APIError) Error() string {
	return e.Message
}

// NewActionError returns a new ActionError
func NewActionError(status int, code int, msg string) *ActionError {
	return &ActionError{
		Status:  status,
		Code:    code,
		Message: msg,
	}
}

// NewActionErrorWithInfo returns a new ActionError with extra info
func NewActionErrorWithInfo(status int, code int, msg string, info M) *ActionError {
	return &ActionError{
		Status:  status,
		Code:    code,
		Message: msg,
		Info:    info,
	}
}

// StatusCode status code
func (e *ActionError) StatusCode() int {
	return e.Status
}

// ErrorCode error code
func (e *ActionError) ErrorCode() int {
	return e.Code
}

// ErrorInfo response info
// returns {"msg":"test err", e.Info...}
func (e *ActionError) ErrorInfo() interface{} {
	if e.Info != nil {
		info := make(M, len(e.Info)+1)
		for key, value := range e.Info {
			info[key] = value
		}
		// 客户端固定使用 msg 字段获取提示信息
		info["msg"] = e.Error()
		return info
	}
	return e.Error()
}

func (e *ActionError) Error() string {
	return e.Message
}

// LoggerError action error with logger
type LoggerError struct {
	Status int
	Code   int
	*ContextError
}

// NewLoggerError new LoggerError
func NewLoggerError(status, code int, msg string) *LoggerError {
	return &LoggerError{
		Status:       status,
		Code:         code,
		ContextError: NewContextError(msg),
	}
}

// StatusCode status code
func (e *LoggerError) StatusCode() int {
	return e.Status
}

// ErrorCode error code
func (e *LoggerError) ErrorCode() int {
	return e.Code
}

// ErrorInfo response info
func (e *LoggerError) ErrorInfo() interface{} {
	if Mode() == ReleaseMode {
		return e.Title()
	}
	return e.Error()
}

type basicResponse struct {
	Code int            `json:"code"`
	Info ActionResponse `json:"info"`
}

// Action structure
type Action struct {
	Method        Method
	Action        ActionFunc
	LoginRequired bool

	handler gin.HandlerFunc
}

// GetHandler for request handler
func (a *Action) GetHandler() gin.HandlerFunc {
	return a.handler
}

// NewAction creates new action
func NewAction(method Method, handler ActionFunc, loginRequired bool) *Action {
	a := Action{Method: method, Action: handler, LoginRequired: loginRequired}
	a.handler = func(c *gin.Context) {
		// init context
		ctx := Context{C: c}

		if a.LoginRequired {
			err := ctx.GetUserError()
			if err == nil {
				err = ErrLoginRequired
			}
			abortError(c, err)
			return
		}

		r, err := a.Action(&ctx)
		if err != nil {
			if err != ErrRawResponse {
				abortError(c, err)
			}
			return
		}
		c.JSON(http.StatusOK, basicResponse{Code: CodeSuccess, Info: r})
	}
	return &a
}

func abortError(c *gin.Context, err error) {
	switch v := err.(type) {
	case (*APIError):
		// TODO: specified code
		abortWithError(c, v.Status, CodeUnknownError, v.Message)
	case ResponseError:
		if l, ok := v.(ContextLogger); ok {
			l.Log(ErrorLevel, c.Request.URL.String())
		}
		c.Status(v.StatusCode())
		abortWithError(c, v.StatusCode(), v.ErrorCode(), v.ErrorInfo())
	default:
		c.Status(http.StatusInternalServerError)
		abortWithError(c, http.StatusInternalServerError, CodeUnknownError, err.Error())
	}
}

func abortWithError(c *gin.Context, status, code int, info interface{}) {
	c.AbortWithStatusJSON(status, basicResponse{Code: code, Info: info})
}
