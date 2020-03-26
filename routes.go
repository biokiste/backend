package main

import (
	"net/http"
)

// Route implements route struct
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes holds routes
type Routes []Route

// GetRoutes defines routes
func GetRoutes(h *Handlers) Routes {
	return Routes{
		Route{
			"Status",
			"GET",
			"/status",
			h.ShowStatus,
		},
		Route{
			"GetUserStates",
			"GET",
			"/states/user",
			h.GetUserStates,
		},
		Route{
			"GetTransactionState",
			"GET",
			"/states/transaction",
			h.GetTransactionStates,
		},
		Route{
			"GetTransactionState",
			"GET",
			"/states/loan",
			h.GetTransactionStates,
		},
		Route{
			"GetTransactionTypes",
			"GET",
			"/types/transaction",
			h.GetTransactionTypes,
		},
	}
}
