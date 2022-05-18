package Drivers

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

// DB is a global variable to hold db connection
var DB *sql.DB

var HUI = "SDSD"

func Asd() string {
	return "Test"
}

func ConnectionDB() *sql.DB {

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_DATABASE"))

	DB, err := sql.Open("postgres", psqlInfo)
	//defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	err = DB.Ping()

	if err != nil {
		log.Fatal(err)
	}

	return DB
}
