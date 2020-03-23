package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	}

	return routes
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
