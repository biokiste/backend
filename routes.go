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
			"/api/status",
			h.ShowStatus,
		},
		Route{
			"GetDoorCode",
			"GET",
			"/api/settings/doorcode",
			h.GetDoorCode,
		},
		Route{
			"UpdateDoorCode",
			"PATCH",
			"/api/settings/doorcode",
			h.UpdateDoorCode,
		},
		Route{
			"ListUsers",
			"GET",
			"/api/users",
			h.ListUsers,
		},
		Route{
			"GetTransactions",
			"GET",
			"/api/transactions",
			h.GetTransactions,
		},
		Route{
			"GetTransactionTypes",
			"GET",
			"/api/transactions/types",
			h.GetTransactionTypes,
		},
		Route{
			"GetTransactionsByUser",
			"GET",
			"/api/transactions/user/{id}",
			h.GetTransactionsByUser,
		},
		Route{
			"AddTransaction",
			"POST",
			"/api/transaction",
			h.AddTransaction,
		},
	}
}
