package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// GetTransactionsRoutes get all routes of path /transactions
func GetTransactionsRoutes(h *Handlers) []Route {
	routes := []Route{
		{
			"get transactions",
			"GET",
			"/transactions",
			h.getTransactions,
		},
		{
			"add transaction",
			"POST",
			"/transactions",
			h.addTransaction,
		},
		{
			"update transaction",
			"PATCH",
			"/transactions/{id}",
			h.updateTransactionByID,
		},
	}

	return routes
}

type transaction struct {
	ID         int       `json:"id" db:"id"`
	UserID     int       `json:"user_id" db:"user_id"`
	CategoryID int       `json:"category_id" db:"category_id"`
	Amount     float32   `json:"amount" db:"amount"`
	StatusID   int       `json:"status_id" db:"status_id"`
	UpdatedBy  int       `json:"updated_by" db:"updated_by"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	Reason     string    `json:"reason" db:"reason"`
}

func (h *Handlers) getTransactions(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	t := params.Get("type")
	s := params.Get("state")
	u := params.Get("user_id")
	c := params.Get("created_at")

	var str strings.Builder
	fmt.Fprint(&str, `SELECT 
	transactions.id, transactions.amount, transactions.type,
	transactions.state, transactions.user_id, transactions.created_at,
	transactions.created_by, 
	COALESCE(transactions.updated_at, '') AS updated_at,
	COALESCE(transactions.updated_by, 0) AS updated_by, 
	COALESCE(transactions.update_comment, '') AS update_comment, 
	users.firstname, users.lastname FROM transactions
	INNER JOIN users ON transactions.user_id=users.id`)

	if t != "" || s != "" || u != "" || c != "" {
		fmt.Fprint(&str, " WHERE ")
	}

	if t != "" {
		fmt.Fprintf(&str, `transactions.type = "%s"`, t)
		if s != "" || u != "" || c != "" {
			fmt.Fprint(&str, " AND ")
		}
	}

	if s != "" {
		fmt.Fprintf(&str, `transactions.state = "%s"`, s)
		if u != "" || c != "" {
			fmt.Fprint(&str, " AND ")
		}
	}

	if u != "" {
		fmt.Fprintf(&str, `transactions.user_id = %s`, u)
		if c != "" {
			fmt.Fprint(&str, " AND ")
		}
	}

	if c != "" {
		fmt.Fprintf(&str, `transactions.created_at >= "%s" AND transactions.created_at < DATE_ADD("%s", INTERVAL 1 DAY)`, c, c)
	}

	query := str.String()

	fmt.Println(query)

	results, err := h.DB.Query(query)
	if err != nil {
		fmt.Println(err)
		respondWithHTTP(w, http.StatusInternalServerError)
		return
	}
	defer results.Close()

	var transactions []transaction

	for results.Next() {
		var t transaction
		err = results.Scan(
			&t.ID,
			&t.Amount,
			&t.UserID,
			&t.CreatedAt,
			&t.UpdatedAt,
			&t.UpdatedBy,
		)
		if err != nil {
			fmt.Println(err)
			respondWithHTTP(w, http.StatusInternalServerError)
			return
		}
		transactions = append(transactions, t)
	}

	if len(transactions) == 0 {
		transactions = make([]transaction, 0)
	}

	respondWithJSON(w, JSONResponse{Body: &transactions})
}

func (h *Handlers) updateTransactionByID(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Amount        float32 `json:"amount,omitempty"`
		Type          string  `json:"type,omitempty"`
		State         string  `json:"state,omitempty"`
		UserID        int     `json:"userId,omitempty"`
		UpdatedBy     int     `json:"updatedBy"`
		UpdateComment string  `json:"updateComment"`
	}

	var b body
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		fmt.Println(err)
		respondWithHTTP(w, http.StatusInternalServerError)
		return
	}

	id, _ := mux.Vars(r)["id"]

	var str strings.Builder

	fmt.Fprint(&str, "UPDATE Transactions SET ")
	if b.Amount != 0 {
		fmt.Fprintf(&str, "Amount = %f, ", b.Amount)
	}
	if b.Type != "" {
		fmt.Fprintf(&str, `Type = "%s", `, b.Type)
	}
	if b.State != "" {
		fmt.Fprintf(&str, `State = "%s", `, b.State)
	}
	if b.UserID != 0 {
		fmt.Fprintf(&str, "UserID = %d, ", b.UserID)
	}

	if b.UpdatedBy == 0 || b.UpdateComment == "" {
		err := SimpleResponseBody{"Some required fields are missing!"}
		fmt.Println(err.Text)
		respondWithJSON(w, JSONResponse{http.StatusBadRequest, &err})
		return
	}

	fmt.Fprintf(&str, `UpdatedAt = CURRENT_TIMESTAMP(), UpdatedBy = %d, UpdateComment = "%s" WHERE ID = %s`, b.UpdatedBy, b.UpdateComment, id)

	query := str.String()

	result, err := h.DB.Exec(query)

	if err != nil {
		fmt.Println(err)
		respondWithHTTP(w, http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()

	respondWithJSON(w, JSONResponse{Body: UpdateResponseBody{int(rowsAffected)}})
}

func (h *Handlers) addTransaction(w http.ResponseWriter, r *http.Request) {
	type body struct {
		ID         int     `json:"id" db:"id"`
		UserID     int     `json:"user_id" db:"user_id"`
		CategoryID int     `json:"category_id" db:"category_id"`
		Amount     float32 `json:"amount" db:"amount"`
		StatusID   int     `json:"status_id" db:"status_id"`
		UpdatedBy  int     `json:"updated_by" db:"updated_by"`
		CreatedAt  string  `json:"created_at" db:"created_at"`
		UpdatedAt  string  `json:"updated_at" db:"updated_at"`
		Reason     string  `json:"reason" db:"reason"`
	}

	var b body
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		fmt.Println(err)
		respondWithHTTP(w, http.StatusInternalServerError)
		return
	}

	// if b.Amount == 0 && b.Reason == "" || b.CategoryID == "" || b.StatusID == "" || b.UserID == 0 || b.CreatedBy == 0 || b.CreatedAt == "" {
	// 	err := SimpleResponseBody{"Some required fields are missing!"}
	// 	fmt.Println(err.Text)
	// 	respondWithJSON(w, JSONResponse{http.StatusBadRequest, &err})
	// 	return
	// }

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var _updatedAt string
	if b.UpdatedAt == "0000-00-00 00:00:00" {
		_updatedAt = b.CreatedAt
	} else {
		_updatedAt = b.UpdatedAt
	}

	result, err := h.DB.ExecContext(ctx,
		`INSERT INTO transactions (amount, category_id, status_id, user_id, created_at, updated_at, reason)
			 VALUES (?,?,?,?,?,?,?)`,
		b.Amount,
		b.CategoryID,
		b.StatusID,
		b.UserID,
		b.CreatedAt,
		_updatedAt,
		b.Reason,
	)

	if err != nil {
		fmt.Println(err)
		respondWithHTTP(w, http.StatusInternalServerError)
		return
	}

	query := fmt.Sprintf(`UPDATE users SET last_activity_at = CURRENT_TIMESTAMP() WHERE id = %d`, b.UserID)
	_, err = h.DB.Exec(query)
	if err != nil {
		fmt.Println(err)
		respondWithHTTP(w, http.StatusInternalServerError)
		return
	}

	type resBody struct {
		Status       string `json:"status"`
		LastInsertId int    `json:"id"`
	}

	id, err := result.LastInsertId()

	if err != nil {
		fmt.Println(err)
		respondWithHTTP(w, http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, JSONResponse{Body: InsertResponseBody{int(id)}})
}
