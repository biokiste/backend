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
		// Route{
		// 	"GetUserById",
		// 	"GET",
		// 	"/api/users/{id}",
		// 	h.GetUser,
		// },
		Route{
			"GetUserByEmail",
			"GET",
			"/users/{email}",
			h.GetUserByEmail,
		},
		Route{
			"CreateAuth0User",
			"POST",
			"/user/auth/create",
			h.CreateUser,
		},
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
			"GetTransactions",
			"GET",
			"/transactions",
			h.GetTransactions,
		},
		Route{
			"GetTransactionsByUser",
			"GET",
			"/transactions/user/{id}",
			h.GetTransactionsByUser,
		},
		Route{
			"AddTransaction",
			"POST",
			"/transaction",
			h.AddTransaction,
		},
		Route{
			"UpdatePayment",
			"PATCH",
			"/transactions/payments",
			h.UpdatePayment,
		},
		Route{
			"GetOpenPayments",
			"GET",
			"/transactions/payments/open",
			h.GetOpenPayments,
		},
		Route{
			"GetGroupTypes",
			"GET",
			"/group/types",
			h.GetGroupTypes,
		},
		Route{
			"GetGroups",
			"GET",
			"/groups",
			h.GetGroups,
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
		Route{
			"GetSettings",
			"GET",
			"/settings",
			h.GetSettings,
		},
		Route{
			"GetSettingByKey",
			"GET",
			"/settings/{key}",
			h.GetSettingByKey,
		},
		Route{
			"UpdateSettingByKey",
			"PATCH",
			"/settings/{key}",
			h.UpdateSettingByKey,
		},
		Route{
			"AddSetting",
			"POST",
			"/settings",
			h.AddSetting,
		},
		// Route{
		// 	"SendMail",
		// 	"GET",
		// 	"/api/send/mail",
		// 	h.SendMail,
		// },
	}
}
