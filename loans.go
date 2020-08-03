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
			h.updateLoanByID,
		},
	}
	return routes
}

type loan struct {
	ID        int     `json:"id"`
	Amount    float32 `json:"amount"`
	State     string  `json:"state"`
	UserID    int     `json:"userId"`
	CreatedAt string  `json:"createdAt"`
	CreatedBy int     `json:"createdBy"`
	UpdatedAt string  `json:"updatedAt,omitempty"`
	UpdatedBy int     `json:"updatedBy,omitempty"`
	Comment   string  `json:"comment,omitempty"`
}

func (h *Handlers) getLoans(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	state := params.Get("state")
	userID := params.Get("user_id")

	var str strings.Builder
	fmt.Fprint(&str, `SELECT ID, Amount, State, UserID, CreatedAt, CreatedBy, COALESCE(UpdatedAt, '') AS UpdatedAt, COALESCE(UpdatedBy, 0) AS UpdatedBy, COALESCE(Comment, '') AS Comment FROM Loans`)

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
		respondWithHTTP(w, http.StatusInternalServerError)
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
			&l.Comment,
		)
		if err != nil {
			fmt.Println(err)
			respondWithHTTP(w, http.StatusInternalServerError)
			return
		}
		loans = append(loans, l)
	}

	if len(loans) == 0 {
		loans = make([]loan, 0)
	}

	respondWithJSON(w, JSONResponse{Body: &loans})
}

func (h Handlers) updateLoanByID(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Amount    float32 `json:"amount,omitempty"`
		State     string  `json:"state,omitempty"`
		UserID    int     `json:"userId,omitempty"`
		UpdatedBy int     `json:"updatedBy"`
		Comment   string  `json:"comment"`
	}

	id, _ := mux.Vars(r)["id"]

	var b body
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		fmt.Println(err)
		respondWithHTTP(w, http.StatusInternalServerError)
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

	if b.UpdatedBy == 0 || b.Comment == "" {
		err := SimpleResponseBody{"Some required fields are missing!"}
		fmt.Println(err.Text)
		respondWithJSON(w, JSONResponse{http.StatusBadRequest, &err})
		return
	}

	fmt.Fprintf(&str, `UpdatedAt = CURRENT_TIMESTAMP(), UpdatedBy = %d, Comment = "%s" WHERE ID = %s`, b.UpdatedBy, b.Comment, id)

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

func (h Handlers) addLoan(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Amount    float32 `json:"amount"`
		State     string  `json:"state"`
		UserID    int     `json:"userId"`
		CreatedBy int     `json:"createdBy"`
		CreatedAt string  `json:"createdAt"`
		Comment   string  `json:"comment,omitempty"`
	}

	var b body
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		fmt.Println(err)
		respondWithHTTP(w, http.StatusInternalServerError)
		return
	}

	if b.Amount == 0 || b.State == "" || b.UserID == 0 || b.CreatedBy == 0 || b.CreatedAt == "" {
		err := SimpleResponseBody{"Some required fields are missing!"}
		fmt.Println(err.Text)
		respondWithJSON(w, JSONResponse{http.StatusBadRequest, &err})
		return
	}

	var insert strings.Builder
	var values strings.Builder

	fmt.Fprint(&insert, "INSERT INTO Loans (Amount, State, UserID, CreatedBy, CreatedAt")
	fmt.Fprintf(&values, `VALUES ("%f", "%s", "%d", "%d", "%s"`,
		b.Amount,
		b.State,
		b.UserID,
		b.CreatedBy,
		b.CreatedAt)

	if b.Comment != "" {
		fmt.Fprint(&insert, ", Comment")
		fmt.Fprintf(&values, ", %q", b.Comment)
	}

	fmt.Fprint(&insert, ")")
	fmt.Fprint(&values, ")")

	query := fmt.Sprintf("%s %s", insert.String(), values.String())

	result, err := h.DB.Exec(query)

	if err != nil {
		fmt.Println(err)
		respondWithHTTP(w, http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()

	if err != nil {
		fmt.Println(err)
		respondWithHTTP(w, http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, JSONResponse{Body: InsertResponseBody{int(id)}})
}
