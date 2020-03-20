package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// GetLoansRoutes get all routes of path /loans
func GetLoansRoutes(h *Handlers) []Route {
	routes := []Route{
		{
			"get loans",
			"GET",
			"/loans",
			h.getLoans,
		},
		{
			"add loan",
			"POST",
			"/loans",
			h.addLoan,
		},
		{
			"update loan",
			"PATCH",
			"/loans/{id}",
			h.updateLoanWithID,
		},
	}
	return routes
}

type loan struct {
	ID            int     `json:"id"`
	Amount        float32 `json:"amount"`
	State         string  `json:"state"`
	UserID        int     `json:"userId"`
	CreatedAt     string  `json:"createdAt"`
	CreatedBy     int     `json:"createdBy"`
	UpdatedAt     string  `json:"updatedAt,omitempty"`
	UpdatedBy     int     `json:"updatedBy,omitempty"`
	UpdateComment string  `json:"updateComment,omitempty"`
}

func (h *Handlers) getLoans(w http.ResponseWriter, r *http.Request) {
	results, err := h.DB.Query(`
		SELECT
			ID,
			Amount,
			State,
			UserID,
			CreatedAt,
			CreatedBy,
			COALESCE(UpdatedAt, '') AS UpdatedAt,
			COALESCE(UpdatedBy, -1) AS UpdatedBy,
			COALESCE(UpdateComment, '') AS UpdateComment
		FROM Loans
	`)
	if err != nil {
		fmt.Println(err)
		printDbError(w)
		return
	}
	defer results.Close()

	var loans []loan

	for results.Next() {
		var l loan
		err = results.Scan(
			&l.ID,
			&l.Amount,
			&l.State,
			&l.UserID,
			&l.CreatedAt,
			&l.CreatedBy,
			&l.UpdatedAt,
			&l.UpdatedBy,
			&l.UpdateComment,
		)
		if err != nil {
			fmt.Println(err)
			printDbError(w)
			return
		}
		loans = append(loans, l)
	}

	if len(loans) == 0 {
		loans = make([]loan, 0)
	}

	printJSON(w, &loans)
}

func (h Handlers) updateLoanWithID(w http.ResponseWriter, r *http.Request) {
	// TODO: Make Amount, State and UserID optional
	type body struct {
		Amount        float32 `json:"amount"`
		State         string  `json:"state"`
		UserID        int     `json:"userId"`
		UpdatedBy     int     `json:"updatedBy"`
		UpdateComment string  `json:"updateComment"`
	}

	id, _ := mux.Vars(r)["id"]

	var b body
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		printDbError(w)
		return
	}

	// TODO: Validate data – e.g. if all fields set

	date := time.Now()

	result, err := h.DB.Exec(
		`UPDATE Loans
			SET Amount = ?, State = ?, UserID = ?, UpdatedAt = ?, UpdatedBy = ?, UpdateComment = ?
		 WHERE ID = ?`,
		b.Amount,
		b.State,
		b.UserID,
		date,
		b.UpdatedBy,
		b.UpdateComment,
		id,
	)

	if err != nil {
		fmt.Println(err)
		printDbError(w)
		return
	}

	rowsAffected, _ := result.RowsAffected()

	type resBody struct {
		Status       string `json:"status"`
		RowsAffected int    `json:"rowsAffected"`
	}

	printJSON(w, &resBody{"ok", int(rowsAffected)})
}

func (h Handlers) addLoan(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Amount    float32 `json:"amount"`
		State     string  `json:"state"`
		UserID    int     `json:"userId"`
		CreatedBy int     `json:"createdBy"`
	}

	var b body
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		fmt.Println(err)
		printDbError(w)
		return
	}

	// TODO: Validate data – e.g. if all fields set

	result, err := h.DB.Exec(
		`INSERT INTO Loans (Amount, State, UserID, CreatedBy)
		 VALUES (?,?,?,?)`,
		b.Amount,
		b.State,
		b.UserID,
		b.CreatedBy,
	)

	if err != nil {
		fmt.Println(err)
		printDbError(w)
		return
	}

	type resBody struct {
		Status       string `json:"status"`
		LastInsertId int    `json:"id"`
	}

	id, err := result.LastInsertId()

	if err != nil {
		printDbError(w)
		return
	}

	printJSON(w, &resBody{"ok", int(id)})
}
