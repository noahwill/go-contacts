package main

import (
	"github.com/gorilla/mux"
	"go-contacts/app"
	"os"
	"fmt"
	"net/http"
	"go-contacts/controllers"
)

// Creates the entry point for the application 
func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")
	router.HandleFunc("/api/contacts/new", controllers.CreateContact).Methods("POST")
	router.HandleFunc("/api/me/contacts", controllers.GetContactsFor).Methods("GET")

	// Attach JWT Authentication Middleware
	router.Use(app.JwTAuthentication)

	// When tested locally, this should return an empty string since no port was specified
	port := os.Getenv("PORT")
	if port == "" { port = "8000"} // localhost

	fmt.Println(port)

	// Launch the app, visit localhost:8000/api
	err := http.ListenAndServe(":" + port, router)
	if err != nil { fmt.Print(err) } 
}