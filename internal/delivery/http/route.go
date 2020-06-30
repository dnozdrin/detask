package http

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
)

// Router is a router implementation wrapper
type Router struct {
	mux *mux.Router
}

// Routes is a slice of Routes
type Routes []Route

// NewRouter is a router wrapper constructor
func NewRouter() *Router {
	router := &Router{
		mux: mux.NewRouter().StrictSlash(true),
	}

	return router
}

// Route is the model for the router setup
type Route struct {
	Pattern     string
	Method      string
	Name        string
	HandlerFunc http.HandlerFunc
}

// GetSubRouter will return a router for the given path prefix
func (r Router) GetSubRouter(pathPrefix string) *Router {
	return &Router{mux: r.mux.PathPrefix(pathPrefix).Subrouter()}
}

// Register will add a new route to the router
func (r Router) Register(route Route) {
	r.mux.Methods(route.Method).Path(route.Pattern).HandlerFunc(route.HandlerFunc).Name(route.Name)
}

// ServeHTTP is a wrapper for a concrete router ServeHTTP method
func (r Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

// GetURL will return an URL for the given route name with given parameters
// or an error
func (r Router) GetURL(routeName string, params ...string) (*url.URL, error) {
	return r.mux.Get(routeName).URL(params...)
}

// GetIDVar will return ID var that was set for the current route or an error
func (r Router) GetIDVar(req *http.Request) (uint, error) {
	ID, err := strconv.Atoi(mux.Vars(req)["id"])

	return uint(ID), err
}
