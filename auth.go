package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/viper"
)

// Auth define routes that need authentication
func Auth(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _ := readTokenFromHeader(r)

		switch {
		case "/api/status" == r.RequestURI:
			inner.ServeHTTP(w, r)
		case "/api/contents" == r.RequestURI:
			inner.ServeHTTP(w, r)
		default:
			_, err := ValidateToken(token)
			if err != nil {
				//inner.ServeHTTP(w, r)
				notAuthorized(w, "not authorized")
			} else {
				inner.ServeHTTP(w, r)
			}
		}

	})
}

// ValidateToken validates request token
func ValidateToken(tokenString string) (bool, error) {

	APIKey := viper.GetString("apikey")

	// check if token is not empty
	if tokenString == "" {
		fmt.Println("empty token")
		return false, errors.New("empty token")
	}
	if tokenString != APIKey {
		return false, errors.New("wrong credentials")
	}
	return true, nil
}

func readTokenFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", nil // empty token
	}

	authHeaderParts := strings.Split(authHeader, " ")
	if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
		return "", errors.New("Authorization header format must be Bearer {token}")
	}

	return authHeaderParts[1], nil
}
