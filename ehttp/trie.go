package ehttp

import "fmt"

// 前缀树路由匹配，生成Get树与Post树

//路由错误
type routerError struct{
	msg string
}

type RouterHandler interface {
	Register(pattern string, f ServerHTTP)
	Search(pattern string) (ServerHTTP,error)
}

type treeNode struct {
	pattern   string           //节点值
	EndNode   bool             //终止节点
	childNode []treeNode       //孩子节点
	hander    ServerHTTP //函数绑定
}


func (err *routerError) Error() string{
	return fmt.Sprintf("%v",err.msg)
}

// 前缀树注册
func (root *treeNode) Register(pattern string, f ServerHTTP){
	currNode:=root
	for _,v := range pattern {
		found:=false
		for i := 0; i < len(currNode.childNode); i++ {
			if currNode.childNode[i].pattern == string(v) {
				currNode = &currNode.childNode[i]
				found = true
				break
			}
		}
		if !found {
			newNode := treeNode{
				pattern:   string(v),
				EndNode:   false,
				childNode: make([]treeNode, 0),
				hander:    nil,
			}
			currNode.childNode = append(currNode.childNode, newNode)
			currNode = &currNode.childNode[len(currNode.childNode)-1]
		}
	}
	currNode.EndNode = true
	currNode.hander = f
}

//路由树查找
func (root *treeNode) Search(pattern string) (ServerHTTP,error ){
	currNode:=root
	var handler ServerHTTP=nil
	var err *routerError=nil
	for _, v := range pattern {
		found:=false
		for i := 0; i < len(currNode.childNode); i++ {
			if currNode.childNode[i].pattern == string(v) {
				currNode = &currNode.childNode[i]
				found = true
				handler=currNode.hander
			}
		}
		if !found {
			err=&routerError{
				msg: "invalid router can't find it",
			}
		}
	}
	if !currNode.EndNode {
		err=&routerError{
			msg: "invalid router!",
		}
	}
	return handler,err
}
