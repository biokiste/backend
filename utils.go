package main

import (
	"encoding/json"
	"net/http"
)

// ErrorMessage holds code and message
type ErrorMessage struct {
	StatusCode    int    `json:"code"`
	StatusMessage string `json:"message"`
}

func printInternalError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/vnd.api+json; charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(ErrorMessage{500, "Internal Server Error"})
}

func printError(w http.ResponseWriter, err interface{}) {
	w.Header().Set("Content-Type", "application/vnd.api+json; charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(err)
}

func printCustomError(w http.ResponseWriter, err interface{}, status int) {
	w.Header().Set("Content-Type", "application/vnd.api+json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(err)
}

func printDbError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/vnd.api+json; charset=UTF-8")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode("database connection failed")
}

func printSuccess(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/vnd.api+json; charset=UTF-8")
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w)
}

func printJSON(w http.ResponseWriter, obj interface{}) {
	w.Header().Set("Content-Type", "application/vnd.api+json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
	json.NewEncoder(w).Encode(obj)
}

func notAuthorized(w http.ResponseWriter, obj interface{}) {
	w.Header().Set("Content-Type", "application/vnd.api+json; charset=UTF-8")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(obj)
}
