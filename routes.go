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

var routes = Routes{
	Route{
		"Status",
		"GET",
		"/api/status",
		ShowStatus,
	},
	Route{
		"GetDoorCode",
		"GET",
		"/api/settings/doorcode",
		GetDoorCode,
	},
	Route{
		"UpdateDoorCode",
		"PATCH",
		"/api/settings/doorcode",
		UpdateDoorCode,
	},
	Route{
		"ListUsers",
		"GET",
		"/api/users",
		ListUsers,
	},
	Route{
		"GetTransactions",
		"GET",
		"/api/transactions",
		GetTransactions,
	},
	Route{
		"GetTransactionsByUser",
		"GET",
		"/api/transactions/user/{id}",
		GetTransactionsByUser,
	},
	// Route{
	// 	"GetAuthToken",
	// 	"POST",
	// 	"/api/token",
	// 	GetAuthToken,
	// },
}
