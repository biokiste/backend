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
	ID            int     `json:"id"`
	Amount        float32 `json:"amount"`
	Type          string  `json:"type"`
	State         string  `json:"state"`
	UserID        int     `json:"userId"`
	CreatedAt     string  `json:"createdAt"`
	CreatedBy     int     `json:"createdBy"`
	UpdatedAt     string  `json:"updatedAt,omitempty"`
	UpdatedBy     int     `json:"updatedBy,omitempty"`
	UpdateComment string  `json:"updateComment,omitempty"`
}

func (h *Handlers) getTransactions(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	t := params.Get("type")
	s := params.Get("state")
	u := params.Get("user_id")
	c := params.Get("created_at")

	var str strings.Builder
	fmt.Fprint(&str, `SELECT ID, Amount, Type, State, UserID, CreatedAt, CreatedBy, COALESCE(UpdatedAt, '') AS UpdatedAt, COALESCE(UpdatedBy, 0) AS UpdatedBy, COALESCE(UpdateComment, '') AS UpdateComment FROM Transactions`)

	if t != "" || s != "" || u != "" || c != "" {
		fmt.Fprint(&str, " WHERE ")
	}

	if t != "" {
		fmt.Fprintf(&str, `Type = "%s"`, t)
		if s != "" || u != "" || c != "" {
			fmt.Fprint(&str, " AND ")
		}
	}

	if s != "" {
		fmt.Fprintf(&str, `State = "%s"`, s)
		if u != "" || c != "" {
			fmt.Fprint(&str, " AND ")
		}
	}

	if u != "" {
		fmt.Fprintf(&str, `UserID = %s`, u)
		if c != "" {
			fmt.Fprint(&str, " AND ")
		}
	}

	if c != "" {
		fmt.Fprintf(&str, `CreatedAt >= "%s" AND CreatedAt < DATE_ADD("%s", INTERVAL 1 DAY)`, c, c)
	}

	query := str.String()

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
			&t.Type,
			&t.State,
			&t.UserID,
			&t.CreatedAt,
			&t.CreatedBy,
			&t.UpdatedAt,
			&t.UpdatedBy,
			&t.UpdateComment,
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
		Amount        float32 `json:"amount"`
		Type          string  `json:"type"`
		State         string  `json:"state"`
		UserID        int     `json:"userId"`
		CreatedBy     int     `json:"createdBy"`
		CreatedAt     string  `json:"createdAt"`
		UpdateComment string  `json:"updateComment"`
	}

	var b body
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		fmt.Println(err)
		respondWithHTTP(w, http.StatusInternalServerError)
		return
	}

	if b.Amount == 0 && b.UpdateComment == "" || b.Type == "" || b.State == "" || b.UserID == 0 || b.CreatedBy == 0 || b.CreatedAt == "" {
		err := SimpleResponseBody{"Some required fields are missing!"}
		fmt.Println(err.Text)
		respondWithJSON(w, JSONResponse{http.StatusBadRequest, &err})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := h.DB.ExecContext(ctx,
		`INSERT INTO Transactions (Amount, Type, State, UserID, CreatedBy, CreatedAt)
		 VALUES (?,?,?,?,?,?)`,
		b.Amount,
		b.Type,
		b.State,
		b.UserID,
		b.CreatedBy,
		b.CreatedAt,
	)

	if err != nil {
		fmt.Println(err)
		respondWithHTTP(w, http.StatusInternalServerError)
		return
	}

	query := fmt.Sprintf(`UPDATE Users SET LastActivityAt = CURRENT_TIMESTAMP() WHERE ID = %d`, b.UserID)
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
