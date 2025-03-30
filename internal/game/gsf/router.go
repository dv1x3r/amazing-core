package gsf

import "sync"

type HandlerFunc func(w ResponseWriter, r *Request) error

type Middleware func(next HandlerFunc) HandlerFunc

type Router interface {
	HandleFunc(int32, int32, HandlerFunc)
	Lookup(int32, int32) (HandlerFunc, bool)
	Use(...Middleware)
}

type router struct {
	mu          *sync.RWMutex
	handlers    map[int32]map[int32]HandlerFunc
	middlewares []Middleware
}

func NewRouter() Router {
	return &router{
		mu:       &sync.RWMutex{},
		handlers: map[int32]map[int32]HandlerFunc{},
	}
}

func (r *router) HandleFunc(svcClass int32, msgType int32, handler HandlerFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.handlers[svcClass] == nil {
		r.handlers[svcClass] = map[int32]HandlerFunc{}
	}
	r.handlers[svcClass][msgType] = r.chain(handler)
}

func (r *router) Lookup(svcClass int32, msgType int32) (HandlerFunc, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if msgMap, ok := r.handlers[svcClass]; ok {
		h, ok := msgMap[msgType]
		return h, ok
	}
	return nil, false
}

func (r *router) Use(xs ...Middleware) {
	r.middlewares = append(r.middlewares, xs...)
}

func (r *router) chain(h HandlerFunc) HandlerFunc {
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		h = r.middlewares[i](h)
	}
	return h
}
