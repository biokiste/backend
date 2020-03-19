package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/didi/gendry/scanner"
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

type tableRow struct {
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
			COALESCE(UpdatedBy, -1) AS UpdatedBy,
			COALESCE(UpdateComment, '') AS UpdateComment
		FROM Settings
	`)
	if err != nil {
		printDbError(w)
		return
	}
	defer results.Close()

	var rows []tableRow

	err = scanner.Scan(results, &rows)
	if err != nil {
		printDbError(w)
		return
	}

	if len(rows) == 0 {
		rows = make([]tableRow, 0)
	}

	printJSON(w, &rows)
}

func (h Handlers) getSettingByKey(w http.ResponseWriter, r *http.Request) {
	key, _ := mux.Vars(r)["key"]

	var tr tableRow
	if err := h.DB.QueryRow(`
		SELECT 
			ID,
			ItemKey,
			ItemValue,
			CreatedAt,
			CreatedBy,
			COALESCE(UpdatedAt, '') AS UpdatedAt,
			COALESCE(UpdatedBy, -1) AS UpdatedBy,
			COALESCE(UpdateComment, '') AS UpdateComment
		FROM Settings
		WHERE ItemKey = ?`, key).Scan(
		&tr.ID,
		&tr.ItemKey,
		&tr.ItemValue,
		&tr.CreatedAt,
		&tr.CreatedBy,
		&tr.UpdatedAt,
		&tr.UpdatedBy,
		&tr.UpdateComment,
	); err != nil {
		fmt.Println(err)
		printCustomError(w, nil, 404)
		return
	}
	printJSON(w, &tr)
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
		printDbError(w)
		return
	}

	// TODO: Validate data – e.g. if all fields set

	date := time.Now()

	result, err := h.DB.Exec(
		`UPDATE Settings
			SET ItemValue = ?, UpdatedAt = ?, UpdatedBy = ?, UpdateComment = ?
		 WHERE ItemKey = ?`,
		b.Value,
		date,
		b.UpdatedBy,
		b.UpdateComment,
		key,
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

func (h Handlers) addSetting(w http.ResponseWriter, r *http.Request) {
	type body struct {
		Key       string `json:"key"`
		Value     string `json:"value"`
		CreatedBy int    `json:"createdBy"`
	}

	var b body
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		printDbError(w)
		return
	}

	// TODO: Validate data – e.g. if all fields set

	result, err := h.DB.Exec(
		`INSERT INTO Settings (ItemKey, ItemValue, CreatedBy)
		 VALUES (?,?,?)`,
		b.Key,
		b.Value,
		b.CreatedBy,
	)

	if err != nil {
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
