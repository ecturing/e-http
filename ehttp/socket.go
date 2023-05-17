package ehttp

import (
	"fmt"
	"net"
)

var (
	connBuf = make([]byte, 1024)
)

// 定义socket绑定点，并将socket交给Linux epoll管理，增强并发率

func Init_Socket(address string) error {
	fmt.Println("socket init")
	tcpadd, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return err
	}
	listener, err := net.ListenTCP("tcp", tcpadd)
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
		fmt.Println("one TCP join the server")
		go readBuf(conn)
	}
}

func readBuf(c *net.TCPConn) {
	for {
		n,err:=c.Read(connBuf)
		if n>0 {
			ReadChannel <- connBuf
		}else{
			fmt.Println("error:",err)
		}
	}
}

func writeBuf(c *net.TCPConn)  {
	c.Write(connBuf)
	defer c.Close()
}