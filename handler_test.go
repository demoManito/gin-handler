package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModeFuncs(t *testing.T) {
	orgMode := Mode()
	defer SetMode(orgMode)
	assert := assert.New(t)
	SetMode(ReleaseMode)
	assert.Equal(ReleaseMode, Mode())
	SetMode(DebugMode)
	SetMode("")
	assert.Equal(ReleaseMode, Mode())
}

func TestHandlerMount(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	middleware := func(c *gin.Context) {
		c.Set("url", c.Request.URL.String())
		c.Next()
	}
	actionTest := func(c *Context) (ActionResponse, error) {
		url, _ := c.C.Get("url")
		return url, nil
	}

	r := gin.Default()
	h3 := Handler{
		Name:    "3",
		Actions: map[string]*Action{"/test": NewAction(GET, actionTest, false)},
	}
	h2 := Handler{
		Name:        "2",
		Actions:     map[string]*Action{"/test": NewAction(GET, actionTest, false)},
		SubHandlers: []Handler{h3},
	}
	h1 := Handler{
		Name:        "1",
		Middlewares: gin.HandlersChain{middleware},
		Actions:     map[string]*Action{"/test": NewAction(POST, actionTest, false)},
		SubHandlers: []Handler{h2},
	}
	h1.Mount(r)
	s := httptest.NewServer(r)
	defer s.Close()

	resp, err := http.Post(fmt.Sprintf("%s/1/test", s.URL), "", nil)
	require.NoError(err)
	body, _ := ioutil.ReadAll(resp.Body)
	assert.JSONEq(`{"code":0,"info":"/1/test"}`, string(body))

	resp, err = http.Get(fmt.Sprintf("%s/1/2/test", s.URL))
	require.NoError(err)
	body, _ = ioutil.ReadAll(resp.Body)
	assert.JSONEq(`{"code":0,"info":"/1/2/test"}`, string(body))

	resp, err = http.Get(fmt.Sprintf("%s/1/2/3/test", s.URL))
	require.NoError(err)
	body, _ = ioutil.ReadAll(resp.Body)
	assert.JSONEq(`{"code":0,"info":"/1/2/3/test"}`, string(body))
}
