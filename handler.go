package handler

import (
	"github.com/gin-gonic/gin"
)

// Method HTTP method
type Method int

// methods
const (
	GET  Method = 1 << iota
	POST Method = 1 << iota
)

// Modes
const (
	DebugMode   = "debug"
	ReleaseMode = "release"
	TestMode    = "test"
)

// Content-Type
const (
	ContentTypeJSON           = "application/json"
	ContentTypeFormUrlencoded = "application/x-www-form-urlencoded"
)

var (
	handlerMode = DebugMode
)

// SetMode set mode
func SetMode(mode string) {
	if mode == "" {
		mode = ReleaseMode
	}
	handlerMode = mode
	gin.SetMode(mode)
}

// Mode return mode
func Mode() string {
	return handlerMode
}

// M for map
type M map[string]interface{}

// S for slice
type S []interface{}

// ActionFunc handle the requests
type ActionFunc func(c *Context) (ActionResponse, error)

// ActionResponse for action response
type ActionResponse interface{}

// Handler contains all routers' info
type Handler struct {
	Name        string
	Middlewares gin.HandlersChain
	Actions     map[string]*Action
	SubHandlers []Handler
}

// Mount mount handler
func (h *Handler) Mount(r *gin.Engine) {
	h.mount(&r.RouterGroup)
}

func (h *Handler) mount(g *gin.RouterGroup) {
	g = g.Group(h.Name)
	g.Use(h.Middlewares...)
	for name, action := range h.Actions {
		actionHandler := action.GetHandler()
		if action.Method&GET != 0 {
			g.GET(name, actionHandler)
		}
		if action.Method&POST != 0 {
			g.POST(name, actionHandler)
		}
	}
	for _, sub := range h.SubHandlers {
		sub.mount(g)
	}
}
