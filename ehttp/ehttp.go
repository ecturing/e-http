package ehttp

import (
	"bufio"
	"bytes"
	"errors"
	"ews/Eerror"
	"ews/logutil"
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
	ReadQueen = make(chan *Event, 10)
	// body_buffer = make([]byte, 1024)
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
                           
// 默认响应头
func (w *E_Response) DefaultHeader() {
	w.Headers["Content-Type"] = "text/plain"
	w.Headers["Server"] = "EWS"
}

// 设置关闭连接请求头
func (w *E_Response)SetConnClose()  {
	w.Headers["Connection"] = "close"
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

// 请求处理
type RequestHandler interface {
	//获取请求头
	GetHeader()
	//获取请求体
	GetBody()
}

// 请求体构造
func NewRequest(b *bufio.Reader) *E_Request {
	line, header := GetLineHeader(b)
	Header := ReadHeader(header)
	logutil.Logger.Info().Msg("请求头:" + header)
	logutil.Logger.Info().Msg("请求行:" + line)
	method, URL, proto, e := ReadRequestLine(line)
	if e != nil {
		logutil.Logger.Error().Err(e).Msg("请求行解析错误")
		return nil
	}
	routerPath, values, err := ReadPathValues(URL)
	if err != nil {
		logutil.Logger.Error().Err(err).Msg("URL解析错误")
	}
	var body io.ReadCloser
	if values == "" {
		logutil.Logger.Info().Msg("没有查询参数")
		value, ok := Header["Content-Length"]
		if ok {
			content_length, err := strconv.Atoi(value)
			if err != nil {
				logutil.Logger.Error().Err(err).Msg("Content-Length参数非法")
			}
			body = ReadBody(b, content_length)
		} else {
			ErrSingal(Eerror.BadRequest)
		}

		logutil.Logger.Info().Msg("请求体:查询")
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
func GetLineHeader(b *bufio.Reader) (string, string) {
	rline, e := b.ReadString('\n')
	rline = strings.TrimSuffix(rline, "\r\n")
	if e != nil {
		return "", ""
	}
	var Headerbuild strings.Builder
	for {
		line, _, err := b.ReadLine()
		if err != nil && err != io.EOF {
			logutil.Logger.Error().Err(err).Msg("读取HTTP请求体错误")
			return "", ""
		}

		// 检查是否读取到了两个CRLF（\r\n），表示HTTP头部已经读取完毕，接下来的内容是消息主体
		if len(line) == 0 {
			break
		} else {
			Headerbuild.WriteString(string(line) + "\r\n")
		}

	}
	return rline, Headerbuild.String()
}

// 读取request流中两次换行符之间的内容，返回一个*bufio.Reader
func ReadBody(reader *bufio.Reader, length int) io.ReadCloser {
	var body_buffer = make([]byte, length)
	_, err := reader.Read(body_buffer)
	if err != nil {
		logutil.Logger.Error().Err(err).Msg("读取HTTP请求体错误")
	}
	return io.NopCloser(bytes.NewReader(body_buffer))
}

// 处理请求行(获取请求的第一行中的请求方法、请求地址、协议类型)
// Method Path Proto
func ReadRequestLine(s string) (method, URL, proto string, e error) {

	defer func() {
		if err := recover(); err != nil {
			logutil.Logger.Error().Msg("行解析错误")
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
			err = errors.New("协议类型丢失")
			logutil.Logger.Error().Msg("协议类型丢失")
		}
	} else {
		err = errors.New("请求方法或丢失")
		logutil.Logger.Error().Msg("请求方法或丢失")
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
func ReadRequest(router *Router, readevent *bufio.Reader) {

	defer func() {
		if err := recover(); err != nil {
			ErrSingal(err.(error))
			logutil.Logger.Error().Err(err.(error)).Msg("ReadRequest error")
		}
	}()

	req := NewRequest(readevent)
	var rep = &E_Response{
		protocal: "HTTP/1.0",
		Status:   200,
		OK:       "OK",
		Headers:  make(map[string]string),
		DataFrom: "",
	}
	switch req.Method {
	case GET:
		// do
		hander, err := router.Search(req.URL.Path, req.Method)
		if err != nil {
			ErrSingal(err)
			return
		}
		logutil.Logger.Info().Msg("进入处理环节")
		rep.DefaultHeader()
		hander(req, rep)
	case POST:
		// do
		hander, err := router.Search(req.URL.Path, req.Method)
		if err != nil {
			ErrSingal(err)
			return
		}
		rep.DefaultHeader()
		hander(req, rep)
	case PUT:
		// do

	case DELETE:
		// do

	default:
		panic(Eerror.UnSupportMethod)
	}
	var event = &Event{
		// Conn:   nil,
		Reader: nil,
		Writer: rep.ResponseSerializer(),
	}
	WriteQueen <- event
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
