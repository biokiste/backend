package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
	// 	GET /transactions/?types={type ?string}&state={state ?string}&user_id={userId ?int}&createdAt={createdAt ?string}
	params := r.URL.Query()
	t := params.Get("type")

	var str strings.Builder
	fmt.Fprint(&str, `SELECT ID, Amount, Type, State, UserID, CreatedAt, CreatedBy, COALESCE(UpdatedAt, '') AS UpdatedAt, COALESCE(UpdatedBy, 0) AS UpdatedBy, COALESCE(UpdateComment, '') AS UpdateComment FROM Transactions`)

	if t != "" {
		fmt.Fprintf(&str, ` WHERE Type = "%s"`, t)
	}

	query := str.String()

	results, err := h.DB.Query(query)
	if err != nil {
		fmt.Println(err)
		printInternalError(w)
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
			printInternalError(w)
			return
		}
		transactions = append(transactions, t)
	}

	if len(transactions) == 0 {
		transactions = make([]transaction, 0)
	}

	printJSON(w, &transactions)

}

func (h *Handlers) addTransaction(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Amount    float32 `json:"amount"`
		Type      string  `json:"type"`
		State     string  `json:"state"`
		UserID    int     `json:"userId"`
		CreatedBy int     `json:"createdBy"`
	}

	var b body
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		fmt.Println(err)
		printInternalError(w)
		return
	}

	if b.Amount == 0 || b.Type == "" || b.State == "" || b.UserID == 0 || b.CreatedBy == 0 {
		code := 400
		msg := "Some required fields are missing!"
		fmt.Println(msg)
		printCustomError(w, ErrorMessage{code, msg}, code)
		return
	}

	result, err := h.DB.Exec(
		`INSERT INTO Transactions (Amount, Type, State, UserID, CreatedBy)
		 VALUES (?,?,?,?,?)`,
		b.Amount,
		b.Type,
		b.State,
		b.UserID,
		b.CreatedBy,
	)

	if err != nil {
		fmt.Println(err)
		printInternalError(w)
		return
	}

	type resBody struct {
		Status       string `json:"status"`
		LastInsertId int    `json:"id"`
	}

	id, err := result.LastInsertId()

	if err != nil {
		fmt.Println(err)
		printInternalError(w)
		return
	}

	printJSON(w, &resBody{"ok", int(id)})
}
