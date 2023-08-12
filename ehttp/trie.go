package ehttp

import (
	"errors"
)

//初始化
var (
	//路由树根节点
	Root = &RootNode{
		methodNodes: make([]methodNode, 0),
	}
	//路由树终止节点
	ENDNODE = &treeNode{
		pattern:   "",
		EndNode:   false,
		childNode: make([]treeNode, 0),
		hander:    nil,
	}
)

//路由树根节点结构体
type RootNode struct {
	methodNodes []methodNode
}

//路由树方法节点结构体
type methodNode struct {
	method RequestMethod //请求方法
	node   *treeNode     //路由树节点
}

// ----------------------------------路由方法节点----------------------------------
type RequestMethod string

const (
	GET  RequestMethod = "GET"
	POST RequestMethod = "POST"
)

//路由处理函数
type RouterHandler interface {
	Register(pattern string, f ServerHTTP)
	Search(pattern string) (ServerHTTP, error)
}

//路由树节点
type treeNode struct {
	pattern   string     //节点值
	EndNode   bool       //终止节点
	childNode []treeNode //孩子节点
	hander    ServerHTTP //函数绑定
}

// 初始化路由树
func InitTree() {
	Root.methodNodes = append(Root.methodNodes, methodNode{
		method: GET,
		node:   nil,
	})
	Root.methodNodes = append(Root.methodNodes, methodNode{
		method: POST,
		node:   nil,
	})
}

// 前缀树注册
func (root *RootNode) Register(pattern string, f ServerHTTP, method RequestMethod) {
	var currNode *treeNode
	switch method {
	case GET:
		if root.methodNodes[0].node == nil {
			root.methodNodes[0].node = ENDNODE
		}
		currNode = root.methodNodes[0].node

	case POST:
		if root.methodNodes[1].node == nil {
			root.methodNodes[1].node = ENDNODE
		}
		currNode = root.methodNodes[1].node

	default:
		// handle other methods
		return
	}
	for _, v := range pattern {
		found := false
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

// 路由树查找
func (root *RootNode) Search(pattern string, method RequestMethod) (ServerHTTP, error) {
	var currNode *treeNode
	var handler ServerHTTP = nil
	var err error = nil

	switch method {
	case GET:
		currNode = root.methodNodes[0].node
	case POST:
		currNode = root.methodNodes[1].node
	default:
		// handle other methods
	}
	for _, v := range pattern {
		found := false
		for i := 0; i < len(currNode.childNode); i++ {
			if currNode.childNode[i].pattern == string(v) {
				currNode = &currNode.childNode[i]
				found = true
				handler = currNode.hander
			}
		}
		if !found {
			err = errors.New("404,routerERR:invalid router can't find it")
		}
	}
	if !currNode.EndNode {
		err = errors.New("404,routerERR:invalid router")
	}
	return handler, err
}
