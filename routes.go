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
		"ListUsers",
		"GET",
		"/api/users",
		ListUsers,
	},
	Route{
		"GetPayments",
		"GET",
		"/api/payments",
		GetPayments,
	},
	Route{
		"GetPaymentsByUser",
		"GET",
		"/api/payments/user/{id}",
		GetPaymentsByUser,
	},
	// Route{
	// 	"GetAuthToken",
	// 	"POST",
	// 	"/api/token",
	// 	GetAuthToken,
	// },
}
