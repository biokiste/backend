package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// Handlers wrapps DB instance
type Handlers struct {
	DB *sql.DB
}

// ShowStatus delivers actual status
func (h Handlers) ShowStatus(w http.ResponseWriter, r *http.Request) {
	printJSON(w, "ok")
}

// GetDoorCode delivers actual door codes
func (h Handlers) GetDoorCode(w http.ResponseWriter, r *http.Request) {
	var doorCode DoorCode
	if err := h.DB.QueryRow(
		`SELECT value 
		 FROM settings
		 WHERE id = ?`, 1).Scan(&doorCode.Value); err != nil {
		printError(w, err)
	}

	printJSON(w, &DoorCodeResponse{DoorCode: doorCode})
}

// UpdateDoorCode sets doorcode
func (h Handlers) UpdateDoorCode(w http.ResponseWriter, r *http.Request) {
	var updatedDoorCode DoorCode
	err := json.NewDecoder(r.Body).Decode(&updatedDoorCode)
	if err != nil {
		printError(w, err.Error())
	}

	_, err = h.DB.Exec(
		`UPDATE settings
	   SET value = ?, updated_at = ?, updated_by = ? 		 
		 WHERE id = 1`, updatedDoorCode.Value, updatedDoorCode.UpdatedAt, updatedDoorCode.UpdatedBy)

	if err != nil {
		fmt.Println(err.Error())
		printError(w, err)
	} else {
		printSuccess(w)
	}

}

// ListUsers delivers user data
func (h Handlers) ListUsers(w http.ResponseWriter, r *http.Request) {
	results, err := h.DB.Query("SELECT id, username FROM users")
	if err != nil {
		printError(w, err.Error())
	}

	defer results.Close()

	var users []User
	for results.Next() {
		var user User
		err = results.Scan(&user.ID, &user.Username)
		if err != nil {
			printError(w, err.Error())
		}
		users = append(users, user)
	}

	printJSON(w, &UsersResponse{Users: users})
}

// GetTransactions delivers all payments
func (h Handlers) GetTransactions(w http.ResponseWriter, r *http.Request) {

	results, err := h.DB.Query(`
		SELECT transactions.id, amount, transactions.created_at, firstname, lastname, transactions.status, transactions.reason, category_id, transactions_category.type
		FROM transactions
		LEFT JOIN transactions_category ON transactions.category_id = transactions_category.id
		LEFT JOIN users ON transactions.user_id = users.id
		WHERE transactions_category.type = "payment"
		AND firstname IS NOT NULL
		OR transactions_category.type = "paymentSepa" AND firstname IS NOT NULL
		OR transactions.status > 1 AND firstname IS NOT NULL
		ORDER BY transactions.created_at desc
	  `)
	if err != nil {
		printError(w, err.Error())
	}

	defer results.Close()
	var transactions []Transaction
	for results.Next() {
		var transaction Transaction

		err = results.Scan(
			&transaction.ID,
			&transaction.Amount,
			&transaction.CreatedAt,
			&transaction.FirstName,
			&transaction.LastName,
			&transaction.Status,
			&transaction.Reason,
			&transaction.CategoryID,
			&transaction.Type)
		if err != nil {
			printError(w, err.Error())
		}
		transactions = append(transactions, transaction)
	}

	printJSON(w, &TransactionResponse{Transactions: transactions})

}

// GetTransactionsByUser delivers payments per user
func (h Handlers) GetTransactionsByUser(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	results, err := h.DB.Query(`
		SELECT transactions.id, amount, transactions.created_at, firstname, lastname, transactions.status, transactions.reason, category_id, transactions_category.type
		FROM transactions
		LEFT JOIN transactions_category ON transactions.category_id = transactions_category.id
		LEFT JOIN users ON transactions.user_id = users.id
		WHERE users.id = ?
		ORDER BY transactions.created_at desc
	  `, id)
	if err != nil {
		printError(w, err.Error())
	}

	defer results.Close()
	var transactions []Transaction
	for results.Next() {
		var transaction Transaction

		err = results.Scan(
			&transaction.ID,
			&transaction.Amount,
			&transaction.CreatedAt,
			&transaction.FirstName,
			&transaction.LastName,
			&transaction.Status,
			&transaction.Reason,
			&transaction.CategoryID,
			&transaction.Type)
		if err != nil {
			printError(w, err.Error())
		}
		transactions = append(transactions, transaction)
	}

	userBalance, err := h.GetBalance(id)
	if err != nil {
		printError(w, err)
	}

	printJSON(w, &UserTransactionResponse{
		UserTransaction: UserTransaction{
			Transactions: transactions,
			Balance:      userBalance},
	})
}

