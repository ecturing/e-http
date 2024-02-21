package ehttp

import (
	"ews/log"
)

// 请求与函数组合+套接字启动
func Server(r *Router, pattern string, f ServerHTTP, m RequestMethod) {
	r.Register(pattern, f, m)
	log.Logger.Info().Msg("server start")
	go r.RouterListen()
}

func Confirm(s string) {
	err := InitSocket(s)
	if err != nil {
		log.Logger.Fatal().Err(err).Msgf("socket error %v", err)
	}
}
