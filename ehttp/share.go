package ehttp

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
)

// 共享接口
type ServerHTTP func(r *Request, rp *Response)

var (
	ReadQueen   = make(chan *Event, 10)
	body_buffer = make([]byte, 10)
	WriteQueen  = make(chan *Event, 10)
)

// 响应体
type Response struct {
	DataFrom string
}
//请求体
type Request struct {
	Method RequestMethod //请求方法
	URL    *url.URL      //请求地址
	Proto  string        //协议类型
	Header http.Header   //请求头
	Body   string        //请求体
}

// 响应处理
type ResponseHandler interface {
	WriteResponse()
}

//请求处理
type RequestHandler interface {
	ReadRequest()
	GetHeader()
	GetBody()
}

// Socket流读取Request
func ReadRequest(readevent *bufio.Reader) {
	req, _ := http.ReadRequest(readevent)
	var rep = &Response{}
	ereq := NewRequest(req.Method, req.URL, req.Proto, req.Header, string(body_buffer))
	hander, err := Root.Search(ereq.URL.Path, RTR(req.Method))
	if err != nil {
		fmt.Println("路由查找失败", err)
	}
	hander(ereq, rep)
	var event = &Event{
		Conn: nil,
		Reader: nil,
		Writer: []byte(rep.DataFrom),
	}
	WriteQueen <- event
}

func RouterListener() {
	fmt.Println("infor:buffer Read In")
	for r := range ReadQueen {
		ReadRequest(r.Reader)
	}
}
// 响应体写入
func (w *Response) ResponseWriter(msg any) {
	w.DataFrom = msg.(string)
}

func NewRequest(method string, url *url.URL, pro string, header http.Header, body string) *Request {
	return &Request{
		Method: RTR(method),
		URL:    url,
		Proto:  pro,
		Header: header,
		Body:   body,
	}
}

func RTR(method string) RequestMethod {
	switch method {
	case "GET":
		return 0
	case "POST":
		return 1
	case "PUT":
		return 2
	case "DELETE":
		return 3
	}
	return -1
}
