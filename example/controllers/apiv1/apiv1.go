package apiv1

import (
	"github.com/gin-gonic/gin"

	"gin-handler"
)

func Handler() handler.Handler {
	return handler.Handler{
		Name:        "v1",
		Middlewares: gin.HandlersChain{
			// 中间件
		},
		Actions: map[string]*handler.Action{
			"api1": handler.NewAction(handler.POST, API1, false), // http:localhost:8080/v1/api1
			"api2": handler.NewAction(handler.POST, API2, false), // http:localhost:8080/v1/api1
			"api3": handler.NewAction(handler.POST, API3, false), // http:localhost:8080/v1/api1
			"api4": handler.NewAction(handler.POST, API4, false), // http:localhost:8080/v1/api1
		},
		SubHandlers: []handler.Handler{oneHandler(), twoHandler(), threeHandler()},
	}
}

func oneHandler() handler.Handler {
	return handler.Handler{
		Name: "one",
		Actions: map[string]*handler.Action{
			"api1": handler.NewAction(handler.POST, API1, false), // http:localhost:8080/v1/one/api1
		},
	}
}

func twoHandler() handler.Handler {
	return handler.Handler{
		Name: "two",
		Actions: map[string]*handler.Action{
			"api2": handler.NewAction(handler.POST, API2, false), // http:localhost:8080/v1/two/api2
		},
	}
}

func threeHandler() handler.Handler {
	return handler.Handler{
		Name: "three",
		Actions: map[string]*handler.Action{
			"api3": handler.NewAction(handler.POST, API2, false), // http:localhost:8080/v1/three/api3
		},
	}
}

func API1(c *handler.Context) (handler.ActionResponse, error) {

	return nil, nil
}

func API2(c *handler.Context) (handler.ActionResponse, error) {

	return nil, nil
}

func API3(c *handler.Context) (handler.ActionResponse, error) {

	return nil, nil
}

func API4(c *handler.Context) (handler.ActionResponse, error) {

	return nil, nil
}
