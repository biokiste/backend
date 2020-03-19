package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/didi/gendry/scanner"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
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

	err = scanner.Scan(results, &transactions)
	if err != nil {
		printError(w, err)
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
	err = scanner.Scan(results, &transactions)
	if err != nil {
		printError(w, err)
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
			printDbError(w)
			return
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
			printDbError(w)
			return
		}
	}

	err = h.LogUserTransaction(transactionRequest.User.ID)
	if err != nil {
		printDbError(w)
		return
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
	err = scanner.Scan(results, &transactions)
	if err != nil {
		printDbError(w)
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
	err = scanner.Scan(results, &groupTypes)
	if err != nil {
		printDbError(w)
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

type setting struct {
	ID            int    `json:"id"`
	ItemKey       string `json:"key"`
	ItemValue     string `json:"value"`
	CreatedAt     string `json:"createdAt"`
	CreatedBy     int    `json:"createdBy"`
	UpdatedAt     string `json:"updatedAt,omitempty"`
	UpdatedBy     int    `json:"updatedBy,omitempty"`
	UpdateComment string `json:"updateComment,omitempty"`
}

// GetSettings get all setting
func (h Handlers) GetSettings(w http.ResponseWriter, r *http.Request) {
	results, err := h.DB.Query(`
		SELECT 
			ID,
			ItemKey,
			ItemValue,
			CreatedAt,
			CreatedBy,
			COALESCE(UpdatedAt, '') AS UpdatedAt,
			COALESCE(UpdatedBy, -1) AS UpdatedBy,
			COALESCE(UpdateComment, '') AS UpdateComment
		FROM Settings
	`)
	if err != nil {
		printDbError(w)
		return
	}
	defer results.Close()

	var settings []setting

	for results.Next() {
		var s setting
		err = results.Scan(
			&s.ID,
			&s.ItemKey,
			&s.ItemValue,
			&s.CreatedAt,
			&s.CreatedBy,
			&s.UpdatedAt,
			&s.UpdatedBy,
			&s.UpdateComment,
		)
		if err != nil {
			printDbError(w)
			return
		}
		settings = append(settings, s)
	}

	if len(settings) == 0 {
		settings = make([]setting, 0)
	}

	printJSON(w, &settings)
}

// GetSettingByKey get setting by key
func (h Handlers) GetSettingByKey(w http.ResponseWriter, r *http.Request) {
	key, _ := mux.Vars(r)["key"]
	var s setting
	if err := h.DB.QueryRow(`
		SELECT 
			ID,
			ItemKey,
			ItemValue,
			CreatedAt,
			CreatedBy,
			COALESCE(UpdatedAt, '') AS UpdatedAt,
			COALESCE(UpdatedBy, -1) AS UpdatedBy,
			COALESCE(UpdateComment, '') AS UpdateComment
		FROM Settings
		WHERE ItemKey = ?`, key).Scan(
		&s.ID,
		&s.ItemKey,
		&s.ItemValue,
		&s.CreatedAt,
		&s.CreatedBy,
		&s.UpdatedAt,
		&s.UpdatedBy,
		&s.UpdateComment,
	); err != nil {
		fmt.Println(err)
		printCustomError(w, nil, 404)
		return
	}
	printJSON(w, &s)
}

// UpdateSettingByKey store new setting
func (h Handlers) UpdateSettingByKey(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Value         string `json:"value"`
		UpdatedBy     int    `json:"updatedBy"`
		UpdateComment string `json:"updateComment"`
	}

	key, _ := mux.Vars(r)["key"]

	var b body
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		printDbError(w)
		return
	}

	// TODO: Validate data – e.g. if all fields set

	t := time.Now()

	res, err := h.DB.Exec(
		`UPDATE Settings
			SET ItemValue = ?, UpdatedAt = ?, UpdatedBy = ?, UpdateComment = ?
		 WHERE ItemKey = ?`,
		b.Value,
		t,
		b.UpdatedBy,
		b.UpdateComment,
		key,
	)

	if err != nil {
		fmt.Println(err)
		printDbError(w)
		return
	}

	rowsAffected, _ := res.RowsAffected()

	type responseBody struct {
		Status       string `json:"status"`
		RowsAffected int    `json:"rowsAffected"`
	}

	printJSON(w, &responseBody{"ok", int(rowsAffected)})
}

// AddSetting store new setting
func (h Handlers) AddSetting(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Key       string `json:"key"`
		Value     string `json:"value"`
		CreatedBy int    `json:"createdBy"`
	}

	var b body
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		printDbError(w)
		return
	}

	// TODO: Validate data – e.g. if all fields set

	res, err := h.DB.Exec(
		`INSERT INTO Settings (ItemKey, ItemValue, CreatedBy)
		 VALUES (?,?,?)`,
		b.Key,
		b.Value,
		b.CreatedBy,
	)

	if err != nil {
		printDbError(w)
		return
	}

	type responseBody struct {
		Status string `json:"status"`
		ID     int    `json:"id"`
	}

	id, err := res.LastInsertId()

	if err != nil {
		printDbError(w)
		return
	}

	printJSON(w, &responseBody{"ok", int(id)})
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
