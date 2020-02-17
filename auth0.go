package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

const usersByMail = "api/v2/users-by-email?email="
const users = "api/v2/users"

func getToken() (string, error) {
	auth0URI := viper.GetString("auth0URI")
	clientID := viper.GetString("clientId")
	clientSecret := viper.GetString("clientSecret")
	clientAudience := viper.GetString("audience")

	payload := strings.NewReader("grant_type=client_credentials&client_id=" + clientID + "&client_secret=" + clientSecret + "&audience=" + clientAudience)
	req, _ := http.NewRequest("POST", auth0URI+"oauth/token", payload)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return res.Status, err
	}

	defer res.Body.Close()
	var authOBearer Auth0Bearer
	json.NewDecoder(res.Body).Decode(&authOBearer)

	return authOBearer.AccessToken, nil
}

// GetAuth0UsersByEmail returns Auth0 user data
func (h Handlers) GetAuth0UsersByEmail(email string) ([]User, error) {
	auth0URI := viper.GetString("auth0URI")
	apikey, err := getToken()
	var users []User

	if err != nil {
		return users, err
	}

	req, err := http.NewRequest("GET", auth0URI+usersByMail+email, nil)
	fmt.Println(req)
	req.Header.Add("Authorization", "Bearer "+apikey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&users)

	return users, nil
}

// CreateAuth0User creates user at auth0
func (h Handlers) CreateAuth0User(user Auth0User) int {
	auth0URI := viper.GetString("auth0URI")
	apikey, err := getToken()

	if err != nil {
		return 500
	}

	jsonValue, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", auth0URI+users, bytes.NewBuffer(jsonValue))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+apikey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 500
	}
	defer resp.Body.Close()

	return resp.StatusCode
}
