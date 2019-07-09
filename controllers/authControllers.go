package controllers

import (
	"encoding/json"
	"go-contacts/models"
	u "go-contacts/utils"
	"net/http"
)

// JSONValidate : validates the JSON data from client
func JSONValidate(w http.ResponseWriter, r *http.Request) (*models.Account, error) {
	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account) // decode the request body into struct
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return nil, err
	}

	return account, nil
}

// CreateAccount ...
var CreateAccount = func(w http.ResponseWriter, r *http.Request) {
	account, err := JSONValidate(w, r)
	if err != nil {
		return
	}

	resp := account.Create()
	u.Respond(w, resp)
}

// Authenticate ...
var Authenticate = func(w http.ResponseWriter, r *http.Request) {
	account, err := JSONValidate(w, r)
	if err != nil {
		return
	}

	resp := models.Login(account.Email, account.Password)
	u.Respond(w, resp)
}
