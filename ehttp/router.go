package ehttp

import (
	"errors"
	"fmt"
)

type RequestMethod int

const (
	NULL RequestMethod = iota
	GET
	POST
	PUT
	DELETE
)

type Router struct {
	root *treeNode
}

// 初始化路由器
func NewRouter() *Router {
	return &Router{&treeNode{
		pattern:   ' ',
		EndNode:   false,
		childNode: make(map[rune]*treeNode),
		hander:    nil,
		method:    NULL,
	}}
}

// ----------------------------------路由方法节点----------------------------------

// 路由处理函数
type RouterHandler interface {
	Register(pattern string, f ServerHTTP)
	Search(pattern string) (ServerHTTP, error)
	RouterListen()
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
func (r *Router) Register(pattern string, f ServerHTTP, method RequestMethod) {
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
					method:    NULL,
				}
			}
			current = current.childNode[v]
		}
		current.EndNode = true
		current.hander = f
		current.method = method
	}
}

// 路由搜索函数
func (r *Router) Search(pattern string) (ServerHTTP, error) {
	if r.root == nil {
		panic("root node is nil")
	} else if pattern == "" {
		panic("pattern is empty")
	} else {
		//从根节点的子节点开始和pattern进行匹配，因为根节点整个路由树的引导点
		if current, yes := r.root.childNode[[]rune (pattern)[0]]; yes {
			cur := current
			for i := 1; i < len(pattern); i++ {
				v := []rune(pattern)[i]
				if _, ok := cur.childNode[v]; !ok {
					return nil, errors.New("no such route")
				}
				cur = cur.childNode[v]
			}
			if cur.EndNode {
				return cur.hander, nil
			} else {
				return nil, errors.New("no such route")
			}
		}

	}
	return nil, errors.New("no such route")
}

func (router *Router) RouterListener() {
	fmt.Println("infor:buffer Read In")
	for read := range ReadQueen {
		ReadRequest(router, read.Reader)
	}
}
