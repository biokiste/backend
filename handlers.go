package main

import (
	"net/http"
)

// ShowStatus delivers actual status
func ShowStatus(w http.ResponseWriter, r *http.Request) {
	printJSON(w, "biokiste api status: ok")
}

// ListUsers delivers user data
func ListUsers(w http.ResponseWriter, r *http.Request) {
	loc := "TODO: should list users"
	printJSON(w, loc)
}
