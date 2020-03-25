package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type user struct {
	ID              int    `json:"id"`
	State           string `json:"state"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Street          string `json:"street"`
	StreetNumber    string `json:"streetNumber"`
	Zip             string `json:"zip"`
	Country         string `json:"country"`
	Birthday        string `json:"birthday"`
	EntranceDate    string `json:"entranceDate"`
	LeavingDate     string `json:"leavingDate,omitempty"`
	AdditionalInfos string `json:"additionalInfos,omitempty"`
	LastActivityAt  string `json:"lastActivityAt,omitempty"`
	CreatedAt       string `json:"createdAt"`
	CreatedBy       int    `json:"createdBy"`
	UpdatedAt       string `json:"updatedAt,omitempty"`
	UpdatedBy       int    `json:"updatedBy,omitempty"`
	UpdateComment   string `json:"updateComment,omitempty"`
}

// GetUsersRoutes get all routes of path /users
func GetUsersRoutes(h *Handlers) []Route {
	routes := []Route{
		{
			"get users",
			"GET",
			"/users",
			h.getUsers,
		},
		{
			"get user",
			"GET",
			"/users/{id}",
			h.getUserByID,
		},
		{
			"add user",
			"POST",
			"/users",
			h.addUser,
		},
		{
			"update user",
			"PATCH",
			"/users/{id}",
			h.updateUserByID,
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

func (h *Handlers) getUsers(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	lastActiveCount := params.Get("last_active")

	var str strings.Builder
	fmt.Fprint(&str, `
		SELECT
			ID,
			State,
			FirstName,
			LastName,
			Email,
			Phone,
			Street,
			StreetNumber,
			Zip,
			Country,
			Birthday,
			EntranceDate,
			COALESCE(LeavingDate, '') as LeavingDate,
			COALESCE(AdditionalInfos, '') as AdditionalInfos,
			COALESCE(LastActivityAt, '') as LastActivityAt,
			CreatedAt,
			CreatedBy,
			COALESCE(UpdatedAt, '') as UpdatedAt,
			COALESCE(UpdatedBy, 0) as UpdatedBy,
			COALESCE(UpdateComment, '') as UpdateComment
		FROM Users`)

	if lastActiveCount != "" {
		date := time.Now().Format("2006-01-02 15:04:05")
		fmt.Fprintf(&str, ` WHERE LastActivityAt <= "%s" ORDER BY LastActivityAt DESC LIMIT %s`, date, lastActiveCount)
	}

	query := str.String()
	results, err := h.DB.Query(query)

	if err != nil {
		fmt.Println(err)
		printInternalError(w)
		return
	}

	var users []user

	for results.Next() {
		var u user
		err = results.Scan(
			&u.ID,
			&u.State,
			&u.FirstName,
			&u.LastName,
			&u.Email,
			&u.Phone,
			&u.Street,
			&u.StreetNumber,
			&u.Zip,
			&u.Country,
			&u.Birthday,
			&u.EntranceDate,
			&u.LeavingDate,
			&u.AdditionalInfos,
			&u.LastActivityAt,
			&u.CreatedAt,
			&u.CreatedBy,
			&u.UpdatedAt,
			&u.UpdatedBy,
			&u.UpdateComment,
		)
		if err != nil {
			fmt.Println(err)
			printInternalError(w)
			return
		}
		users = append(users, u)
	}

	if len(users) == 0 {
		users = make([]user, 0)
	}

	printJSON(w, &users)

}

func (h *Handlers) getUserByID(w http.ResponseWriter, r *http.Request) {
	id, _ := mux.Vars(r)["id"]
	var u user

	query := fmt.Sprintf(`
		SELECT
			ID,
			State,
			FirstName,
			LastName,
			Email,
			Phone,
			Street,
			StreetNumber,
			Zip,
			Country,
			Birthday,
			EntranceDate,
			COALESCE(LeavingDate, '') as LeavingDate,
			COALESCE(AdditionalInfos, '') as AdditionalInfos,
			COALESCE(LastActivityAt, '') as LastActivityAt,
			CreatedAt,
			CreatedBy,
			COALESCE(UpdatedAt, '') as UpdatedAt,
			COALESCE(UpdatedBy, 0) as UpdatedBy,
			COALESCE(UpdateComment, '') as UpdateComment
		FROM Users
		WHERE ID = %s`, id,
	)

	row := h.DB.QueryRow(query)

	row.Scan(
		&u.ID,
		&u.State,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Phone,
		&u.Street,
		&u.StreetNumber,
		&u.Zip,
		&u.Country,
		&u.Birthday,
		&u.EntranceDate,
		&u.LeavingDate,
		&u.AdditionalInfos,
		&u.LastActivityAt,
		&u.CreatedAt,
		&u.CreatedBy,
		&u.UpdatedAt,
		&u.UpdatedBy,
		&u.UpdateComment,
	)

	if u.ID == 0 {
		code := 404
		msg := "not found"
		printCustomError(w, ErrorMessage{code, msg}, code)
		return
	}
	printJSON(w, u)
}

func (h *Handlers) updateUserByID(w http.ResponseWriter, r *http.Request) {
	type body struct {
		State           string `json:"state,omitempty"`
		FirstName       string `json:"firstName,omitempty"`
		LastName        string `json:"lastName,omitempty"`
		Birthday        string `json:"birthday,omitempty"`
		Password        string `json:"password,omitempty"`
		Email           string `json:"email,omitempty"`
		Phone           string `json:"phone,omitempty"`
		Street          string `json:"street,omitempty"`
		StreetNumber    string `json:"streetNumber,omitempty"`
		Zip             string `json:"zip,omitempty"`
		Country         string `json:"country,omitempty"`
		EntranceDate    string `json:"entranceDate,omitempty"`
		LeavingDate     string `json:"leavingDate,omitempty"`
		AdditionalInfos string `json:"additionalInfos,omitempty"`
		UpdatedBy       int    `json:"updatedBy"`
		UpdateComment   string `json:"updateComment"`
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

	if b.Email != "" || b.Password != "" {
		var u Auth0User
		u.Connection = `Connection: "Username-Password-Authentication"`
		if b.Email != "" {
			u.Email = b.Email
		}
		if b.Password != "" {
			u.Password = b.Password
		}
		err = h.UpdateAuth0User(u, "auth|"+id)
		if err != nil {
			fmt.Println(err)
			printInternalError(w)
			return
		}
	}

	var str strings.Builder

	fmt.Fprint(&str, "UPDATE Users SET ")

	if b.State != "" {
		fmt.Fprintf(&str, `State = "%s", `, b.State)
	}
	if b.FirstName != "" {
		fmt.Fprintf(&str, `FirstName = "%s", `, b.FirstName)
	}
	if b.LastName != "" {
		fmt.Fprintf(&str, `LastName = "%s", `, b.LastName)
	}
	if b.Birthday != "" {
		fmt.Fprintf(&str, `Birthday = "%s", `, b.Birthday)
	}
	if b.Phone != "" {
		fmt.Fprintf(&str, `Phone = "%s", `, b.Phone)
	}
	if b.Street != "" {
		fmt.Fprintf(&str, `Street = "%s", `, b.Street)
	}
	if b.StreetNumber != "" {
		fmt.Fprintf(&str, `StreetNumber = "%s", `, b.StreetNumber)
	}
	if b.Zip != "" {
		fmt.Fprintf(&str, `Zip = "%s", `, b.Zip)
	}
	if b.Country != "" {
		fmt.Fprintf(&str, `Country = "%s", `, b.Country)
	}
	if b.EntranceDate != "" {
		fmt.Fprintf(&str, `EntranceDate = "%s", `, b.EntranceDate)
	}
	if b.LeavingDate != "" {
		fmt.Fprintf(&str, `LeavingDate = "%s", `, b.LeavingDate)
	}
	if b.AdditionalInfos != "" {
		fmt.Fprintf(&str, `AdditionalInfos = "%s", `, b.AdditionalInfos)
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

func deleteUser(db *sql.DB, id int64) (bool, error) {
	query := fmt.Sprintf("DELETE FROM Users WHERE ID = %d", id)
	_, err := db.Exec(query)
	if err != nil {
		return false, err
	}
	return true, nil
}
