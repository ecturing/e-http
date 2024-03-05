package ehttp

import (
	"bufio"
	"bytes"
	"errors"
	"ews/Eerror"
	"ews/log"
	"io"
	"strconv"
	"strings"
)

const (
	NEWLINE = "\r\n"
	SPACE   = " "
)

// 共享接口
type ServerHTTP func(r *E_Request, rp *E_Response)

var (
	ReadQueen  = make(chan *Event, 10)
	CloseQueen = make(chan struct{}, 10)
	WriteQueen = make(chan *Event, 10)
)

// --------------响应体--------------------
type E_Response struct {
	protocal string
	Status   int
	OK       string
	Headers  map[string]string
	DataFrom string
}

// 响应处理
type ResponseHandler interface {
	//响应体序列化
	ResponseSerializer() *bytes.Buffer
	// 默认响应头
	DefaultHeader()
}

// 响应体序列化
func (w *E_Response) ResponseSerializer() *bytes.Buffer {
	//对Response整个结构体序列化为[]bytes
	var b, header strings.Builder
	for k, v := range w.Headers {
		header.WriteString(k + ": " + v + NEWLINE)
	}

	// 拼接响应体,header已经加上一个换行符了，所以这里只需要加一个换行符
	b.WriteString(w.protocal + SPACE + strconv.Itoa(w.Status) + SPACE + w.OK +
		NEWLINE +
		header.String() +
		NEWLINE +
		w.DataFrom)
	buf := bytes.NewBuffer([]byte(b.String()))
	return buf
}

func (w *E_Response) DefaultHeader() {
	w.Headers["Content-Type"] = "text/plain"
	w.Headers["Server"] = "EWS"

	buf := bytes.NewBufferString(w.DataFrom)
	w.Headers["Content-Length"] = strconv.Itoa(buf.Len())
}

// 请求体
type E_Request struct {
	Method RequestMethod     //请求方法
	URL    *E_URL            //请求地址
	Proto  string            //协议类型
	Header map[string]string //请求头
	Body   io.ReadCloser     //请求体
}

// URL
type E_URL struct {
	Path  string
	Query string
}

// 请求体构造
func NewRequest(b *bufio.Reader) (req *E_Request) {
	line, header, breakerr := GetLineHeader(b)
	if breakerr != nil {
		log.Logger.Error().Err(breakerr).Msg("请求行与请求头解析错误")
		return
	}
	Header := ReadHeader(header)
	log.Logger.Info().Msg("request header:" + header)
	log.Logger.Info().Msg("request line:" + line)
	method, URL, proto, e := ReadRequestLine(line)
	if e != nil {
		log.Logger.Error().Err(e).Msg("request parse failed")
		return nil
	}
	routerPath, values, err := ReadPathValues(URL)
	if err != nil {
		log.Logger.Error().Err(err).Msg("URL parse failed")
	}
	var body io.ReadCloser
	if values == "" {
		log.Logger.Info().Msg("don`t have query args")
		value, ok := Header["Content-Length"]
		if ok {
			content_length, err := strconv.Atoi(value)
			if err != nil {
				log.Logger.Error().Err(err).Msg("Content-Length args invalid")
			}
			body = ReadBody(b, content_length)
		} else {
			ErrSingal(Eerror.BadRequest)
		}
		log.Logger.Info().Msg("请求体:查询")
	}
	return &E_Request{
		Method: RTR(method),
		URL: &E_URL{
			Path:  routerPath,
			Query: values,
		},
		Proto:  proto,
		Header: Header,
		Body:   body,
	}
}

// 分离请求行与请求头
func GetLineHeader(b *bufio.Reader) (string, string, error) {
	rline, e := b.ReadString('\n')
	if e != nil {
		return "", "", e
	}
	rline = strings.TrimSuffix(rline, "\r\n")
	var Headerbuild strings.Builder
	for {
		line, _, err := b.ReadLine()
		if err != nil && err != io.EOF {
			log.Logger.Error().Err(err).Msg("read header error")
			return "", "", err
		}

		// 检查是否读取到了两个CRLF（\r\n），表示HTTP头部已经读取完毕，接下来的内容是消息主体
		if len(line) == 0 {
			break
		} else {
			Headerbuild.WriteString(string(line) + "\r\n")
		}

	}
	return rline, Headerbuild.String(), nil
}

