package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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
			h.updateSettingByKey,
		},
	}
	return routes
}

type setting struct {
	ID            int    `json:"id"`
	Key           string `json:"key"`
	Value         string `json:"value"`
	Description   string `json:"description"`
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
			SettingKey,
			Value,
			Description
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
			&s.Key,
			&s.Value,
			&s.Description,
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
			SettingKey,
			Value,
			Description,
			CreatedAt,
			CreatedBy,
			COALESCE(UpdatedAt, '') AS UpdatedAt,
			COALESCE(UpdatedBy, 0) AS UpdatedBy,
			COALESCE(UpdateComment, '') AS UpdateComment
		FROM Settings
		WHERE SettingKey = ?`, key).Scan(
		&s.ID,
		&s.Key,
		&s.Value,
		&s.Description,
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

func (h Handlers) updateSettingByKey(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Value         string `json:"value,omitempty"`
		Description   string `json:"description,omitempty"`
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

	if b.UpdatedBy == 0 || b.UpdateComment == "" {
		code := 400
		msg := "Some required fields are missing!"
		fmt.Println(msg)
		printCustomError(w, ErrorMessage{code, msg}, code)
		return
	}

	var str strings.Builder
	fmt.Fprint(&str, "UPDATE Settings SET ")
	if b.Value != "" {
		fmt.Fprintf(&str, "Value = '%s', ", b.Value)
	}
	if b.Description != "" {
		fmt.Fprintf(&str, "Description = '%s', ", b.Description)
	}

	fmt.Fprintf(&str, `UpdatedAt = CURRENT_TIMESTAMP(), UpdatedBy = %d, UpdateComment = "%s" WHERE SettingKey = "%s"`, b.UpdatedBy, b.UpdateComment, key)

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

func (h Handlers) addSetting(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Key         string `json:"key"`
		Value       string `json:"value"`
		Description string `json:"description"`
		CreatedBy   int    `json:"createdBy"`
	}

	var b body
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		fmt.Println(err)
		printInternalError(w)
		return
	}

	if b.Key == "" || b.Value == "" || b.Description == "" || b.CreatedBy == 0 {
		code := 400
		msg := "Some required fields are missing!"
		fmt.Println(msg)
		printCustomError(w, ErrorMessage{code, msg}, code)
		return
	}

	result, err := h.DB.Exec(
		`INSERT INTO Settings (SettingKey, Value, Description, CreatedBy)
		 VALUES (?,?,?,?)`,
		b.Key,
		b.Value,
		b.Description,
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
