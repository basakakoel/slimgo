package slimgo

import (
	"strings"
)

type RouterTree struct {
	subrouters map[string]*RouterTree
	leaves     []*routerLeaf
}

type routerLeaf struct {
	toRun interface{} //run info (controllerinfo)
}

//new tree
func NewRouterTree() *RouterTree {
	return &RouterTree{
		subrouters: make(map[string]*RouterTree),
	}
}

//add tree base
func (this *RouterTree) addTreeBase(pattern []string, tree *RouterTree) {
	if len(pattern) == 0 {
		panic("No pattern to add tree")
	}
	node := pattern[0]
	if len(pattern) == 1 {
		this.subrouters[node] = tree
		return
	}
	subTree := NewRouterTree()
	this.subrouters[node] = subTree
	subTree.addTreeBase(pattern[1:], tree)
}

//add tree
func (this *RouterTree) AddTree(pattern string, tree *RouterTree) {
	this.addTreeBase(pattern2node(pattern), tree)
}

//add router base
func (this *RouterTree) addRouterBase(nodes []string, torun interface{}) {
	if len(nodes) == 0 {
		this.leaves = append(this.leaves, &routerLeaf{toRun: torun})
		return
	}
	node := nodes[0]
	subtree, ok := this.subrouters[node]
	if !ok {
		subtree = NewRouterTree()
		this.subrouters[node] = subtree
	}
	subtree.addRouterBase(nodes[1:], torun)
}

//add router
func (this *RouterTree) AddRouter(pattern string, torun interface{}) {
	this.addRouterBase(pattern2node(pattern), torun)
}

func (this *RouterTree) FindRouter(pattern string) interface{} {
	if len(pattern) == 0 || pattern[0] != '/' {
		return nil
	}
	return this.findRouterBase(pattern2node(strings.ToLower(pattern)))
}

func (this *RouterTree) findRouterBase(nodes []string) interface{} {
	if len(nodes) == 0 {
		//leaves
		if this.leaves == nil || len(this.leaves) == 0 {
			subtree, ok := this.subrouters["index"]
			if ok {
				return subtree.findRouterBase([]string{})
			}
			return nil
		} else {
			return this.leaves[0].toRun
		}
	}
	nodeNow := nodes[0]
	subtree, ok := this.subrouters[nodeNow]
	if ok {
		return subtree.findRouterBase(nodes[1:])
	} else {
		subtree, ok = this.subrouters["index"]
		if ok {
			return subtree.findRouterBase(nodes[1:])
		}
	}
	return nil
}

func pattern2node(pattern string) []string {
	if pattern == "" {
		return []string{}
	}
	nodes := strings.Split(pattern, "/")
	if nodes[0] == "" {
		nodes = nodes[1:]
	}
	if nodes[len(nodes)-1] == "" {
		nodes = nodes[:len(nodes)-1]
	}
	return nodes
}
