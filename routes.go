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
			"LastActiveUsers",
			"GET",
			"/api/users/lastactive",
			h.LastActiveUsers,
		},
		// Route{
		// 	"GetUserById",
		// 	"GET",
		// 	"/api/users/{id}",
		// 	h.GetUser,
		// },
		Route{
			"GetUserByEmail",
			"GET",
			"/api/users/{email}",
			h.GetUserByEmail,
		},
		Route{
			"CreateAuth0User",
			"POST",
			"/api/user/auth/create",
			h.CreateUser,
		},
		Route{
			"GetAuth0User",
			"GET",
			"/api/user/auth/{id}",
			h.GetAuth0User,
		},
		Route{
			"UpdateUser",
			"PATCH",
			"/api/user",
			h.UpdateUser,
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
		Route{
			"UpdatePayment",
			"PATCH",
			"/api/transactions/payments",
			h.UpdatePayment,
		},
		Route{
			"GetOpenPayments",
			"GET",
			"/api/transactions/payments/open",
			h.GetOpenPayments,
		},
		Route{
			"GetGroupTypes",
			"GET",
			"/api/group/types",
			h.GetGroupTypes,
		},
		Route{
			"GetGroups",
			"GET",
			"/api/groups",
			h.GetGroups,
		},
		Route{
			"GetUserStates",
			"GET",
			"/api/states/users",
			h.GetUserStates,
		},
		Route{
			"GetTransactionState",
			"GET",
			"/api/states/transaction",
			h.GetTransactionStates,
		},
		// Route{
		// 	"SendMail",
		// 	"GET",
		// 	"/api/send/mail",
		// 	h.SendMail,
		// },
	}
}
