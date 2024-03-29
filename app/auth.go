// JWT creates unique tokens for each authenticated user
// to be included in the header fo the request made by the
// API to the server.

package app

import (
	"context"
	"fmt"
	"go-contacts/models"
	u "go-contacts/utils"
	"net/http"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

// JwTAuthentication : Creates a middleware to intercept requests, check/verify a JWT token,
// Sends an error if the token is invalid or malformed,
// Proceeds to serve the request otherwise
var JwTAuthentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		AuthRespond := func(w http.ResponseWriter, forbidden bool, status bool, message string) {
			response := u.Message(status, message)
			if forbidden {
				w.WriteHeader(http.StatusForbidden)
			}
			u.Respond(w, response)
			return
		}

		// List of endpoints (identifiers) that don't require auth
		notAuth := []string{"/api/user/new", "/api/user/login"}
		requestPath := r.URL.Path

		// Check if request does not need authentication, serve request if not
		for _, value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" { // Missing token, returns with error code 403 Unauthorized
			AuthRespond(w, false, false, "Missing authentication token")
		}

		// Token normally has format 'Bearer {token-body}'
		// Check if the retrieved token matches this format
		split := strings.Split(tokenHeader, " ")
		if len(split) != 2 {
			AuthRespond(w, true, false, "Invalid/Malformed authentication token")
		}

		tokenPart := split[1]
		tk := &models.Token{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil { // Malformed toke, returns with http code 403
			AuthRespond(w, true, false, "Malformed authentication token")
		}

		if !token.Valid { // Invalid token, may not sign onto server
			AuthRespond(w, true, false, "Invalid authentication token")
		}

		// Proceed with request, set caller to the user retrieved from the token
		fmt.Printf("User %d", tk.UserID)
		contxt := context.WithValue(r.Context(), "user", tk.UserID)
		r = r.WithContext(contxt)
		next.ServeHTTP(w, r)
	})
}
