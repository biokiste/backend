package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// GetGroupsRoutes get all routes of path /groups
func GetGroupsRoutes(h *Handlers) []Route {
	routes := []Route{
		{
			"add group",
			"POST",
			"/groups",
			h.addGroup,
		},
		{
			"update group",
			"PATCH",
			"/groups/{id}",
			h.updateGroupByID,
		},
		{
			"get groups",
			"GET",
			"/groups",
			h.getGroups,
		},
		{
			"get grogroupups",
			"GET",
			"/groups/{id}",
			h.getGroupByID,
		},
	}

	return routes
}

type group struct {
	ID            int    `json:"id"`
	Key           string `json:"key"`
	Email         string `json:"email"`
	CreatedAt     string `json:"createdAt"`
	CreatedBy     int    `json:"createdBy"`
	UpdatedAt     string `json:"updatedAt,omitempty"`
	UpdatedBy     int    `json:"updatedBy,omitempty"`
	UpdateComment string `json:"updateComment,omitempty"`
	Users         []int  `json:"users"`
	Leaders       []int  `json:"leaders"`
}

type groupUser struct {
	ID       int
	GroupID  int
	UserID   int
	IsLeader int
}

func (h *Handlers) addGroup(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Key       string `json:"key"`
		Email     string `json:"email"`
		CreatedBy int    `json:"createdBy"`
	}

	var b body
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		fmt.Println(err)
		respondWithHTTP(w, http.StatusInternalServerError)
		return
	}

	if b.Key == "" || b.Email == "" || b.CreatedBy == 0 {
		text := "Some required fields are missing!"
		fmt.Println(text)
		respondWithJSON(w, JSONResponse{http.StatusBadRequest, SimpleResponseBody{Text: text}})
		return
	}

	result, err := h.DB.Exec(
		`INSERT INTO Groups (GroupKey, Email, CreatedBy)
		 VALUES (?,?,?)`,
		b.Key,
		b.Email,
		b.CreatedBy,
	)

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

func (h *Handlers) updateGroupByID(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Key           string `json:"key,omitempty"`
		Email         string `json:"amount,omitempty"`
		UpdatedBy     int    `json:"updatedBy"`
		UpdateComment string `json:"updateComment"`
	}

	id, _ := mux.Vars(r)["id"]

	var b body
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		fmt.Println(err)
		respondWithHTTP(w, http.StatusInternalServerError)
		return
	}

	if b.UpdatedBy == 0 || b.UpdateComment == "" {
		err := SimpleResponseBody{"Some required fields are missing!"}
		fmt.Println(err.Text)
		respondWithJSON(w, JSONResponse{http.StatusBadRequest, &err})
		return
	}

	var str strings.Builder

	fmt.Fprint(&str, "UPDATE Groups SET ")
	if b.Key != "" {
		fmt.Fprintf(&str, `GroupKey = "%s", `, b.Key)
	}
	if b.Email != "" {
		fmt.Fprintf(&str, "Email = %s, ", b.Email)
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

	resObj := UpdateResponseBody{int(rowsAffected)}
	respondWithJSON(w, JSONResponse{Body: &resObj})
}

func (h *Handlers) getGroups(w http.ResponseWriter, r *http.Request) {
	query := "SELECT ID, GroupKey, Email, CreatedAt, CreatedBy, COALESCE(UpdatedAt, '') AS UpdatedAt, COALESCE(UpdatedBy, 0) AS UpdatedBy, COALESCE(UpdateComment, '') AS UpdateComment FROM Groups"
	results, err := h.DB.Query(query)

	if err != nil {
		fmt.Println(err)
		respondWithHTTP(w, http.StatusInternalServerError)
		return
	}
	defer results.Close()

	var groups []group
	for results.Next() {
		var g group
		err := results.Scan(
			&g.ID,
			&g.Key,
			&g.Email,
			&g.CreatedAt,
			&g.CreatedBy,
			&g.UpdatedAt,
			&g.UpdatedBy,
			&g.UpdateComment,
		)
		if err != nil {
			fmt.Println(err)
			respondWithHTTP(w, http.StatusInternalServerError)
			return
		}

		q := fmt.Sprintf("SELECT ID, GroupID, UserID, IsLeader FROM GroupUsers WHERE GroupID = %d", g.ID)
		result, err := h.DB.Query(q)
		if err != nil {
			fmt.Println(err)
			respondWithHTTP(w, http.StatusInternalServerError)
			return
		}
		defer result.Close()
		for result.Next() {
			var gu groupUser
			err := result.Scan(
				&gu.ID,
				&gu.GroupID,
				&gu.UserID,
				&gu.IsLeader,
			)
			if err != nil {
				fmt.Println(err)
				respondWithHTTP(w, http.StatusInternalServerError)
				return
			}
			g.Users = append(g.Users, gu.UserID)
			if gu.IsLeader == 1 {
				g.Leaders = append(g.Leaders, gu.UserID)
			}
		}
		if len(g.Users) == 0 {
			g.Users = make([]int, 0)
		}
		if len(g.Leaders) == 0 {
			g.Leaders = make([]int, 0)
		}
		groups = append(groups, g)
	}

	if len(groups) == 0 {
		groups = make([]group, 0)
	}
	respondWithJSON(w, JSONResponse{Body: &groups})
}

func (h *Handlers) getGroupByID(w http.ResponseWriter, r *http.Request) {
	id, _ := mux.Vars(r)["id"]

	query := fmt.Sprintf("SELECT ID, GroupKey, Email, CreatedAt, CreatedBy, COALESCE(UpdatedAt, '') AS UpdatedAt, COALESCE(UpdatedBy, 0) AS UpdatedBy, COALESCE(UpdateComment, '') AS UpdateComment FROM Groups WHERE ID = %s", id)

	row := h.DB.QueryRow(query)

	var g group

	row.Scan(
		&g.ID,
		&g.Key,
		&g.Email,
		&g.CreatedAt,
		&g.CreatedBy,
		&g.UpdatedAt,
		&g.UpdatedBy,
		&g.UpdateComment,
	)

	if g.ID == 0 {
		respondWithHTTP(w, http.StatusNotFound)
		return
	}

	q := fmt.Sprintf("SELECT ID, GroupID, UserID, IsLeader FROM GroupUsers WHERE GroupID = %d", g.ID)

	result, err := h.DB.Query(q)
	if err != nil {
		fmt.Println(err)
		respondWithHTTP(w, http.StatusInternalServerError)
		return
	}
	defer result.Close()
	for result.Next() {
		var gu groupUser
		err := result.Scan(
			&gu.ID,
			&gu.GroupID,
			&gu.UserID,
			&gu.IsLeader,
		)
		if err != nil {
			fmt.Println(err)
			respondWithHTTP(w, http.StatusInternalServerError)
			return
		}
		g.Users = append(g.Users, gu.UserID)
		if gu.IsLeader == 1 {
			g.Leaders = append(g.Leaders, gu.UserID)
		}
	}
	if len(g.Users) == 0 {
		g.Users = make([]int, 0)
	}
	if len(g.Leaders) == 0 {
		g.Leaders = make([]int, 0)
	}

	respondWithJSON(w, JSONResponse{Body: &g})
}
