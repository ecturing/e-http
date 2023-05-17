package ehttp

import "fmt"

// 请求与函数组合+套接字启动
func Server(pattern string, f ServerHTTP) {
	Root.Register(pattern, f)
}

func Confirm(s string) {
	go func() {
		err := InitSocket(s)
		if err != nil {
			fmt.Println("error:%v", err)
		}
	}()
}
