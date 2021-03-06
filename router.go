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
	subRouter := router.PathPrefix("/api").Subrouter()

	var routes []Route

	for _, r := range GetSettingsRoutes(config.Handlers) {
		routes = append(routes, r)
	}

	for _, r := range GetLoansRoutes(config.Handlers) {
		routes = append(routes, r)
	}

	for _, r := range GetTransactionsRoutes(config.Handlers) {
		routes = append(routes, r)
	}

	for _, r := range GetGroupsRoutes(config.Handlers) {
		routes = append(routes, r)
	}

	for _, r := range GetUsersRoutes(config.Handlers) {
		routes = append(routes, r)
	}

	for _, r := range GetAssetsRoutes(config.Handlers) {
		routes = append(routes, r)
	}

	for _, r := range GetRoutes(config.Handlers) {
		routes = append(routes, r)
	}

	for _, r := range routes {
		var h http.Handler
		h = r.HandlerFunc
		h = Auth(Logger(h, r.Name))
		subRouter.
			Methods(r.Method).
			Path(r.Pattern).
			Name(r.Name).
			Handler(h)
	}

	// load routes from routes.go

	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./assets/"))))

	return router
}
