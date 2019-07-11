package models

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	// Blank Import
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

var db *gorm.DB

// Retrieves connection information from .env, builds connection
// string to then connect to the database
// (init are auto called by Go)
func init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Print(err)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_pass")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")

	// Connection String
	dbURI := fmt.Sprintf("host=%s \n user=%s \n dbname=%s \n sslmode=disable \n password=%s",
		dbHost, username, dbName, password)
	fmt.Println(dbURI)

	conn, err := gorm.Open("postgres", dbURI)
	if err != nil {
		fmt.Print(err)
	}

	db = conn
	db.Debug().AutoMigrate(&Account{}, &Contact{})
}

// GetDB : returns the database
func GetDB() *gorm.DB { return db }
