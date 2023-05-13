package ehttp

type ServerHTTP func(r Request,rp *Response)

//请求与函数组合
func Server(pattern string,f ServerHTTP){
	root:=&treeNode{
		pattern: "/",
		EndNode: false,
		childNode: nil,
		hander: nil,
	}
	root.Register(pattern,f)
}