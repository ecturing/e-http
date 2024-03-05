package ehttp

import (
	"ews/log"
)

type ServerArg struct {
	addr   string
	Router *Router
}

type ServerHandle interface {
	ServerHTTP(address string) error
}

// 请求与函数组合+套接字启动
func ServerMux(r *Router, pattern string, f ServerHTTP, m RequestMethod) {
	r.Register(pattern, f, m)
	log.Logger.Info().Msg("server start")
}

func ListenAddr(addr string, r *Router) {
	Server := &ServerArg{
		addr:   addr,
		Router: r,
	}
	err := Server.ServerHTTP(addr)
	if err != nil {
		log.Logger.Fatal().Err(err).Msgf("socket error %v", err)
	}
}
