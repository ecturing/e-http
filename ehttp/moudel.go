package ehttp

import (
	"ews/socket"
)

type Response struct{
   
}

type Request struct{
	
}



type ResponseHandler interface{
	ReadResquest();
	GetHeader();
	GetBody();
}


//端口配置
func Confirm(address string){
	socket.Init_Socket(address)
}


//Socket流读取
func (r *Request) ReadResquest()  {
	
}