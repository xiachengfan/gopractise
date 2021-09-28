package framework

import (
	"net/http"
)

type HandlerFunc func(*Context)

type Core struct {
	RGroup
	router          *router
	RemoteIPHeaders []string
}

var _ IGroup = &Core{}

func NewCore() *Core {
	engine := &Core{router: newRouter()}
	engine.RGroup = RGroup{
		Handlers: nil,
		parent:   "",
		root:     true,
		core:     engine,
	}
	return engine
}

func (c *Core) addRoute(method string, pattern string, handler HandlerFunc) {
	c.router.addRoute(method, pattern, handler)
}

func (c *Core) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	context := newContext(response, request)
	c.router.handle(context)
}

// Run defines the method to start a http server
func (c *Core) Run(addr string) (err error) {
	return http.ListenAndServe(addr, c)
}
