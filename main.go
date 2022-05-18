package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/subosito/gotenv"
	"go-sand/Controllers"
	"go-sand/Drivers"
	"log"
	"net/http"
)

// DB is a global variable to hold db connection
var DB *sql.DB

func init() {
	gotenv.Load(".env.example")
}

func main() {

	DB = Drivers.ConnectionDB()

	controller := Controllers.Controller{}

	router := mux.NewRouter()

	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/register", controller.Register(DB)).Methods("POST")
	router.HandleFunc("/login", controller.Login(DB)).Methods("POST")

	router.HandleFunc("/profile", controller.TokenVerifyMiddleware(controller.Profile())).Methods("GET")

	log.Println("Listen on port localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", router))

}

func index(response http.ResponseWriter, request *http.Request) {

	fmt.Println("Invoked index method")

	response.Write([]byte("Success"))

}
