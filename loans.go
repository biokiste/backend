package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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
	params := r.URL.Query()
	state := params.Get("state")
	userID := params.Get("user_id")

	var str strings.Builder
	fmt.Fprint(&str, `SELECT ID, Amount, State, UserID, CreatedAt, CreatedBy, COALESCE(UpdatedAt, '') AS UpdatedAt, COALESCE(UpdatedBy, 0) AS UpdatedBy, COALESCE(UpdateComment, '') AS UpdateComment FROM Loans`)

	if state != "" || userID != "" {
		fmt.Fprint(&str, " WHERE ")
	}

	if state != "" && userID != "" {
		fmt.Fprintf(&str, `State = "%s" AND UserID = %s`, state, userID)
	} else if state != "" {
		fmt.Fprintf(&str, `State = "%s"`, state)
	} else if userID != "" {
		fmt.Fprintf(&str, `UserID = %s`, userID)
	}

	query := str.String()

	results, err := h.DB.Query(query)
	if err != nil {
		fmt.Println(err)
		printInternalError(w)
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
			printInternalError(w)
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
	type body struct {
		Amount        float32 `json:"amount,omitempty"`
		State         string  `json:"state,omitempty"`
		UserID        int     `json:"userId,omitempty"`
		UpdatedBy     int     `json:"updatedBy"`
		UpdateComment string  `json:"updateComment"`
	}

	id, _ := mux.Vars(r)["id"]

	var b body
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		fmt.Println(err)
		printInternalError(w)
		return
	}

	var str strings.Builder

	fmt.Fprint(&str, "UPDATE Loans SET ")
	if b.Amount != 0 {
		fmt.Fprintf(&str, "Amount = %f, ", b.Amount)
	}
	if b.State != "" {
		fmt.Fprintf(&str, `State = "%s", `, b.State)
	}
	if b.UserID != 0 {
		fmt.Fprintf(&str, "UserID = %d, ", b.UserID)
	}

	if b.UpdatedBy == 0 || b.UpdateComment == "" {
		code := 400
		msg := "Some required fields are missing!"
		fmt.Println(msg)
		printCustomError(w, ErrorMessage{code, msg}, code)
		return
	}

	fmt.Fprintf(&str, `UpdatedAt = CURRENT_TIMESTAMP(), UpdatedBy = %d, UpdateComment = "%s" WHERE ID = %s`, b.UpdatedBy, b.UpdateComment, id)

	query := str.String()

	result, err := h.DB.Exec(query)

	if err != nil {
		fmt.Println(err)
		printInternalError(w)
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
		printInternalError(w)
		return
	}

	if b.Amount == 0 || b.State == "" || b.UserID == 0 || b.CreatedBy == 0 {
		code := 400
		msg := "Some required fields are missing!"
		fmt.Println(msg)
		printCustomError(w, ErrorMessage{code, msg}, code)
		return
	}

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
