package main

import (
	"github.com/didi/gendry/scanner"
	"github.com/spf13/viper"
)

// UserStates returns possible user states
func (h Handlers) UserStates() ([]string, error) {
	var states []string
	states = viper.GetStringSlice("userstate")

	return states, nil
}

// TransactionStates delivers transaction states
func (h Handlers) TransactionStates() ([]TransactionState, error) {
	var transactionStates []TransactionState
	results, err := h.DB.Query(`
		SELECT id, type
		FROM transactions_status
	`)
	if err != nil {

		return transactionStates, err
	}
	defer results.Close()

	scanner.Scan(results, &transactionStates)
	if err != nil {
		return transactionStates, err
	}

	return transactionStates, nil
}
