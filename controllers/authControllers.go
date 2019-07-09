package controllers

import (
	"encoding/json"
	"go-contacts/models"
	u "go-contacts/utils"
	"net/http"
)

// CreateAccount : decode JSON request into an account
var CreateAccount = func(w http.ResponseWriter, r *http.Request) {
	account := &models.Account{}
	// Decode request body into a struct
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := account.Create()
	u.Respond(w, resp)
}

// Authenticate : uses JSON decoded account to log in
var Authenticate = func(w http.ResponseWriter, r *http.Request) {
	account := &models.Account{}
	// Decode request body into a struct
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := models.Login(account.Email, account.Password)
	u.Respond(w, resp)
}
