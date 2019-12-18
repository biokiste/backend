package main

import (
	"database/sql"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// UsersResponse JSON API Spec Wrapper
type UsersResponse struct {
	Users []User `json:"data"`
}

// User holds properties
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

// ShowStatus delivers actual status
func ShowStatus(w http.ResponseWriter, r *http.Request) {
	printJSON(w, "ok")
}

// ListUsers delivers user data
func ListUsers(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:8889)/foodkoop_biokiste")
	if err != nil {
		printError(w, err.Error())
	}
	defer db.Close()

	results, err := db.Query("SELECT id, username FROM users")
	if err != nil {
		printError(w, err.Error())
	}

	var users []User
	for results.Next() {
		var user User
		err = results.Scan(&user.ID, &user.Username)
		if err != nil {
			printError(w, err.Error())
		}
		users = append(users, user)
	}

	printJSON(w, &UsersResponse{Users: users})
}
