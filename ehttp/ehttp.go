package ehttp

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	NEWLINE = "\r\n"
	SPACE   = " "
)

// 共享接口
type ServerHTTP func(r *Request, rp *Response)

var (
	ReadQueen   = make(chan *Event, 10)
	body_buffer = make([]byte, 10)
	WriteQueen  = make(chan *Event, 10)
)

// --------------响应体--------------------
type Response struct {
	protocal string
	Status   int
	OK       string
	Headers  map[string]string
	DataFrom string
}

// 响应处理
type ResponseHandler interface {
	// 响应体写入
	WriteResponse(msg any)
	//响应体序列化
	ResponseSerializer()
}

func (w *Response) ResponseSerializer() *bytes.Buffer {
	//对Response整个结构体序列化为[]bytes
	w.protocal = "HTTP/1.0"
	w.Status = 200
	w.OK = "OK"
	w.Headers = make(map[string]string)
	w.Headers["Content-Type"] = "text/plain"
	var b strings.Builder
	b.WriteString(w.protocal + 
		SPACE+fmt.Sprintf("%d", w.Status) + 
		SPACE+w.OK+
		NEWLINE+
		"Content-Type: " + 
		w.Headers["Content-Type"] + 
		NEWLINE+NEWLINE+NEWLINE+w.DataFrom)
	buf := bytes.NewBuffer([]byte(b.String()))
	return buf
}

// 请求体
type Request struct {
	Method RequestMethod //请求方法
	URL    *url.URL      //请求地址
	Proto  string        //协议类型
	Header http.Header   //请求头
	Body   string        //请求体
}

// 请求处理
type RequestHandler interface {
	ReadRequest()
	GetHeader()
	GetBody()
}

// 请求体构造
func NewRequest(method string, url *url.URL, pro string, header http.Header, body string) *Request {
	return &Request{
		Method: RTR(method),
		URL:    url,
		Proto:  pro,
		Header: header,
		Body:   body,
	}
}

// Socket流读取Request
func ReadRequest(router *Router, readevent *bufio.Reader) {
	req, _ := http.ReadRequest(readevent)
	var rep = &Response{}
	ereq := NewRequest(req.Method, req.URL, req.Proto, req.Header, string(body_buffer))
	hander, err := router.Search(ereq.URL.Path)
	if err != nil {
		fmt.Println("路由查找失败", err)
	}
	hander(ereq, rep)
	var event = &Event{
		Conn:   nil,
		Reader: nil,
		Writer: rep.ResponseSerializer(),
	}
	WriteQueen <- event
}

// 传统request method转自有request method
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
