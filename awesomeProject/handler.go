package main

import "net/http"

type Routable interface {
	Route(method string, pattern string, handlerFunc handlerFunc)
}

type Handler interface {
	ServeHTTP(c *Context)
	Routable
}

// 确保实现handler接口
var _ Handler = &HandlerBasedOnMap{}

type HandlerBasedOnMap struct {
	Handlers map[string]handlerFunc
	root Filter
}

func (h *HandlerBasedOnMap) Route(method string, pattern string, handlerFunc handlerFunc) {
	key := h.key(method, pattern)
	h.Handlers[key] = handlerFunc
}

func (h *HandlerBasedOnMap) ServeHTTP(c *Context) {
	// 分发路由
	request := c.R
	writer := c.W
	key := h.key(request.Method, request.URL.Path)

	if handler, ok := h.Handlers[key]; ok {

		// ctx := NewContext(writer, request)
		handler(c)
	}else {
		writer.WriteHeader(http.StatusNotFound)
		writer.Write([]byte("Not Found!"))
	}
}

func (h *HandlerBasedOnMap) key(method string, pattern string) string{
	return method + "#" + pattern
}

func NewHandlerBasedOnMap() Handler {

	return &HandlerBasedOnMap{
		Handlers: make(map[string]handlerFunc),
	}
}
