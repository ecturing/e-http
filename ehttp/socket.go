package ehttp

import (
	"bufio"
	"bytes"
	"context"
	e_context "ews/context"
	"ews/log"
	"net"
	"time"
)

const (
	DEADLINETIME = 10 * time.Second
)

type ConnClose bool

type Event struct {
	// Conn   *net.TCPConn
	Reader *bufio.Reader
	Writer *bytes.Buffer
}

type ServerCtx struct {
	conn      *net.TCPConn
	serverArg ServerArg
}

// 定义socket绑定点，并将socket交给Linux epoll管理，增强并发率
func (s *ServerArg) ServerHTTP(address string) error {
	log.Logger.Info().Msgf("Starting TCP server on %s", address)
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return err
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}
	defer listener.Close()

	listenerHandler(s, listener)
	return nil
}

func listenerHandler(s *ServerArg, listener *net.TCPListener) {
	log.Logger.Info().Msg("Tcp listener start")
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Logger.Error().Err(err).Msg("Error accepting connection")
			continue
		}
		log.Logger.Info().Msgf("Accepted connection from %s", conn.RemoteAddr().String())
		// conn.SetDeadline(time.Now().Add(DEADLINETIME))
		ctx := context.WithValue(e_context.RootCtx, "conn", &ServerCtx{conn, *s})
		// todo
		// 未来考虑采用IO多路复用技术
		go connectionReadBuf(ctx)
		go connectWriteBuf(ctx)
	}
}

// 读取链接缓冲区准备读取的事件并发送到管道
func connectionReadBuf(ctx context.Context) {
	serverCtx := ctx.Value("conn").(*ServerCtx)
	c := serverCtx.conn
	r := serverCtx.serverArg.Router
	readevent := bufio.NewReader(c)
	for {
		ok := ReadRequest(r, readevent)
		if ok {
			CloseQueen <- struct{}{}
			break
		}
	}
}

// 写入链接的缓冲区
func connectWriteBuf(ctx context.Context) {
	sctx := ctx.Value("conn").(*ServerCtx)
	for {
		select {
		case s := <-WriteQueen:
			// do
			w := bufio.NewWriter(sctx.conn)
			w.Write(s.Writer.Bytes())
			w.Flush()
			log.Logger.Info().Msgf("%s TCP the server return msg", sctx.conn.RemoteAddr().String())
		case s := <-ErrorSingal:
			// do
			w := bufio.NewWriter(sctx.conn)
			w.Write(s.Writer.Bytes())
			w.Flush()
			sctx.conn.Close()
			log.Logger.Info().Msgf("%s TCP leave the server", sctx.conn.RemoteAddr().String())
		case <-CloseQueen:
			return
		}
	}
}
