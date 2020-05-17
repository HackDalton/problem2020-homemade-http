package router

import (
	"net/http"
	"strings"
)

// Router is a simple router
type Router struct {
	routes []Route
}

func CreateRouter() Router {
	return Router{make([]Route, 0)}
}

func (rtr *Router) AddRoute(r Route) {
	rtr.routes = append(rtr.routes, r)
}

func (rtr *Router) Handle(req http.Request) http.Response {
	rw := NewResponseWriter()

	if !strings.HasSuffix(req.RequestURI, "/") {
		req.RequestURI = req.RequestURI + "/"
	}

	for _, route := range rtr.routes {
		if strings.HasPrefix(req.RequestURI, route.route) && req.Method == route.method {
			route.handler.ServeHTTP(rw, &req)
			return *rw.resp
		}
	}
	rw.WriteHeader(http.StatusNotFound)
	rw.Write([]byte("Page not found."))
	return *rw.resp
}
