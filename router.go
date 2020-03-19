package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

// RouterConfig implements handlers
type RouterConfig struct {
	Handlers *Handlers
}

// APIRouter registers routes
func APIRouter(config *RouterConfig) *mux.Router {
	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api").Subrouter()

	// load routes from routes.go
	for _, route := range GetRoutes(config.Handlers) {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Auth(Logger(handler, route.Name))

		apiRouter.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./assets/"))))

	return router
}
