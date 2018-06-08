package simplerouter

import (
	"net/http"
)

// Param ...
type Param struct {
	Key   string
	Value string
}

// Params ...
type Params []Param

// Handle ...
type Handle func(res http.ResponseWriter, req *http.Request, p Params)

// Router ...
type Router struct {
	tree *node
}

// GetValue ...
func (p Params) GetValue(key string) (string, bool) {
	for _, param := range p {
		if key == param.Key {
			return param.Value, true
		}
	}
	return "", false
}

func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	handler, params, found := r.tree.search(req.Method, req.URL.Path)
	if found {
		handler(res, req, params)
	} else {
		res.WriteHeader(http.StatusNotFound)
	}
}

// Get ...
func (r *Router) Get(path string, handle Handle) {
	r.Register("GET", path, handle)
}

// Post ...
func (r *Router) Post(path string, handle Handle) {
	r.Register("POST", path, handle)
}

// Put ...
func (r *Router) Put(path string, handle Handle) {
	r.Register("Put", path, handle)
}

// Delete ...
func (r *Router) Delete(path string, handle Handle) {
	r.Register("Delete", path, handle)
}

// Register method and path with a handler
func (r *Router) Register(method string, path string, handle Handle) {
	if path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}

	r.tree.addNode(method, path, handle)

}

// New ...
func New() *Router {
	node := node{component: "/", isNamedParam: false, methods: make(map[string]Handle)}
	return &Router{tree: &node}
}
