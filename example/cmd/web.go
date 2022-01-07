package cmd

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gin-handler"
	"gin-handler/example/controllers/apiv1"
	"gin-handler/example/controllers/apiv2"
)

func Run() {
	srv := newServer()
	srv.ListenAndServe()
}

func newServer() *http.Server {
	router := gin.New()
	router.Use(gin.Logger())
	mount(router)

	return &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
}

// 挂载路由
func mount(r *gin.Engine) {
	handlers := []handler.Handler{
		apiv1.Handler(),
		apiv2.Handler(),
	}
	for _, h := range handlers {
		// 每一个路由分别挂载
		h.Mount(r)
	}
}
