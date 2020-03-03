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
const usersByID = "api/v2/users/"
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

// Auth0GetUser returns Auth0 user data
func (h Handlers) Auth0GetUser(id string) (Auth0User, error) {
	auth0URI := viper.GetString("auth0URI")
	apikey, err := getToken()
	var user Auth0User

	if err != nil {
		return user, err
	}

	req, err := http.NewRequest("GET", auth0URI+usersByID+id, nil)
	req.Header.Add("Authorization", "Bearer "+apikey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&user)

	return user, nil
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

// UpdateAuth0User updates user at auth0
func (h Handlers) UpdateAuth0User(user Auth0User) error {
	auth0URI := viper.GetString("auth0URI")
	apikey, err := getToken()

	if err != nil {
		return err
	}
	fmt.Println(user)
	jsonValue, _ := json.Marshal(user)
	req, err := http.NewRequest("PATCH", auth0URI+usersByID+user.UserID, bytes.NewBuffer(jsonValue))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+apikey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	fmt.Println(resp)
	defer resp.Body.Close()
	return nil
}
