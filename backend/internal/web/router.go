package web

import (
	"fmt"
	"nekolog/internal/log"
	"net/http"
)

type HandlerFunc func(*Context)

type RouterGroup struct {
	path        string
	middlewares []HandlerFunc
	router      *Router
}

type Router struct {
	RouterGroup
	routes       map[string]map[string][]HandlerFunc
	sessionStore *SessionStore
}

func NewRouter() *Router {
	r := &Router{
		routes: map[string]map[string][]HandlerFunc{
			"GET":    make(map[string][]HandlerFunc),
			"POST":   make(map[string][]HandlerFunc),
			"PATCH":  make(map[string][]HandlerFunc),
			"DELETE": make(map[string][]HandlerFunc),
		},
		sessionStore: NewSessionStore(),
	}

	r.RouterGroup = RouterGroup{
		path:        "",
		middlewares: nil,
		router:      r,
	}

	return r
}

func (g *RouterGroup) Use(middlewares ...HandlerFunc) {
	g.middlewares = append(g.middlewares, middlewares...)
}

func (g *RouterGroup) Group(relativePath string, middlewares ...HandlerFunc) *RouterGroup {
	return &RouterGroup{
		path:        g.path + relativePath,
		middlewares: append(g.middlewares, middlewares...),
		router:      g.router,
	}
}

func (g *RouterGroup) addRoute(method, relativePath string, handler HandlerFunc) {
	absolutePath := g.path + relativePath

	var chain []HandlerFunc
	chain = append(chain, g.middlewares...)
	chain = append(chain, handler)

	g.router.routes[method][absolutePath] = chain
}

func (g *RouterGroup) GET(path string, handler HandlerFunc)    { g.addRoute("GET", path, handler) }
func (g *RouterGroup) POST(path string, handler HandlerFunc)   { g.addRoute("POST", path, handler) }
func (g *RouterGroup) PATCH(path string, handler HandlerFunc)  { g.addRoute("PATCH", path, handler) }
func (g *RouterGroup) DELETE(path string, handler HandlerFunc) { g.addRoute("DELETE", path, handler) }

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handlers, exists := r.routes[req.Method]
	if !exists {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	chain, exists := handlers[req.URL.Path]
	if !exists {
		http.Error(w, "Handler Not Found", http.StatusNotFound)
		return
	}

	c := &Context{
		Writer:    w,
		Request:   req,
		store:     r.sessionStore,
		session:   nil,
		sessionID: "",
		handlers:  chain,
		index:     -1,
	}

	c.Next()
}

func (r *Router) Serve(addr string) error {
	log.Info("Register routes")

	for method, route := range r.routes {
		for route_name, route_handler_chain := range route {
			log.Info(
				"%-10s %-40s → %d handlers",
				fmt.Sprintf("[%s]", method),
				route_name,
				len(route_handler_chain),
			)
		}
	}

	log.Info("Server listening on %s", addr)

	return http.ListenAndServe(addr, r)
}
