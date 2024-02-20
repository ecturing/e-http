package state

import (
	"ews/ehttp"
	"net"
)


// tcp连接的状态机模型 
type State int

const(
	//  闲置状态
	Idle State =iota 
	// 活动状态
	Active
	// 关闭状态
	Close
)


type ConnState struct{
	c net.Conn 
	connstate State
}
// 初始化状态是闲置状态
func NewState(c net.Conn) *ConnState{
	return &ConnState{c,Idle}
	
}

func (c *ConnState) GetConnState() State {
	return c.connstate
}
func (c *ConnState)GetConnStateConn() net.Conn {
	return c.c
}
// 1. BlockState -> ActiveState
func (c *ConnState) BlockStateToActiveState(){
	//to do
	c.connstate = Active
}

// 2. ActiveState -> CloseState
func (c *ConnState) ActiveStateToCloseState(){
	//to do
	req:=&ehttp.E_Response{}
	req.SetConnClose()
	c.connstate = Close
}

// 3. CloseState -> BlockState
func (c *ConnState) BlockStateToCloseState() {
	//to do
	c.connstate = Close
}
