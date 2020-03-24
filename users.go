package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// GetUsersRoutes get all routes of path /users
func GetUsersRoutes(h *Handlers) []Route {
	routes := []Route{
		{
			"add group",
			"POST",
			"/users",
			h.addUser,
		},
		{
			"add user",
			"POST",
			"/user/auth/create",
			h.addUser,
		},
	}

	return routes
}

func (h *Handlers) addUser(w http.ResponseWriter, r *http.Request) {
	type body struct {
		State           string `json:"state"`
		FirstName       string `json:"firstName"`
		LastName        string `json:"lastName"`
		Birthday        string `json:"birthday"`
		Password        string `json:"password"`
		Email           string `json:"email"`
		Phone           string `json:"phone"`
		Street          string `json:"street"`
		StreetNumber    string `json:"streetNumber"`
		Zip             string `json:"zip"`
		Country         string `json:"country"`
		EntranceDate    string `json:"entranceDate"`
		LeavingDate     string `json:"leavingDate,omitempty"`
		AdditionalInfos string `json:"additionalInfos,omitempty"`
		CreatedBy       int    `json:"createdBy"`
	}

	var b body
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		fmt.Println(err)
		printInternalError(w)
		return
	}

	if b.State == "" || b.FirstName == "" || b.LastName == "" || b.Birthday == "" || b.Password == "" || b.Email == "" || b.Phone == "" || b.Street == "" || b.StreetNumber == "" || b.Zip == "" || b.Country == "" || b.EntranceDate == "" || b.CreatedBy == 0 {
		code := 400
		msg := "Some required fields are missing!"
		fmt.Println(msg)
		printCustomError(w, ErrorMessage{code, msg}, code)
		return
	}

	var insert strings.Builder
	var values strings.Builder

	fmt.Fprint(&insert, "INSERT INTO Users (State, FirstName, LastName, Birthday, Email, Phone, Street, StreetNumber, Zip, Country, EntranceDate, CreatedBy")
	fmt.Fprintf(&values, `VALUES ("%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", %d`, b.State, b.FirstName, b.LastName, b.Birthday, b.Email, b.Phone, b.Street, b.StreetNumber, b.Zip, b.Country, b.EntranceDate, b.CreatedBy)

	if b.LeavingDate != "" {
		fmt.Fprint(&insert, ", LeavingDate")
		fmt.Fprintf(&values, ", %s", b.LeavingDate)
	}

	if b.AdditionalInfos != "" {
		fmt.Fprint(&insert, ", AdditionalInfos")
		fmt.Fprintf(&values, ", %s", b.AdditionalInfos)
	}

	fmt.Fprint(&insert, ")")
	fmt.Fprint(&values, ")")

	query := fmt.Sprintf("%s %s", insert.String(), values.String())

	result, err := h.DB.Exec(query)
	if err != nil {
		fmt.Println(err)
		printError(w, err)
		return
	}
	id, _ := result.LastInsertId()

	auth0User := Auth0User{
		UserID:     strconv.Itoa(int(id)),
		Password:   b.Password,
		Email:      b.Email,
		Connection: "Username-Password-Authentication",
	}

	token, err := getToken()

	if err != nil {
		deleteUser(h.DB, id)
		message := fmt.Sprintf(`Creating user at auth provider failed with "%s"`, err)
		code := 500
		fmt.Println(err)
		printCustomError(w, ErrorMessage{code, message}, code)
		return
	}

	_, err = h.CreateAuth0User(auth0User, token)

	if err != nil {
		deleteUser(h.DB, id)
		message := fmt.Sprintf(`Creating user at auth provider failed with "%s"`, err)
		code := 500
		fmt.Println(message)
		printCustomError(w, ErrorMessage{code, message}, code)
		return
	}

	type resBody struct {
		Status       string `json:"status"`
		LastInsertId int    `json:"id"`
	}

	printJSON(w, &resBody{"ok", int(id)})
}

func deleteUser(db *sql.DB, id int64) (bool, error) {
	query := fmt.Sprintf("DELETE FROM Users WHERE ID = %d", id)
	_, err := db.Exec(query)
	if err != nil {
		return false, err
	}
	return true, nil
}
