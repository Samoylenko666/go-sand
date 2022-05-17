package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"

	//"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// DB is a global variable to hold db connection
var DB *sql.DB

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

type JwtToken struct {
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
	//defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()

	if err != nil {
		log.Fatal(err)
	}

	DB = db

	router := mux.NewRouter()

	router.HandleFunc("/", index).Methods("GET")
	router.HandleFunc("/register", register).Methods("POST")
	router.HandleFunc("/login", login).Methods("POST")

	router.HandleFunc("/profile", TokenVerifyMiddleware(profile)).Methods("GET")

	log.Println("Listen on port localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", router))

}

func responseWithError(response http.ResponseWriter, status int, error Error) {

	response.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(response).Encode(error)
}

func responseWithSucces(response http.ResponseWriter, data interface{}) {
	response.Header().Set("Content-type", "application/json")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(data)
}

func generateToken(user User) (string, error) {

	// Create the JWT key used to create the signature
	var jwtKey = []byte("secret")

	//  iss - чувствительная к регистру строка или URI, которая является уникальным идентификатором стороны, генерирующей токен
	claims := jwt.MapClaims{
		"email": user.Email,
		"iss":   "WASD",
		"exp":   time.Now().Add(time.Minute * 60).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(jwtKey)

	if err != nil {
		log.Fatal(err)
	}

	return tokenStr, nil

}

func index(response http.ResponseWriter, request *http.Request) {

	fmt.Println("Invoked index method")

	response.Write([]byte("Success cc fdf ! ds><> "))

}

func register(response http.ResponseWriter, request *http.Request) {

	var user User

	var error Error
	json.NewDecoder(request.Body).Decode(&user)

	if user.Email == "" {
		error.Message = "Email is missing"
		responseWithError(response, http.StatusBadRequest, error)
		return
	}

	if user.Password == "" {
		error.Message = "Password is missing"
		responseWithError(response, http.StatusBadRequest, error)
		return
	}

	if user.Name == "" {
		error.Message = "Message is missing"
		responseWithError(response, http.StatusBadRequest, error)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

	if err != nil {
		log.Fatal(err)
	}

	user.Password = string(hash)

	stmt := "insert into users (email, password, name) values($1, $2, $3) RETURNING id;"

	err = DB.QueryRow(stmt, user.Email, user.Password, user.Name).Scan(&user.Id)

	if err != nil {
		error.Message = "Something went wrong"

		responseWithError(response, http.StatusInternalServerError, error)
	}

	user.Password = ""

	responseWithSucces(response, user)
}

func login(response http.ResponseWriter, request *http.Request) {
	var user User
	var jwt JwtToken
	var error Error

	json.NewDecoder(request.Body).Decode(&user)

	if user.Email == "" {
		error.Message = "Email is missing"
		responseWithError(response, http.StatusBadRequest, error)
		return
	}

	if user.Password == "" {
		error.Message = "Password is missing"
		responseWithError(response, http.StatusBadRequest, error)
		return
	}

	password := user.Password

	stmt := "select * from users where email = $1"

	row := DB.QueryRow(stmt, user.Email)

	err := row.Scan(&user.Id, &user.Email, &user.Password, &user.Name)

	if err != nil {
		if err == sql.ErrNoRows {
			error.Message = "The user dosen't exists"
			responseWithError(response, http.StatusBadRequest, error)
			return
		}

		log.Fatal(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		error.Message = "The password is invalid"
		responseWithError(response, http.StatusBadRequest, error)
		return
	}

	token, err := generateToken(user)

	if err != nil {
		log.Fatal(err)
	}

	jwt.Token = token

	responseWithSucces(response, jwt)
	fmt.Println(token)

}

func profile(response http.ResponseWriter, request *http.Request) {
	fmt.Println("Invoked profile method")

	response.Write([]byte("Success zaebis "))

	return
}

func TokenVerifyMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		var errorObj Error
		authHeader := request.Header.Get("Authorization")

		bearerToken := strings.Split(authHeader, " ")

		if len(bearerToken) == 2 {
			authToken := bearerToken[1]
			token, error := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {

				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Pizda")
				}
				return []byte("secret"), nil
			})

			if error != nil {
				errorObj.Message = error.Error()

				responseWithError(response, http.StatusBadRequest, errorObj)
				return
			}

			if token.Valid {
				next.ServeHTTP(response, request)
			} else {
				errorObj.Message = error.Error()
				responseWithError(response, http.StatusBadRequest, errorObj)
				return
			}

		} else {
			errorObj.Message = "Invalid token"
			responseWithError(response, http.StatusUnauthorized, errorObj)
			return
		}

	})

}
