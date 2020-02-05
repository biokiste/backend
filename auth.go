package main

import (
	"fmt"
	"net/http"

	auth0 "github.com/auth0-community/go-auth0"
	"gopkg.in/square/go-jose.v2"
)

// Auth define routes that need authentication
func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth0Domain := "https://biokiste.eu.auth0.com/"
		jwkClient := auth0.NewJWKClient(auth0.JWKClientOptions{URI: fmt.Sprintf("%s.well-known/jwks.json", auth0Domain)}, nil)
		validator := auth0.NewValidator(auth0.NewConfiguration(jwkClient, []string{}, auth0Domain, jose.RS256), nil)
		token, err := validator.ValidateRequest(r)

		switch {
		case "/api/status" == r.RequestURI:
			next.ServeHTTP(w, r)
		case "/api/contents" == r.RequestURI:
			next.ServeHTTP(w, r)
		default:
			if err != nil {
				fmt.Println(err)
				fmt.Println("Token is not valid:", token)
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
			} else {
				next.ServeHTTP(w, r)
			}
		}

	})
}
