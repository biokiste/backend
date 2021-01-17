package main

import (
	"encoding/json"
	"net/http"
)

// SimpleResponseBody Body for text response
type SimpleResponseBody struct {
	Text string `json:"text"`
}

// UpdateResponseBody Body for updates
type UpdateResponseBody struct {
	RowsAffected int `json:"rowsAffected"`
}

// InsertResponseBody Body for inserts
type InsertResponseBody struct {
	ID int `json:"id"`
}

// InsertResponseArrayBody Body for array inserts
type InsertResponseArrayBody struct {
	IDs []int64 `json:"ids"`
}

// JSONResponse Response container
type JSONResponse struct {
	StatusCode int
	Body       interface{}
}

func respondWithHTTP(w http.ResponseWriter, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(http.StatusText(statusCode)))
}

func respondWithJSON(w http.ResponseWriter, res JSONResponse) {
	w.Header().Set("Content-Type", "application/vnd.api+json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	if res.StatusCode == 0 {
		res.StatusCode = http.StatusOK
	}
	w.WriteHeader(res.StatusCode)
	json.NewEncoder(w).Encode(res.Body)
}
