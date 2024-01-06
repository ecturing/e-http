package ehttp

import (
	"github.com/rs/zerolog/log"
)

// 请求与函数组合+套接字启动
func Server(r *Router, pattern string, f ServerHTTP,m RequestMethod) {
	r.Register(pattern, f,m)
	log.Info().Msg("server start")
	go r.RouterListen() 
}

func Confirm(s string) {
	err := InitSocket(s)
	if err != nil {
		log.Fatal().Err(err).Msgf("socket error %v", err)
	}
}
