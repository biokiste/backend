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
	respondWithHTTP(w, http.StatusOK)
}

// GetUserStates returns user states
func (h Handlers) GetUserStates(w http.ResponseWriter, r *http.Request) {
	var states []string
	states = viper.GetStringSlice("user_states")
	respondWithJSON(w, JSONResponse{Body: &states})
}

// GetTransactionStates returns transaction states
func (h Handlers) GetTransactionStates(w http.ResponseWriter, r *http.Request) {
	var states []string
	states = viper.GetStringSlice("transaction_states")
	respondWithJSON(w, JSONResponse{Body: &states})
}

// GetTransactionTypes return transaction types
func (h Handlers) GetTransactionTypes(w http.ResponseWriter, r *http.Request) {
	var types []string
	types = viper.GetStringSlice("transaction_types")
	respondWithJSON(w, JSONResponse{Body: &types})
}
