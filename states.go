package main

import (
	"github.com/didi/gendry/scanner"
)

// UserStates returns possible user states
func (h Handlers) UserStates() ([]UserState, error) {
	var states []UserState
	results, err := h.DB.Query(`
		SELECT id, description
		FROM user_states`)
	if err != nil {
		return states, err
	}

	defer results.Close()

	err = scanner.Scan(results, &states)
	if err != nil {
		return states, err
	}

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
