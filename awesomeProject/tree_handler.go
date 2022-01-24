package main

import (
	"net/http"
	"strings"
)

type node struct {
	path string
	children []*node

	//该节点注册的方法
	handler handlerFunc

}

func (n *node)findMatchChild(path string) (*node, bool) {
	var wildcardNode *node
	for _, child := range n.children {
		if child.path == path {
			return child, true
		}
		if child.path == "*" {
			wildcardNode = child
		}
	}
	return wildcardNode, wildcardNode != nil
}

func (n *node) createSubTree(paths []string, handlerFunc handlerFunc) {
	cur := n
	for _, path := range paths {
		nn := newNode(path)
		cur.children = append(cur.children, nn)
		cur = nn
	}
	cur.handler = handlerFunc
}

func newNode(path string) *node {
	return &node{
		path: path,
		children: make([]*node,0,4),
	}
}


func (h *HandlerBasedOnTree) ServeHTTP(c *Context) {
	//查找路由
	path := c.R.URL.Path
	handlerFunc, ok := h.findRouter(path)
	if !ok {
		c.W.WriteHeader(http.StatusNotFound)
		c.W.Write([]byte("Not Found"))
	}
	handlerFunc(c)
}

func (h *HandlerBasedOnTree) findRouter(path string) (handlerFunc, bool) {
	paths := strings.Split(strings.Trim(path, "/"),"/")
	cur := h.root
	for _, p := range paths {
		matchChild, ok := cur.findMatchChild(p)
		if !ok {
			return nil, false
		}
		cur = matchChild
	}
	if cur.handler == nil {
		//注册/usr/profile
		//访问/usr
		return nil, false
	}
	return cur.handler, true
}

// 注册路由
func (h *HandlerBasedOnTree) Route(method string, pattern string, handlerFunc handlerFunc) {
	pattern = strings.Trim(pattern,"/")
	paths := strings.Split(pattern,"/")

	cur := h.root
	for index, path := range paths {
		matchChild, ok := cur.findMatchChild(path)
		if ok {
			cur = matchChild
		}else {
			cur.createSubTree(paths[index:], handlerFunc)
			return
		}
	}
}


type HandlerBasedOnTree struct {
	root *node
}