// GetTransactionTypes delivers transaction categories
func (h Handlers) GetTransactionTypes(w http.ResponseWriter, r *http.Request) {
	results, err := h.DB.Query(`
		SELECT id, type, description
		FROM transactions_category
	`)
	if err != nil {
		printError(w, err.Error())
	}
	defer results.Close()

	var transactionCategories []TransactionCategory
	for results.Next() {
		var transactionCategory TransactionCategory
		err = results.Scan(
			&transactionCategory.ID,
			&transactionCategory.Type,
			&transactionCategory.Description)
		if err != nil {
			printError(w, err.Error())
		}
		transactionCategories = append(transactionCategories, transactionCategory)
	}
	printJSON(w, &TransactionCategoryResponse{TransactionCategories: transactionCategories})
}

// AddTransaction updates user balance with transactions
func (h Handlers) AddTransaction(w http.ResponseWriter, r *http.Request) {
	var transactionRequest TransactionRequest
	err := json.NewDecoder(r.Body).Decode(&transactionRequest)
	if err != nil {
		printError(w, err.Error())
	}

	for _, t := range transactionRequest.Transactions {
		stmt, err := h.DB.Prepare(`
			INSERT INTO transactions(
				user_id,
				category_id,
				amount,
				status,
				created_at,
				reason) VALUES(?,?,?,?,?,?)
		`)
		if err != nil {
			// TODO push error to error[] and send printError
			log.Fatal(err)
		}
		_, err = stmt.Exec(
			transactionRequest.User.ID,
			t.CategoryID,
			t.Amount,
			t.Status,
			t.CreatedAt,
			t.Reason,
		)
		if err != nil {
			// TODO push error to error[] and send printError
			log.Fatal(err)
		}
	}

	balance, err := h.GetBalance(transactionRequest.User.ID)
	printJSON(w, &UserTransactionResponse{
		UserTransaction{Balance: balance},
	})
}

// UpdatePayment updates payment possibly as accepted or rejected
func (h Handlers) UpdatePayment(w http.ResponseWriter, r *http.Request) {
	var transaction TransactionRequest
	err := json.NewDecoder(r.Body).Decode(&transaction)
	if err != nil {
		printError(w, err.Error())
	}
	err = h.UpdateTransaction(transaction)
	if err != nil {
		printError(w, err.Error)
	} else {
		printSuccess(w)
	}
}

// GetOpenPayments delivers all open user transactions (payments)
func (h Handlers) GetOpenPayments(w http.ResponseWriter, r *http.Request) {
	results, err := h.DB.Query(`
		SELECT transactions.id, amount, transactions.created_at, firstname, lastname, transactions.status, transactions.reason, category_id, transactions_category.type
		FROM transactions
		LEFT JOIN transactions_category ON transactions.category_id = transactions_category.id
		LEFT JOIN users ON transactions.user_id = users.id
		WHERE transactions.status = 2
		AND firstname IS NOT NULL		
		ORDER BY transactions.created_at desc
	  `)
	if err != nil {
		printError(w, err.Error())
	}

	defer results.Close()
	var transactions []Transaction
	for results.Next() {
		var transaction Transaction

		err = results.Scan(
			&transaction.ID,
			&transaction.Amount,
			&transaction.CreatedAt,
			&transaction.FirstName,
			&transaction.LastName,
			&transaction.Status,
			&transaction.Reason,
			&transaction.CategoryID,
			&transaction.Type)
		if err != nil {
			printError(w, err.Error())
		}
		transactions = append(transactions, transaction)
	}

	printJSON(w, &TransactionResponse{Transactions: transactions})

}
