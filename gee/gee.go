package gee

import (
	"log"
	"net/http"
	"strings"
)

type HandlerFunc func(*Context)

type RouterGroup struct {
	prefix string
	middlewares []HandlerFunc
	engine *Engine
}

type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup
}

func (engine *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(request.RequestURI, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	context := newContext(writer, request)
	context.handlers = middlewares
	engine.router.handle(context)
}

func New() *Engine {
	engine := Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: &engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return &engine
}

func (group *RouterGroup) Use(handlers ...HandlerFunc) {
	group.middlewares = append(group.middlewares, handlers...)
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{prefix: group.prefix + prefix, engine: engine}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (group *RouterGroup) Run(port string) (err error) {
	return http.ListenAndServe(port, group.engine)
}
