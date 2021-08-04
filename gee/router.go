package gee

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc), roots: make(map[string]*node)}
}

func parsePattern(pattern string) []string {
	splitted := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range splitted {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = NewNode()
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) getRoute(method string, pattern string) (*node, map[string]string) {
	parts := parsePattern(pattern)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	params := make(map[string]string)
	n := root.search(parts, 0)
	if n != nil {
		newParts := parsePattern(n.pattern)
		for index, newPart := range newParts {
			if newPart[0] == ':' {
				params[newPart[1:]] = parts[index]
			}

			if newPart[0] == '*' && len(newPart) > 1 {
				params[newPart[1:]] = strings.Join(parts[index:], "/")
			}
		}
		return n, params
	}

	return nil, nil
}

func (r *router) getRoutes(method string) []*node{
	root, ok := r.roots[method]
	if !ok {
		return nil
	}

	nodes := make([]*node, 0)
	root.travel(&nodes)
	return nodes
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(context *Context) {
			context.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}
