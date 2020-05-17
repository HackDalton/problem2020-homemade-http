package router

import "net/http"

// Route is a route that is used by a Router
type Route struct {
	method  string
	route   string
	handler http.Handler
}

// GET creates a GET route
func GET(route string, handler http.Handler) Route {
	return Route{
		method:  "GET",
		route:   route,
		handler: handler,
	}
}

// POST creates a POST route
func POST(route string, handler http.Handler) Route {
	return Route{
		method:  "POST",
		route:   route,
		handler: handler,
	}
}

// New creates a new route
func New(method string, route string, handler http.Handler) Route {
	return Route{method, route, handler}
}
