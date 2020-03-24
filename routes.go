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
			"ListUsers",
			"GET",
			"/users",
			h.ListUsers,
		},
		Route{
			"LastActiveUsers",
			"GET",
			"/users/lastactive",
			h.LastActiveUsers,
		},
		Route{
			"GetUserByEmail",
			"GET",
			"/users/{email}",
			h.GetUserByEmail,
		},
		// Route{
		// 	"CreateAuth0User",
		// 	"POST",
		// 	"/user/auth/create",
		// 	h.CreateUser,
		// },
		Route{
			"GetAuth0User",
			"GET",
			"/user/auth/{id}",
			h.GetAuth0User,
		},
		Route{
			"UpdateUser",
			"PATCH",
			"/user",
			h.UpdateUser,
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
		// Route{
		// 	"SendMail",
		// 	"GET",
		// 	"/api/send/mail",
		// 	h.SendMail,
		// },
	}
}
