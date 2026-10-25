package gsf

import "slices"

type HandlerFunc func(w ResponseWriter, r *Request) error

type Middleware func(next HandlerFunc) HandlerFunc

type routeKey struct {
	svcClass int32
	msgType  int32
}

type Router struct {
	handlers    map[routeKey]HandlerFunc
	middlewares []Middleware
}

func NewRouter() *Router {
	return &Router{
		handlers: map[routeKey]HandlerFunc{},
	}
}

func (r *Router) HandleFunc(svcClass int32, msgType int32, handler HandlerFunc) {
	r.handlers[routeKey{svcClass, msgType}] = handler
}

func (r *Router) Use(xs ...Middleware) {
	r.middlewares = append(r.middlewares, xs...)
}

func (r *Router) Lookup(svcClass int32, msgType int32) (HandlerFunc, bool) {
	h, ok := r.handlers[routeKey{svcClass, msgType}]
	if !ok {
		return nil, false
	}
	for _, v := range slices.Backward(r.middlewares) {
		h = v(h)
	}
	return h, true
}
