package ehttp

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
)

type Event struct {
	Conn   *net.TCPConn
	Reader *bufio.Reader
	Writer *bytes.Buffer
}

// 定义socket绑定点，并将socket交给Linux epoll管理，增强并发率
func InitSocket(address string) error {
	fmt.Println("socket init")
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return err
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}
	defer listener.Close()
	
	listenerHandler(listener)
	return nil
}

func listenerHandler(listener *net.TCPListener) {
	fmt.Println("TCP LISTENING")
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Printf("Error accepting connection: %s\n", err.Error())
			continue
		}
		fmt.Printf("%s TCP join the server\n",conn.RemoteAddr().String())
		go connectionReadBuf(conn)
		go connectWriteBuf(conn)
	}
}

// 读取链接缓冲区准备读取的事件并发送到管道
func connectionReadBuf(c *net.TCPConn) {
		readevent := bufio.NewReader(c)
		var event = &Event{
			Conn: c, 
			Reader: readevent, 
			Writer: nil,
		}
		ReadQueen <- event
}

// 写入链接的缓冲区
func connectWriteBuf(c *net.TCPConn) {
	for s := range WriteQueen {
		w := bufio.NewWriter(c)
		w.Write(s.Writer.Bytes())
		w.Flush()
		c.Close()
		fmt.Printf("%s TCP leave the server\n",c.RemoteAddr().String())
	}
}
