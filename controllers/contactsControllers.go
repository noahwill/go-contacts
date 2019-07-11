package controllers

import (
	"encoding/json"
	"go-contacts/models"
	u "go-contacts/utils"
	"net/http"
)

// CreateContact : Decode JSON body into a Contact struct
var CreateContact = func(w http.ResponseWriter, r *http.Request) {

	// user_id that sent the request
	user := r.Context().Value("user").(uint)
	contact := &models.Contact{}

	if err := json.NewDecoder(r.Body).Decode(contact); err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	contact.UserID = user
	resp := contact.Create()
	u.Respond(w, resp)
}

// GetContactsFor ...
var GetContactsFor = func(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("user").(uint)
	data := models.GetContacts(id)
	resp := u.Message(true, "Success")
	resp["data"] = data
	u.Respond(w, resp)
}
