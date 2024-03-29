package models

import (
	"fmt"
	u "go-contacts/utils"

	"github.com/jinzhu/gorm"
)

// Contact : Name, Phone, UserID
type Contact struct {
	gorm.Model
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	UserID uint   `json:"user_id"`
}

// Validate : validates the passed contact
func (contact *Contact) Validate() (map[string]interface{}, bool) {
	if contact.Name == "" {
		return u.Message(false, "Contact name should be on the payload"), false
	}

	if contact.Phone == "" {
		return u.Message(false, "Phone number should be on the payload"), false
	}

	if contact.UserID <= 0 {
		return u.Message(false, "User is not recognized"), false
	}

	// contact is correctly constructed
	return u.Message(true, "Success"), true
}

// Create : creates a contact
func (contact *Contact) Create() map[string]interface{} {
	if resp, ok := contact.Validate(); !ok {
		return resp
	}

	GetDB().Create(contact)

	resp := u.Message(true, "Success")
	resp["contact"] = contact
	return resp
}

// GetContact : gets a contact associated with a specific id
func GetContact(id uint) *Contact {
	contact := &Contact{}
	if err := GetDB().Table("contacts").Where("id = ?", id).First(contact).Error; err != nil {
		return nil
	}
	return contact
}

// GetContacts : gets all the contacts associated with a user id
func GetContacts(user uint) []*Contact {
	contacts := make([]*Contact, 0)
	if err := GetDB().Table("contacts").Where("user_id = ?", user).Find(&contacts).Error; err != nil {
		fmt.Println(err)
		return nil
	}
	return contacts
}
