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
	ID              int     `json:"id"`
	State           string  `json:"state,omitempty"`
	FirstName       string  `json:"firstName,omitempty"`
	LastName        string  `json:"lastName,omitempty"`
	Email           string  `json:"email,omitempty"`
	Phone           string  `json:"phone,omitempty"`
	Street          string  `json:"street,omitempty"`
	StreetNumber    string  `json:"streetNumber,omitempty"`
	Zip             string  `json:"zip,omitempty"`
	Country         string  `json:"country,omitempty"`
	Birthday        string  `json:"birthday,omitempty"`
	EntranceDate    string  `json:"entranceDate,omitempty"`
	LeavingDate     string  `json:"leavingDate,omitempty"`
	AdditionalInfos string  `json:"additionalInfos,omitempty"`
	LastActivityAt  string  `json:"lastActivityAt,omitempty"`
	CreatedAt       string  `json:"createdAt,omitempty"`
	CreatedBy       int     `json:"createdBy,omitempty"`
	UpdatedAt       string  `json:"updatedAt,omitempty"`
	UpdatedBy       int     `json:"updatedBy,omitempty"`
	UpdateComment   string  `json:"updateComment,omitempty"`
	Balance         float32 `json:"balance,omitempty"`
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
		{
			"add user to group",
			"POST",
			"/users/{id}/groups",
			h.addUserToGroup,
		},
		{
			"update group user",
			"PATCH",
			"/users/{id}/groups/{groupId}",
			h.updateGroupUser,
		},
		{
			"remove user from group",
			"DELETE",
			"/users/{id}/groups/{groupId}",
			h.removeUserFromGroup,
		},
	}

	return routes
}

func (h *Handlers) addUser(w http.ResponseWriter, r *http.Request) {
	type body struct {
		ID              int    `json:"id" db:"id"`
		CreatedAt       string `json:"created_at" db:"created_at"`
		UpdatedAt       string `json:"updated_at" db:"updated_at"`
		State           string `json:"state" db:"state"`
		Firstname       string `json:"firstname" db:"firstname"`
		Lastname        string `json:"lastname" db:"lastname"`
		Email           string `json:"email" db:"email"`
		Phone           string `json:"phone" db:"phone"`
		Street          string `json:"street" db:"street"`
		Zip             string `json:"zip" db:"zip"`
		Country         string `json:"country" db:"country"`
		Birthday        string `json:"birthday" db:"birthday"`
		EntranceDate    string `json:"entrance_date" db:"entrance_date"`
		LeavingDate     string `json:"leaving_date" db:"leaving_date"`
		AdditionalInfos string `json:"additional_infos" db:"additional_infos"`
		LastActivityAt  string `json:"last_activity_at" db:"last_activity_at"`
		Password        string `json:"password" db:"password"`
	}

	var b body
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		fmt.Println(err)
		respondWithHTTP(w, http.StatusInternalServerError)
		return
	}

	if b.State == "" || b.Firstname == "" || b.Lastname == "" || b.Password == "" || b.Email == "" || b.Phone == "" || b.Street == "" || b.Zip == "" || b.Country == "" {
		err := SimpleResponseBody{"Some required fields are missing!"}
		fmt.Println(err.Text)
		respondWithJSON(w, JSONResponse{http.StatusBadRequest, &err})
		return
	}

	var insert strings.Builder
	var values strings.Builder

	fmt.Fprint(&insert, "INSERT INTO users (state, firstname, lastname, birthday, email, phone, street, zip, country, entrance_date, created_at, updated_at")
	fmt.Fprintf(&values, `VALUES ("%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s", "%s"`, b.State, b.Firstname, b.Lastname, b.Birthday, b.Email, b.Phone, b.Street, b.Zip, b.Country, b.EntranceDate, b.CreatedAt, b.UpdatedAt)

	if b.LeavingDate != "" {
		fmt.Fprint(&insert, ", leaving_date")
		fmt.Fprintf(&values, ", %q", b.LeavingDate)
	}

	if b.AdditionalInfos != "" {
		fmt.Fprint(&insert, ", additional_infos")
		fmt.Fprintf(&values, ", %s", b.AdditionalInfos)
	}

	fmt.Fprint(&insert, ")")
	fmt.Fprint(&values, ")")

	query := fmt.Sprintf("%s %s", insert.String(), values.String())

	fmt.Println(query)

	result, err := h.DB.Exec(query)
	if err != nil {
		fmt.Println(err)
		respondWithHTTP(w, http.StatusInternalServerError)
		return
	}
	id, _ := result.LastInsertId()

	auth0User := Auth0User{
		UserID:      strconv.Itoa(int(id)),
		Password:    b.Password,
		Email:       b.Email,
		Connection:  "Username-Password-Authentication",
		VerifyEmail: false,
	}

	fmt.Println("try creating auth0 user: ", auth0User)

	token, err := getToken()

	if err != nil {
		deleteUser(h.DB, id)
		text := fmt.Sprintf(`Creating user at auth provider failed with "%s"`, err)
		respondWithJSON(w, JSONResponse{http.StatusInternalServerError, SimpleResponseBody{text}})
		return
	}

	_, err = h.CreateAuth0User(auth0User, token)

	if err != nil {
		deleteUser(h.DB, id)
		text := fmt.Sprintf(`Creating user at auth provider failed with "%s"`, err)
		respondWithJSON(w, JSONResponse{http.StatusInternalServerError, SimpleResponseBody{text}})
		return
	}

	respondWithJSON(w, JSONResponse{Body: InsertResponseBody{int(id)}})
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
		respondWithHTTP(w, http.StatusInternalServerError)
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
			respondWithHTTP(w, http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}

	if len(users) == 0 {
		users = make([]user, 0)
	}

	respondWithJSON(w, JSONResponse{Body: &users})
}

