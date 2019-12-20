package main

import (
	"database/sql"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// UsersResponse JSON API Spec Wrapper
type UsersResponse struct {
	Users []User `json:"data"`
}

// User holds properties
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

// TransactionResponse JSON API Spec Wrapper
type TransactionResponse struct {
	Transactions []Transaction `json:"data"`
}

// Transaction holds properties
type Transaction struct {
	ID         int     `json:"id"`
	Amount     float32 `json:"amount"`
	CreatedAt  string  `json:"created_at"`
	FirstName  string  `json:"firstname"`
	LastName   string  `json:"lastname"`
	Status     int     `json:"status"`
	Reason     string  `json:"reason"`
	CategoryID int     `json:"category_id"`
	Type       string  `json:"type"`
}

// ShowStatus delivers actual status
func ShowStatus(w http.ResponseWriter, r *http.Request) {
	printJSON(w, "ok")
}

// ListUsers delivers user data
func ListUsers(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/foodkoop_biokiste")
	if err != nil {
		printError(w, err.Error())
	}
	defer db.Close()

	results, err := db.Query("SELECT id, username FROM users")
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

// GetPayments delivers all payments
func GetPayments(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/foodkoop_biokiste")
	if err != nil {
		printError(w, err.Error())
	}
	defer db.Close()

	results, err := db.Query(`
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

// GetPaymentsByUser delivers payments per user
func GetPaymentsByUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/foodkoop_biokiste")
	if err != nil {
		printError(w, err.Error())
	}
	defer db.Close()

	results, err := db.Query(`
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

	printJSON(w, &TransactionResponse{Transactions: transactions})
}
