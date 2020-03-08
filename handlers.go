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
		printDbError(w)
		return
	}

	printJSON(w, &DoorCodeResponse{DoorCode: doorCode})
}

// UpdateDoorCode sets doorcode
func (h Handlers) UpdateDoorCode(w http.ResponseWriter, r *http.Request) {
	var updatedDoorCode DoorCode
	err := json.NewDecoder(r.Body).Decode(&updatedDoorCode)
	if err != nil {
		printDbError(w)
		return
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
	users, err := h.GetAllUser()
	if err != nil {
		printDbError(w)
	} else {
		printJSON(w, &UsersResponse{Users: users})
	}
}

// LastActiveUsers delivers last ten active users
func (h Handlers) LastActiveUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.GetLastActiveUsers()
	if err != nil {
		printDbError(w)
	} else {
		printJSON(w, &UsersResponse{Users: users})
	}
}

// GetUser delivers user for id
// func (h Handlers) GetUser(w http.ResponseWriter, r *http.Request) {
// 	id, _ := strconv.Atoi(mux.Vars(r)["id"])
// 	user, err := h.GetSingleUser(id)
// 	if err != nil {
// 		printDbError(w)
// 	} else {
// 		printJSON(w, &UserResponse{User: user})
// 	}
// }

// GetUserByEmail delivers user for email
func (h Handlers) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	email, _ := mux.Vars(r)["email"]
	user, err := h.GetSingleUserByIEmail(email)
	if err != nil {
		printDbError(w)
	} else {
		printJSON(w, &UserResponse{User: user})
	}
}

// GetAuth0User delivers auth0 user data
func (h Handlers) GetAuth0User(w http.ResponseWriter, r *http.Request) {
	id, _ := (mux.Vars(r)["id"])
	var user Auth0User
	user, err := h.Auth0GetUser(id)
	if err != nil {
		printError(w, err)
	} else {
		printJSON(w, user)
	}
}

// CreateUser creates Auth0 user and user in app database
func (h Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		printDbError(w)
		return
	}

	// first create auth0 user
	auth0User := Auth0User{
		UserID:     user.UserID,
		Password:   user.Password,
		Email:      user.Email,
		Connection: "Username-Password-Authentication",
	}
	statusCode := h.CreateAuth0User(auth0User)

	if statusCode != 201 {
		printCustomError(w, err, statusCode)
		return
	}

	// then create user in app db
	id, err := h.CreateUserData(user)
	if err != nil {
		printError(w, err)
		return
	}

	printJSON(w, &User{
		ID: int(id),
	})
}

// UpdateUser updates user
func (h Handlers) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		printError(w, err)
		return
	}

	err = h.UpdateUserData(user)
	if err != nil {
		printError(w, err.Error)
		return
	}
	printSuccess(w)
}

// GetTransactions delivers all payments
func (h Handlers) GetTransactions(w http.ResponseWriter, r *http.Request) {

	results, err := h.DB.Query(`
		SELECT transactions.id, amount, transactions.created_at, firstname, lastname, transactions.status, transactions.reason, category_id, transactions_category.type
		FROM transactions
		LEFT JOIN transactions_category ON transactions.category_id = transactions_category.id
		LEFT JOIN users ON transactions.user_id = users.id
		WHERE NOT transactions_category.type = "correction"
		AND firstname IS NOT NULL AND transactions.status = 1 		
		ORDER BY transactions.created_at desc
	  `)
	if err != nil {
		printDbError(w)
		return
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
			printError(w, err)
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
		printDbError(w)
		return
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
			printError(w, err)
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
		printDbError(w)
		return
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
			printDbError(w)
		}
		transactionCategories = append(transactionCategories, transactionCategory)
	}
	printJSON(w, &TransactionCategoryResponse{TransactionCategories: transactionCategories})
}

// GetTransactionStates delivers transaction states
func (h Handlers) GetTransactionStates(w http.ResponseWriter, r *http.Request) {
	results, err := h.DB.Query(`
		SELECT id, type
		FROM transactions_status
	`)
	if err != nil {
		printDbError(w)
		return
	}
	defer results.Close()

	var transactionStates []TransactionState
	for results.Next() {
		var transactionState TransactionState
		err = results.Scan(
			&transactionState.ID,
			&transactionState.Type)
		if err != nil {
			printDbError(w)
		}
		transactionStates = append(transactionStates, transactionState)
	}
	printJSON(w, &TransactionStateResponse{TransactionStates: transactionStates})
}

// AddTransaction updates user balance with transactions
func (h Handlers) AddTransaction(w http.ResponseWriter, r *http.Request) {
	var transactionRequest TransactionRequest
	err := json.NewDecoder(r.Body).Decode(&transactionRequest)
	if err != nil {
		printDbError(w)
		return
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
		printDbError(w)
		return
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
		printDbError(w)
		return
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
			printDbError(w)
		}
		transactions = append(transactions, transaction)
	}

	printJSON(w, &TransactionResponse{Transactions: transactions})

}

// GetGroupTypes returns all type of groups
func (h Handlers) GetGroupTypes(w http.ResponseWriter, r *http.Request) {
	results, err := h.DB.Query(`
		SELECT id, name, description
		FROM groups
	`)
	if err != nil {
		printDbError(w)
		return
	}
	defer results.Close()

	var groupTypes []GroupType
	for results.Next() {
		var groupType GroupType
		err = results.Scan(
			&groupType.ID,
			&groupType.Name,
			&groupType.Description)
		if err != nil {
			printDbError(w)
		}
		groupTypes = append(groupTypes, groupType)
	}
	printJSON(w, &GroupTypesRequest{GroupTypes: groupTypes})
}

// GetGroups returns all groups
func (h Handlers) GetGroups(w http.ResponseWriter, r *http.Request) {
	groups, err := h.GetGroupsWithUsers()
	if err != nil {
		printDbError(w)
		return
	}
	printJSON(w, &GroupRequest{Groups: groups})
}

// SendMail sends emails
// func (h Handlers) SendMail(w http.ResponseWriter, r *http.Request) {
// 	mailRecipient := "sebastian.koslitz@gmail.com"
// 	err := h.SendEMail(mailRecipient)

// 	if err != nil {
// 		printError(w, err)
// 	}
// 	printSuccess(w)
// }
