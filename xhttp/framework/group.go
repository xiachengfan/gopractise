package framework

import (
	"net/http"
	"path/filepath"
)

type IGroup interface {
	Get(string, HandlerFunc)
	Post(string, HandlerFunc)
	Put(string, HandlerFunc)
	Delete(string, HandlerFunc)
}

type RGroup struct {
	Handlers HandlerFunc
	core     *Core  // all groups share a Engine instance
	parent   string //指向上一个Group的路径，如果有的话
	root     bool
}

func (group *RGroup) Group(relativePath string) *RGroup {
	engine := group.core
	newGroup := &RGroup{
		parent: "" + relativePath,
		core:   engine,
	}
	return newGroup
}

func (group *RGroup) calculateAbsolutePath(relativePath string) string {
	return filepath.Join(group.parent, relativePath)
}

func (group *RGroup) ParentPath() string {
	return group.parent
}

func (group *RGroup) handle(httpMethod, relativePath string, handlers HandlerFunc) IGroup {
	absolutePath := group.calculateAbsolutePath(relativePath)
	group.core.addRoute(httpMethod, absolutePath, handlers)
	return group.returnObj()
}

// 实现Get方法
func (g *RGroup) Get(relativePath string, handler HandlerFunc) {
	g.handle(http.MethodGet, relativePath, handler)
}

// 实现Post方法
func (g *RGroup) Post(relativePath string, handler HandlerFunc) {
	g.handle(http.MethodPost, relativePath, handler)
}

// 实现Put方法
func (g *RGroup) Put(relativePath string, handler HandlerFunc) {
	g.handle(http.MethodPut, relativePath, handler)
}

// 实现Delete方法
func (g *RGroup) Delete(relativePath string, handler HandlerFunc) {
	g.handle(http.MethodDelete, relativePath, handler)
}

func (group *RGroup) returnObj() IGroup {
	if group.root {
		return group.core
	}
	return group
}
