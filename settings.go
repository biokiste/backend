package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// GetSettingsRoutes get all routes of path /settings
func GetSettingsRoutes(h *Handlers) []Route {
	routes := []Route{
		{
			"get settings",
			"GET",
			"/settings",
			h.getSettings,
		},
		{
			"get setting by key",
			"GET",
			"/settings/{key}",
			h.getSettingByKey,
		},
		{
			"add setting",
			"POST",
			"/settings",
			h.addSetting,
		},
		{
			"update setting",
			"PATCH",
			"/settings/{key}",
			h.updateSettingWithKey,
		},
	}
	return routes
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

func (h Handlers) getSettings(w http.ResponseWriter, r *http.Request) {
	results, err := h.DB.Query(`
		SELECT 
			ID,
			ItemKey,
			ItemValue,
			CreatedAt,
			CreatedBy,
			COALESCE(UpdatedAt, '') AS UpdatedAt,
			COALESCE(UpdatedBy, 0) AS UpdatedBy,
			COALESCE(UpdateComment, '') AS UpdateComment
		FROM Settings
	`)
	if err != nil {
		fmt.Println(err)
		printInternalError(w)
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
			fmt.Println(err)
			printInternalError(w)
			return
		}
		settings = append(settings, s)
	}

	if len(settings) == 0 {
		settings = make([]setting, 0)
	}

	printJSON(w, &settings)
}

func (h Handlers) getSettingByKey(w http.ResponseWriter, r *http.Request) {
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
			COALESCE(UpdatedBy, 0) AS UpdatedBy,
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
		printInternalError(w)
		return
	}
	printJSON(w, &s)
}

func (h Handlers) updateSettingWithKey(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Value         string `json:"value"`
		UpdatedBy     int    `json:"updatedBy"`
		UpdateComment string `json:"updateComment"`
	}

	key, _ := mux.Vars(r)["key"]

	var b body
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		fmt.Println(err)
		printInternalError(w)
		return
	}

	if b.Value == "" || b.UpdatedBy == 0 || b.UpdateComment == "" {
		code := 400
		msg := "Some required fields are missing!"
		fmt.Println(msg)
		printCustomError(w, ErrorMessage{code, msg}, code)
		return
	}

	result, err := h.DB.Exec(`
		UPDATE Settings
			SET ItemValue = ?, UpdatedAt = CURRENT_TIMESTAMP(), UpdatedBy = ?, UpdateComment = ?
		WHERE ItemKey = ?
	`, b.Value, b.UpdatedBy, b.UpdateComment, key)

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

func (h Handlers) addSetting(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Key       string `json:"key"`
		Value     string `json:"value"`
		CreatedBy int    `json:"createdBy"`
	}

	var b body
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		fmt.Println(err)
		printInternalError(w)
		return
	}

	if b.Key == "" || b.Value == "" || b.CreatedBy == 0 {
		code := 400
		msg := "Some required fields are missing!"
		fmt.Println(msg)
		printCustomError(w, ErrorMessage{code, msg}, code)
		return
	}

	result, err := h.DB.Exec(
		`INSERT INTO Settings (ItemKey, ItemValue, CreatedBy)
		 VALUES (?,?,?)`,
		b.Key,
		b.Value,
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
