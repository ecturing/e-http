package ehttp


//请求与函数组合+套接字启动
func Server(pattern string, f ServerHTTP) {
	Root.Register(pattern, f)
}

func Confirm(s string) error{
	err:=Init_Socket(s)
	return err
}