// 读取request流中两次换行符之间的内容，返回一个*bufio.Reader
func ReadBody(reader *bufio.Reader, length int) (body io.ReadCloser) {
	var body_buffer = make([]byte, length)
	_, err := reader.Read(body_buffer)
	if err != nil {
		log.Logger.Error().Err(err).Msg("read body error")
	}
	return io.NopCloser(bytes.NewReader(body_buffer))
}

// 处理请求行(获取请求的第一行中的请求方法、请求地址、协议类型)
// Method Path Proto
func ReadRequestLine(s string) (method, URL, proto string, e error) {

	defer func() {
		if err := recover(); err != nil {
			log.Logger.Error().Msg("ReadRequestLine error")
		}
	}()

	var emethod, url, protocol string
	var err error
	me, rest, ok := strings.Cut(s, " ")
	if ok {
		emethod = me
		u, p, ok2 := strings.Cut(rest, " ")
		if ok2 {
			url = u
			protocol = p
		} else {
			err = errors.New("protocol type lost")
			log.Logger.Error().Msg("protocol type lost")
		}
	} else {
		err = errors.New("请求方法或丢失")
		log.Logger.Error().Msg("请求方法或丢失")
	}
	return emethod, url, protocol, err
}

// 从标准url获取url的路由与查询参数(注意类似POST这种是不在url中写参数，要做好校验)
// 标准url格式 scheme://host:port/path?query
func ReadPathValues(url string) (routerPath, values string, err error) {
	var path, value string
	index := strings.Index(url, "?")
	if index == -1 {
		path = url
		value = ""
	} else {
		path = url[:index]
		value = url[index+1:]
	}
	return path, value, nil
}

// 读取请求头
func ReadHeader(b string) map[string]string {
	reader := bufio.NewReader(strings.NewReader(b))
	headers := make(map[string]string)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		parts := strings.SplitN(string(line), ": ", 2)
		if len(parts) == 2 {
			headers[parts[0]] = parts[1]
		}
	}
	return headers
}

// Socket流读取Request
func ReadRequest(router *Router, reader *bufio.Reader) ConnClose {

	defer func() {
		if err := recover(); err != nil {
			ErrSingal(err.(error))
			log.Logger.Error().Err(err.(error)).Msg("ReadRequest error")
		}
	}()

	req := NewRequest(reader)
	if req == nil {
		return ConnClose(true)
	}
	var rep = &E_Response{
		protocal: "HTTP/1.1",
		Status:   200,
		OK:       "OK",
		Headers:  make(map[string]string),
		DataFrom: "",
	}
	switch req.Method {
	case GET:
		// do
		handler, err := router.Search(req.URL.Path, req.Method)
		if err != nil {
			ErrSingal(err)
			return ConnClose(true)
		}
		log.Logger.Info().Msg("pre handle get request")
		handler(req, rep)
		rep.DefaultHeader()
	case POST:
		// do
		handler, err := router.Search(req.URL.Path, req.Method)
		if err != nil {
			ErrSingal(err)
			return ConnClose(true)
		}
		handler(req, rep)
		rep.DefaultHeader()
	case PUT:
		// do

	case DELETE:
		// do

	default:
		panic(Eerror.UnSupportMethod)
	}
	proto := req.Proto
	switch proto {
	case "HTTP/1.0":
		// do
		rep.Send(true)
		return ConnClose(true)
	case "HTTP/1.1":
		// do
		rep.Send(false)
		return ConnClose(false)
	default:
		log.Logger.Debug().Msg("不支持的协议类型")
		ErrSingal(Eerror.UnSupportProto)
		return ConnClose(true)
	}
}

// 传统request method转自有request method
func RTR(method string) RequestMethod {
	switch method {
	case "GET":
		return 1
	case "POST":
		return 2
	case "PUT":
		return 3
	case "DELETE":
		return 4
	}
	return -1
}

func send(rep *E_Response) {
	var event = &Event{
		Reader: nil,
		Writer: rep.ResponseSerializer(),
	}
	WriteQueen <- event
}

func (rep *E_Response) Send(close bool) {
	if close {
		rep.Headers["Connection"] = "close"
	}
	send(rep)
}
