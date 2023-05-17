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
	TCPBuffer  = make(chan *bufio.Reader, 10)
	connBuffer = make([]byte, 10)
)

type Response struct {
	DataFrom []byte
}

type Request struct {
	Method   string      //请求方法
	URL      *url.URL    //请求地址
	Proto    string      //协议类型
	Header   http.Header //请求头
	Body     string      //请求体
	Response *Response   //响应体
}

type ResponseHandler interface {
	ReadRequest()
	GetHeader()
	GetBody()
}

type RequestHandler interface {
}

// ReadRequest Socket流读取
func ReadRequest(b *bufio.Reader) {
	req, _ := http.ReadRequest(b)
	var rep = &Response{}
	req.Body.Read(connBuffer)
	ereq := NewRequest(req.Method, req.URL, req.Proto, req.Header, string(connBuffer), rep)
	search, err := Root.Search(ereq.URL.Path)
	if err != nil {
		fmt.Println(err)
	}
	search(ereq, rep)
}

func RouterListener() {
	fmt.Println("infor:buffer Read")
	for {
		select {
		case r := <-TCPBuffer:
			ReadRequest(r)
		}
	}
}

func (w *Response) ResponseWriter(any string) {
	w.DataFrom = []byte(any)
}

func NewRequest(method string, url *url.URL, pro string, header http.Header, body string, rep *Response) *Request {
	return &Request{
		Method:   method,
		URL:      url,
		Proto:    pro,
		Header:   header,
		Body:     body,
		Response: rep,
	}
}
