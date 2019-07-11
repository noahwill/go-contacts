package models

import (
	u "go-contacts/utils"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// Token : JWT Token Claim
type Token struct {
	UserID uint
	jwt.StandardClaims
}

// Account : User account
type Account struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token";sql:"-"`
}

// AccErr : populates the Account email, returns an error if needed
func AccErr(acc *Account) error {
	return GetDB().Table("accounts").Where("email = ?", acc.Email).First(acc).Error
}

// Validate : Ensures valid email/password sent from clients
func (account *Account) Validate() (map[string]interface{}, bool) {
	if !strings.Contains(account.Email, "@") {
		return u.Message(false, "Email address is required"), false
	}

	if len(account.Password) < 6 {
		return u.Message(false, "Password is required"), false
	}

	temp := &Account{}

	// Not possible at line 87 because this statement does not save the err
	// after it is evaluated. At 96, the err from 87 is needed again
	if err := AccErr(temp); err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error, Please retry"), false
	}

	if temp.Email != "" {
		return u.Message(false, "Email address already in use by another user."), false
	}

	return u.Message(false, "Requirement Passed"), true
}

// Create : new account and JWT token to send back to client
func (account *Account) Create() map[string]interface{} {
	if resp, ok := account.Validate(); !ok {
		return resp
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
	account.Password = string(hashedPassword)

	GetDB().Create(account)

	if account.ID <= 0 {
		return u.Message(false, "Failed to create account, connection error.")
	}

	// Create new JWT token for the newly registered account
	tk := &Token{UserID: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString

	account.Password = "" // delete password

	response := u.Message(true, "Account has been created")
	response["account"] = account
	return response
}

// Login : authenticate an existing user, generate JWT token if authenticated
func Login(email, password string) map[string]interface{} {
	account := &Account{}
	err := AccErr(account)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email address not found")
		}
		return u.Message(false, "Connection error, Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { // password does not match
		return u.Message(false, "Invalid login credentials. Please try again")
	}

	// Success,
	account.Password = ""

	//Create JWT token
	tk := &Token{UserID: account.ID}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	account.Token = tokenString // put token in the response

	resp := u.Message(true, "Logged In")
	resp["account"] = account
	return resp
}

// GetUser ...
func GetUser(u int) *Account {
	acc := &Account{}
	GetDB().Table("accounts").Where("id = ?", u).First(acc)
	if acc.Email == "" {
		return nil
	}

	// Success
	acc.Password = ""
	return acc
}
