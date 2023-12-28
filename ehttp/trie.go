package ehttp

import (
	"errors"
)

type RequestMethod int

const (
	GET RequestMethod = iota
	POST
	PUT
	DELETE
)

type Router struct {
	root *treeNode
}

// 初始化路由器
func NewRouter() *Router {
	return &Router{&treeNode{}}
}

// ----------------------------------路由方法节点----------------------------------

// 路由处理函数
type RouterHandler interface {
	Register(pattern string, f ServerHTTP)
	Search(pattern string) (ServerHTTP, error)
}

// 路由树节点
type treeNode struct {
	pattern   rune               //节点值
	EndNode   bool               //终止节点
	childNode map[rune]*treeNode //孩子节点
	hander    ServerHTTP         //函数绑定
	method    RequestMethod      //请求方法
}

// 路由器注册函数
func (r *Router) Register(pattern string, f ServerHTTP) {
	if r.root == nil {
		panic("root node is nil")
	} else if pattern == "" {
		panic("pattern is empty")
	} else if f == nil {
		panic("hander is nil")
	} else {
		//遍历pattern，使用pattern的每个字符来使用前缀树构建路由树
		current := r.root
		for _, v := range pattern {
			if _, ok := current.childNode[v]; !ok {
				current.childNode[v] = &treeNode{
					pattern:   v,
					EndNode:   false,
					childNode: make(map[rune]*treeNode),
					hander:    nil,
					method:    GET,
				}
			}
			current = current.childNode[v]
		}
		current.EndNode = true
		current.hander = f
	}
}

// 路由搜索函数
func (r *Router) Search(pattern string) (ServerHTTP, error) {
	if r.root == nil {
		panic("root node is nil")
	} else if pattern == "" {
		panic("pattern is empty")
	} else {
		current := r.root
		//前缀树搜索，从根节点开始，使用pattern的字符与节点进行比较，直到找到EndNode，返回hander，否则返回错误
		for _, v := range pattern {
			if current.pattern == v {
				current = current.childNode[v]
			} else {
				if current.EndNode {
					return current.hander, nil
				} else {
					return nil, errors.New("no such route")
				}
			}
		}
	}
	return nil, errors.New("no such route")
}
