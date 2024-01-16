package ehttp

import (
	"ews/Eerror"
	"ews/logutil"
)

var (
	ErrorSingal = make(chan *Event, 4)
)

// 错误处理,返回值为HTTP状态码
func errHandler(e error) int {
	switch err := e.(type) {
	case Eerror.NetError:
		logutil.Logger.Error().Err(err).Msg(err.Msg)
		return err.Code
	case Eerror.ServerError:
		logutil.Logger.Error().Err(err).Msg(err.Msg)
		return err.Code
	case nil:
		return Eerror.OK.Code
	default:
		return Eerror.SERVERERR.Code
	}

}

func ErrSingal(err error) {
	er := &E_Response{
		protocal: "HTTP/1.0",
		Status:   errHandler(err),
		OK:       "OK",
		Headers:  make(map[string]string),
		DataFrom: err.Error(),
	}
	er.DefaultHeader()
	er.ResponseSerializer()
	e := &Event{
		Reader: nil,
		Writer: er.ResponseSerializer(),
	}
	ErrorSingal <- e
}
