package main

import (
	"net/http"
)

type Server interface {
	Routable
	Start(address string) error
}

type sdkHttpServer struct{
	Name string
	handler Handler
	root Filter
}

// Route 核心 只依赖一些接口
// 路由支持 POST GET PUT DELETE 等请求方法检查
// Route 委托给Handler处理
func (s *sdkHttpServer) Route(method string, pattern string, handlerFunc handlerFunc) {
	// 注册路由
	s.handler.Route(method, pattern, handlerFunc)
/*
	key := s.handler.key(method, pattern)
	s.handler.Handlers[key] = handlerFunc*/
}

func (s *sdkHttpServer) Start(address string) error {
	http.HandleFunc("/", func(writer http.ResponseWriter,
		request *http.Request) {
		c := NewContext(writer, request)
		s.root(c)
	})
	return http.ListenAndServe(address, nil)
}

func NewHttpServer(name string, builders... FilterBuilder) Server{

    handler := NewHandlerBasedOnMap()
    // 最后的业务逻辑也需要放在这里
	var root Filter = handler.ServeHTTP
	// 反向链接
	for i := len(builders)-1; i >=0 ; i-- {
		b := builders[i]
		root = b(root)
	}

	return &sdkHttpServer{
		Name: name,
		handler: handler,
		root: root,
	}
}


