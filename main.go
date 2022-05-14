package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "secret"
	dbname   = "demo"
)

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	Id    int    `json:"id"`
	Token string `json:"token"`
}

type Error struct {
	Message string `json:"message"`
}

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println(db)

	router := mux.NewRouter()

	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/register", register).Methods("POST")
	router.HandleFunc("/login", login).Methods("POST")

	router.HandleFunc("/profile", TokenVerifyMiddleware(profile)).Methods("GET")

	log.Println("Listen on port localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", router))

}

func index(response http.ResponseWriter, request *http.Request) {

	fmt.Println("Invoked index method")

	response.Write([]byte("Success cc fdf ! ds><> "))

}

func responseWithError(response http.ResponseWriter, status int, error Error) {

	error.Message = "Email is missing"
	response.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(response).Encode(error)


}

func register(response http.ResponseWriter, request *http.Request) {

	var user User

	var error Error
	json.NewDecoder(request.Body).Decode(&user)

	if user.Email == "" {
		error.Message = "Email is missing"
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(error)

		return
	}

	if user.Password == "" {
		error.Message = "Password is missing"
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(error)

		return
	}

	if user.Name == "" {
		error.Message = "Message is missing"
		response.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(response).Encode(error)

		return
	}

}

func login(response http.ResponseWriter, request *http.Request) {
	fmt.Println("Invoked login method")
}

func profile(response http.ResponseWriter, request *http.Request) {
	fmt.Println("Invoked profile method")
}

func TokenVerifyMiddleware(next http.HandlerFunc) http.HandlerFunc {

	fmt.Println("Invoked TokenVerifyMiddleware method")
	return nil
}