func (h *Handlers) getUserByID(w http.ResponseWriter, r *http.Request) {
	id, _ := mux.Vars(r)["id"]
	var u user

	query := fmt.Sprintf(`
		SELECT
			Users.ID,
			Users.State,
			Users.FirstName,
			Users.LastName,
			Users.Email,
			Users.Phone,
			Users.Street,
			Users.StreetNumber,
			Users.Zip,
			Users.Country,
			Users.Birthday,
			Users.EntranceDate,
			COALESCE(Users.LeavingDate, '') as LeavingDate,
			COALESCE(Users.AdditionalInfos, '') as AdditionalInfos,
			COALESCE(Users.LastActivityAt, '') as LastActivityAt,
			Users.CreatedAt,
			Users.CreatedBy,
			COALESCE(Users.UpdatedAt, '') as UpdatedAt,
			COALESCE(Users.UpdatedBy, 0) as UpdatedBy,
			COALESCE(Users.UpdateComment, '') as UpdateComment,
			CAST(SUM(Transactions.Amount) AS DECIMAL(5,2)) AS Balance
		FROM Users
		INNER JOIN Transactions ON Users.ID=Transactions.UserID
		WHERE Users.ID = %s`, id,
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
		&u.Balance,
	)

	if u.ID == 0 {
		respondWithHTTP(w, http.StatusNotFound)
		return
	}
	respondWithJSON(w, JSONResponse{Body: u})
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
		respondWithHTTP(w, http.StatusInternalServerError)
		return
	}

	if b.UpdatedBy == 0 || b.UpdateComment == "" {
		err := SimpleResponseBody{"Some required fields are missing!"}
		fmt.Println(err.Text)
		respondWithJSON(w, JSONResponse{http.StatusBadRequest, &err})
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
			respondWithHTTP(w, http.StatusInternalServerError)
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
		respondWithHTTP(w, http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()

	respondWithJSON(w, JSONResponse{Body: UpdateResponseBody{int(rowsAffected)}})
}

func (h *Handlers) addUserToGroup(w http.ResponseWriter, r *http.Request) {

	type body struct {
		GroupID  int  `json:"groupId"`
		IsLeader bool `json:"isLeader,omitempty"`
	}

	userID, _ := mux.Vars(r)["id"]

	var b body
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		fmt.Println(err)
		respondWithHTTP(w, http.StatusInternalServerError)
		return
	}

	if b.GroupID == 0 {
		err := SimpleResponseBody{"Some required fields are missing!"}
		fmt.Println(err.Text)
		respondWithJSON(w, JSONResponse{http.StatusBadRequest, &err})
		return
	}

	isLeader := 0
	if b.IsLeader {
		isLeader = 1
	}

	result, err := h.DB.Exec(
		`INSERT INTO groupusers (group_id, user_id, is_leader)
		 VALUES (?,?,?)`,
		b.GroupID,
		userID,
		isLeader,
	)

	if err != nil {
		respondWithHTTP(w, http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()

	if err != nil {
		fmt.Println(err)
		respondWithHTTP(w, http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, JSONResponse{Body: InsertResponseBody{int(id)}})
}

func (h *Handlers) updateGroupUser(w http.ResponseWriter, r *http.Request) {
	type body struct {
		IsLeader bool `json:"isLeader"`
	}

	userID, _ := mux.Vars(r)["id"]
	groupID, _ := mux.Vars(r)["groupId"]

	var b body
	err := json.NewDecoder(r.Body).Decode(&b)
	if err != nil {
		fmt.Println(err)
		respondWithHTTP(w, http.StatusInternalServerError)
		return
	}

	isLeader := 0
	if b.IsLeader {
		isLeader = 1
	}

	query := fmt.Sprintf("UPDATE GroupUsers SET IsLeader = %d WHERE GroupID = %s AND UserID = %s", isLeader, groupID, userID)

	result, err := h.DB.Exec(query)

	if err != nil {
		fmt.Println(err)
		respondWithHTTP(w, http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()

	respondWithJSON(w, JSONResponse{Body: UpdateResponseBody{int(rowsAffected)}})
}

func (h *Handlers) removeUserFromGroup(w http.ResponseWriter, r *http.Request) {
	userID, _ := mux.Vars(r)["id"]
	groupID, _ := mux.Vars(r)["groupId"]

	query := fmt.Sprintf("DELETE FROM GroupUsers WHERE GroupID = %s AND UserID = %s", groupID, userID)

	result, err := h.DB.Exec(query)

	if err != nil {
		fmt.Println(err)
		respondWithHTTP(w, http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()

	respondWithJSON(w, JSONResponse{Body: UpdateResponseBody{int(rowsAffected)}})
}

func deleteUser(db *sql.DB, id int64) (bool, error) {
	query := fmt.Sprintf("DELETE FROM Users WHERE ID = %d", id)
	_, err := db.Exec(query)
	if err != nil {
		return false, err
	}
	return true, nil
}
