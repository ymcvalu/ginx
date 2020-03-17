package ginx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
)

type Router interface {
	Use(...gin.HandlerFunc) Router
	Group(string, ...gin.HandlerFunc) Router
	Any(string, ...interface{}) Router
	GET(string, ...interface{}) Router
	POST(string, ...interface{}) Router
	DELETE(string, ...interface{}) Router
	PATCH(string, ...interface{}) Router
	PUT(string, ...interface{}) Router
	OPTIONS(string, ...interface{}) Router
	HEAD(string, ...interface{}) Router
	StaticFile(string, string) Router
	Static(string, string) Router
	StaticFS(string, http.FileSystem) Router
}

type group struct {
	r  Renderer
	rg *gin.RouterGroup
}

func (g *group) Use(ms ...gin.HandlerFunc) Router {
	g.rg.Use(ms...)
	return g
}

func (g *group) Group(p string, ms ...gin.HandlerFunc) Router {
	rg := g.rg.Group(p, ms...)
	return &group{
		rg: rg,
	}
}

func (g *group) Any(p string, hs ...interface{}) Router {
	g.rg.Any(p, g.handlers(hs)...)
	return g
}

func (g *group) GET(p string, hs ...interface{}) Router {
	g.rg.GET(p, g.handlers(hs)...)
	return g
}

func (g *group) POST(p string, hs ...interface{}) Router {
	g.rg.POST(p, g.handlers(hs)...)
	return g
}

func (g *group) DELETE(p string, hs ...interface{}) Router {
	g.rg.DELETE(p, g.handlers(hs)...)
	return g
}

func (g *group) PATCH(p string, hs ...interface{}) Router {
	g.rg.PATCH(p, g.handlers(hs)...)
	return g
}

func (g *group) PUT(p string, hs ...interface{}) Router {
	g.rg.PUT(p, g.handlers(hs)...)
	return g
}

func (g *group) OPTIONS(p string, hs ...interface{}) Router {
	g.rg.OPTIONS(p, g.handlers(hs)...)
	return g
}

func (g *group) HEAD(p string, hs ...interface{}) Router {
	g.rg.HEAD(p, g.handlers(hs)...)
	return g
}

func (g *group) StaticFile(p string, fpath string) Router {
	g.rg.StaticFile(p, fpath)
	return g
}

func (g *group) Static(p string, root string) Router {
	g.rg.Static(p, root)
	return g
}

func (g *group) StaticFS(p string, fs http.FileSystem) Router {
	g.rg.StaticFS(p, fs)
	return g
}

func XRouter(h *gin.Engine, r Renderer) Router {
	if r == nil {
		r = defRenderer{}
	}
	return &group{
		r:  r,
		rg: h.Group("/"),
	}
}

func (g *group) handlers(hs []interface{}) []gin.HandlerFunc {
	n := len(hs)
	if n == 0 {
		return nil
	}

	_hs := make([]gin.HandlerFunc, n)

	for i := 0; i < n-1; i++ {
		if h, ok := hs[i].(gin.HandlerFunc); !ok {
			panic(fmt.Errorf("illegal middleware signature: %s", reflect.TypeOf(hs[i]).String()))
		} else {
			_hs[i] = h
		}
	}

	_hs[n-1] = wrapper(hs[n-1], g.r)
	return _hs
}
