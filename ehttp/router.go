package ehttp

import (
	"ews/Eerror"
	"ews/logutil"
)

// 请求方法重定义
type RequestMethod int

// 方法枚举
const (
	NULL RequestMethod = iota //默认空
	GET
	POST
	PUT
	DELETE
)

// 路由
type Router struct {
	root *treeNode
}

// 初始化路由器,根节点为引导节点，不做存储
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
	Register(pattern string, f ServerHTTP, method RequestMethod)
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
	defer func() {
		if err := recover(); err != nil {
			logutil.Logger.Error().Err(err.(error)).Msg("Register error")
		}
	}()

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
func (r *Router) Search(pattern string, method RequestMethod) (ServerHTTP, error) {

	defer func() {
		if err := recover(); err != nil {
			logutil.Logger.Error().Err(err.(error)).Msg("Search error")
		}
	}()

	if r.root == nil {
		panic("root node is nil")
	} else if pattern == "" {
		panic("pattern is empty")
	} else {
		//从根节点的子节点开始和pattern进行匹配，因为根节点整个路由树的引导点
		if current, yes := r.root.childNode[[]rune(pattern)[0]]; yes {
			cur := current
			for i := 1; i < len(pattern); i++ {
				v := []rune(pattern)[i]
				if _, ok := cur.childNode[v]; !ok {
					return nil, Eerror.NotFound
				}
				cur = cur.childNode[v]
			}
			if cur.EndNode && cur.method == method {
				return cur.hander, nil
			} else if !cur.EndNode {
				return nil, Eerror.NotFound

			} else {
				return nil, Eerror.MethodNotAllow
			}
		}

	}
	logutil.Logger.Error().Err(Eerror.NotFound).Msg("Search error")
	return nil, Eerror.NotFound
}

// 路由监听函数
func (r *Router) RouterListen() {
	logutil.Logger.Info().Msg("Starting Router Listening...")
	for read := range ReadQueen {
		ReadRequest(r, read.Reader)
	}
}
