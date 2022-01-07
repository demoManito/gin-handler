package handler

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// pageSize limit value and default value
const (
	DefaultPageSize int64 = 20
	MaxPageSize     int64 = 100
)

// Context request context
type Context struct {
	C *gin.Context

	equipOnce sync.Once
	equipErr  error

	buvid     string
	buvidOnce sync.Once
	buvidErr  error

	userOnce sync.Once
	userErr  error
}

// UserAgent returns the "User-Agent" header
func (ctx *Context) UserAgent() string {
	return ctx.C.GetHeader("User-Agent")
}

// ClientIP returns the client IP
func (ctx *Context) ClientIP() string {
	return ctx.C.ClientIP()
}

// Request return the request
func (ctx *Context) Request() *http.Request {
	return ctx.C.Request
}

// GetEquipError returns the error associated with the util.Equipment in Context
func (ctx *Context) GetEquipError() error {
	return ctx.equipErr
}

// GetUserError returns the error associated with the user.User in Context
func (ctx *Context) GetUserError() error {
	return ctx.userErr
}

// BindJSON decodes the body into v as json.Unmarshal does
func (ctx *Context) BindJSON(v interface{}) error {
	return ctx.C.BindJSON(v)
}

// Bind checks the Content-Type to select a binding engine automatically
func (ctx *Context) Bind(v interface{}) error {
	return ctx.C.Bind(v)
}

// GetDefaultParam gets params from query if successful, otherwise post form, otherwise defaultValue
func (ctx *Context) GetDefaultParam(key string, defaultValue string) string {
	value, ok := ctx.C.GetQuery(key)
	if ok {
		return value
	}
	return ctx.C.DefaultPostForm(key, defaultValue)
}

// GetParam gets params from query if successful, otherwise post form
func (ctx *Context) GetParam(key string) (string, bool) {
	value, ok := ctx.C.GetQuery(key)
	if ok {
		return value, true
	}
	return ctx.C.GetPostForm(key)
}

// GetDefaultParamString returns string without leading and trailing white space if successful,
// otherwise post form, otherwise defaultValue
func (ctx *Context) GetDefaultParamString(key string, defaultValue string) string {
	value, ok := ctx.GetParamString(key)
	if ok {
		return value
	}
	return defaultValue
}

// GetParamString returns string without leading and trailing white space
func (ctx *Context) GetParamString(key string) (string, bool) {
	raw, ok := ctx.GetParam(key)
	if !ok {
		return "", false
	}
	return strings.TrimSpace(raw), true
}

// GetDefaultParamInt gets params from query if successful, otherwise post form, otherwise defaultValue
func (ctx *Context) GetDefaultParamInt(key string, defaultValue int) (int, error) {
	val, ok := ctx.GetParam(key)
	if !ok {
		return defaultValue, nil
	}
	value, err := strconv.Atoi(val)
	return value, err
}

// GetParamInt gets params from query if successful, otherwise post form
func (ctx *Context) GetParamInt(key string) (int, error) {
	val, ok := ctx.GetParam(key)
	if !ok {
		return 0, ErrEmptyValue
	}
	value, err := strconv.Atoi(val)
	return value, err
}

// GetDefaultParamInt64 gets params from query if successful, otherwise post form, otherwise defaultValue
func (ctx *Context) GetDefaultParamInt64(key string, defaultValue int64) (int64, error) {
	val, ok := ctx.GetParam(key)
	if !ok {
		return defaultValue, nil
	}
	value, err := strconv.ParseInt(val, 10, 64)
	return value, err
}

// GetParamInt64 gets params from query if successful, otherwise post form
func (ctx *Context) GetParamInt64(key string) (int64, error) {
	val, ok := ctx.GetParam(key)
	if !ok {
		return 0, ErrEmptyValue
	}
	value, err := strconv.ParseInt(val, 10, 64)
	return value, err
}

// GetParamDateRange gets the start_date and end_date parameters
func (ctx *Context) GetParamDateRange(defaultStartDate, defaultEndDate string) (startDate, endDate time.Time, err error) {
	startDateStr, ok := ctx.GetParam("start_date")
	if !ok && defaultStartDate != "" {
		startDateStr = defaultStartDate
	}
	endDateStr, ok := ctx.GetParam("end_date")
	if !ok && defaultEndDate != "" {
		endDateStr = defaultEndDate
	}

	startDate, err = time.ParseInLocation("2006-01-02", startDateStr, time.Local)
	if err != nil {
		return
	}
	endDate, err = time.ParseInLocation("2006-01-02", endDateStr, time.Local)
	if err != nil {
		return
	}
	if endDate.Before(startDate) {
		err = ErrInvalidDateRange
		return
	}
	return
}

// GetParamMonthRange gets the start_month and end_month parameters
func (ctx *Context) GetParamMonthRange(defaultStartDate, defaultEndDate string) (startTime, endTime time.Time, err error) {
	startTimeStr, ok := ctx.GetParam("start_month")
	if !ok && defaultStartDate != "" {
		startTimeStr = defaultStartDate
	}
	endTimeStr, ok := ctx.GetParam("end_month")
	if !ok && defaultEndDate != "" {
		endTimeStr = defaultEndDate
	}

	startTime, err = time.ParseInLocation("2006-01", startTimeStr, time.Local)
	if err != nil {
		return
	}
	endTime, err = time.ParseInLocation("2006-01", endTimeStr, time.Local)
	if err != nil {
		return
	}
	if endTime.Before(startTime) {
		err = ErrInvalidDateRange
		return
	}
	return
}

// PageOption param page option
type PageOption struct {
	DefaultPageSize int64 // pageSize, 默认 20
	MaxPageSize     int64 // 被允许使用的最大 pageSize, 默认 100
}

// GetParamPage gets the p and page_size parameters
// NOTICE: return: p 默认值: 1, pageSize 默认值: 20
func (ctx *Context) GetParamPage(opt ...*PageOption) (p, pageSize int64, err error) {
	p, err = ctx.GetDefaultParamInt64("p", 1)
	if err != nil || p <= 0 {
		return 0, 0, ErrInvalidParam
	}
	option := getPageOption(opt...)
	pageSize, err = getPageSize(ctx, option)
	if err != nil || pageSize <= 0 || pageSize > option.MaxPageSize {
		return 0, 0, ErrInvalidParam
	}
	return
}

func getPageSize(ctx *Context, option *PageOption) (pageSize int64, err error) {
	size, exists := ctx.GetParam("pagesize")
	if !exists {
		size, exists = ctx.GetParam("page_size")
		if !exists {
			return option.DefaultPageSize, nil
		}
	}
	return strconv.ParseInt(size, 10, 64)
}

func getPageOption(opt ...*PageOption) *PageOption {
	length := len(opt)
	if length == 0 {
		return &PageOption{DefaultPageSize: DefaultPageSize, MaxPageSize: MaxPageSize}
	}
	if length > 1 {
		panic("invalid page option")
	}
	if opt[0].DefaultPageSize == 0 {
		opt[0].DefaultPageSize = DefaultPageSize
	}
	if opt[0].MaxPageSize == 0 {
		opt[0].MaxPageSize = MaxPageSize
	}
	if opt[0].DefaultPageSize < 0 || opt[0].MaxPageSize < opt[0].DefaultPageSize {
		panic("invalid page option")
	}
	return opt[0]
}

// Token 从 cookie 中获取 token, 没有返回空字符串
func (ctx *Context) Token() string {
	token, _ := ctx.C.Cookie("token")
	return token
}
