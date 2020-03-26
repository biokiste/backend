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
		printInternalError(w)
		return
	}

	if b.Key == "" || b.Email == "" || b.CreatedBy == 0 {
		code := 400
		msg := "Some required fields are missing!"
		fmt.Println(msg)
		printCustomError(w, ErrorMessage{code, msg}, code)
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
		printInternalError(w)
		return
	}

	if b.UpdatedBy == 0 || b.UpdateComment == "" {
		code := 400
		msg := "Some required fields are missing!"
		fmt.Println(msg)
		printCustomError(w, ErrorMessage{code, msg}, code)
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

func (h *Handlers) getGroups(w http.ResponseWriter, r *http.Request) {
	query := "SELECT ID, GroupKey, Email, CreatedAt, CreatedBy, COALESCE(UpdatedAt, '') AS UpdatedAt, COALESCE(UpdatedBy, 0) AS UpdatedBy, COALESCE(UpdateComment, '') AS UpdateComment FROM Groups"
	results, err := h.DB.Query(query)

	if err != nil {
		fmt.Println(err)
		printInternalError(w)
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
			printInternalError(w)
			return
		}

		q := fmt.Sprintf("SELECT ID, GroupID, UserID, IsLeader FROM GroupUsers WHERE GroupID = %d", g.ID)
		r, err := h.DB.Query(q)
		if err != nil {
			fmt.Println(err)
			printInternalError(w)
			return
		}
		defer r.Close()
		for r.Next() {
			var gu groupUser
			err := r.Scan(
				&gu.ID,
				&gu.GroupID,
				&gu.UserID,
				&gu.IsLeader,
			)
			if err != nil {
				fmt.Println(err)
				printInternalError(w)
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

	printJSON(w, &groups)
}

// // GetGroupsWithUsers return groups with user ids and ids of group leaders
// func (h Handlers) GetGroupsWithUsers() ([]Group, error) {
// 	results, err := h.DB.Query(`
// 		SELECT
// 			group_id, user_id, position_id
// 		FROM
// 			groups_users
// 		WHERE
// 			active=1
// 	`)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer results.Close()

// 	var groups []Group

// 	for results.Next() {
// 		var entry GroupUserEntry
// 		err = results.Scan(&entry.GroupID, &entry.UserID, &entry.PositionID)

// 		if err != nil {
// 			panic(err.Error())
// 		}

// 		var idx int = -1

// 		for i, g := range groups {
// 			if g.ID == entry.GroupID {
// 				idx = i
// 				break
// 			}
// 		}

// 		if idx == -1 {
// 			var newGroup Group
// 			newGroup.ID = entry.GroupID
// 			newGroup.UserIDs = append(newGroup.UserIDs, entry.UserID)
// 			if entry.PositionID == 1 {
// 				newGroup.LeaderIDs = append(newGroup.LeaderIDs, entry.UserID)
// 			}
// 			groups = append(groups, newGroup)
// 		} else {
// 			var group = groups[idx]
// 			group.UserIDs = append(group.UserIDs, entry.UserID)
// 			if entry.PositionID == 1 {
// 				group.LeaderIDs = append(group.LeaderIDs, entry.UserID)
// 			}
// 			groups[idx] = group
// 		}

// 	}

// 	return groups, err
// }
