package ehttp

import (
	"bufio"
	"fmt"
	"net"
)

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
	go RouterListener() //Socket启动成功，启动路由系统
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
		reader := bufio.NewReader(c)
		TCPBuffer <- reader
	}
}
