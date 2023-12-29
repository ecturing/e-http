package ehttp

import "fmt"

// 请求与函数组合+套接字启动
func Server(r *Router, pattern string, f ServerHTTP,m RequestMethod) {
	r.Register(pattern, f,m)
	go r.RouterListener() //启动路由系统
}

func Confirm(s string) {
	err := InitSocket(s)
	if err != nil {
		fmt.Println("socket初始化失败", err)
	}
}
