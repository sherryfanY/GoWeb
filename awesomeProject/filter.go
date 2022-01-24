package main

import (
	"fmt"
	"time"
)

type handlerFunc func(c *Context)

type Filter func(c *Context)

type FilterBuilder func(next Filter) Filter

var _ FilterBuilder = MatricFilterBuilder

func MatricFilterBuilder(next Filter) Filter{
	return func(c *Context) {
		start := time.Now().Nanosecond()
		// 调用next参数
		next(c)
		end := time.Now().Nanosecond()
		fmt.Printf("执行用了 %d 纳秒", start - end)
	}
}