package Eerror

// 错误码枚举，无符号短整型
var (
	// 生成Http错误码枚举
	NotFound        = NetError{Code: 404, Msg: "Not Found"}
	MethodNotAllow  = NetError{Code: 405, Msg: "Method Not Allow"}
	OK              = NetError{Code: 200, Msg: "OK"}
	UnSupportMethod = NetError{Code: 501, Msg: "UNSUPPORT METHOD"}
	BadRequest      = NetError{Code: 400, Msg: "Bad Request"}

	SERVERERR = ServerError{Code: 500, Msg: "Server Error"}
)

type NetError struct {
	// 错误码
	Code int
	// 错误信息
	Msg string
}

type ServerError struct {
	// 错误码
	Code int
	// 错误信息
	Msg string
}

func (e NetError) Error() string {
	return e.Msg
}
func (e *NetError) GetCode() int {
	return e.Code
}

func (e ServerError) Error() string {
	return e.Msg
}
func (e *ServerError) GetCode() int {
	return e.Code
}
