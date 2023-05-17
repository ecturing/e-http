package ehttp
import (
	"fmt"
	"net/http"
	"net/url"
)

//共享接口
type ServerHTTP func(r *Request, rp *Response)

var (
	ReadChannel  = make(chan []byte, 10)
	WriteChannel = make(chan []byte, 10)
)

type Response struct {
}

type Request struct {
	Method   string      //请求方法
	URL      *url.URL    //请求地址
	Proto    string      //协议类型
	Header   http.Header //请求头
	Body     string      //请求体
	Response Response    //响应体
}

type ResponseHandler interface {
	ReadResquest()
	GetHeader()
	GetBody()
}

type RequestHandler interface {
}

//Socket流读取
func ReadResquest(b []byte) *Request {
	pattern:=string(b)
	Root.Search(pattern)
}

func BufferRead() {
	fmt.Println("infor:buffer Read")
	for {
		select {
		case r := <-ReadChannel:
			ReadResquest(r)
		}
	}
}
