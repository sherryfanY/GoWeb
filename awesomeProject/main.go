package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func SignUp(ctx *Context) {
	// json 读取数据
    req := &signUpReq{}
	err := ctx.ReadJson(req)
	if err != nil {
		ctx.BadRequestJson(err)
		return
	}
	//返回json数据
	ctx.OKJson(req.Email)
}

func main() {
	server := NewHttpServer("test-server")
	server.Route(http.MethodGet,"/home", SignUp)
	err := server.Start("8080")
	if err != nil {
		// 快速失败
		panic(err)
	}
}
