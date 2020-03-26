package main

import (
	"database/sql"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

// Handlers wrapps DB instance
type Handlers struct {
	DB *sql.DB
}

// ShowStatus delivers actual status
func (h Handlers) ShowStatus(w http.ResponseWriter, r *http.Request) {
	printJSON(w, "ok")
}

// GetUserStates returns user states
func (h Handlers) GetUserStates(w http.ResponseWriter, r *http.Request) {
	var states []string
	states = viper.GetStringSlice("user_states")
	printJSON(w, &states)
}

// GetTransactionStates returns transaction states
func (h Handlers) GetTransactionStates(w http.ResponseWriter, r *http.Request) {
	var states []string
	states = viper.GetStringSlice("transaction_states")
	printJSON(w, &states)
}

// GetTransactionTypes return transaction types
func (h Handlers) GetTransactionTypes(w http.ResponseWriter, r *http.Request) {
	var types []string
	types = viper.GetStringSlice("transaction_types")
	printJSON(w, &types)
}
