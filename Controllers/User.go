package Controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go-sand/Models"
	userRepository "go-sand/Repository/Postgresql"
	"go-sand/Utils"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"os"
	"strings"
)

type Controller struct {
}

func (c Controller) Register(DB *sql.DB) http.HandlerFunc {

	return func(response http.ResponseWriter, request *http.Request) {

		var user Models.User

		var error Models.Error
		json.NewDecoder(request.Body).Decode(&user)

		if user.Email == "" {
			error.Message = "Email is missing"
			Utils.ResponseWithError(response, http.StatusBadRequest, error)
			return
		}

		if user.Password == "" {
			error.Message = "Password is missing"
			Utils.ResponseWithError(response, http.StatusBadRequest, error)
			return
		}

		if user.Name == "" {
			error.Message = "Message is missing"
			Utils.ResponseWithError(response, http.StatusBadRequest, error)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

		if err != nil {
			log.Fatal(err)
		}

		user.Password = string(hash)

		userRepo := userRepository.UserRepository{}

		user = userRepo.Register(DB, user)

		Utils.ResponseWithSucces(response, user)
	}
}

func (c Controller) Login(DB *sql.DB) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		var user Models.User
		var jwt Models.JwtToken
		var error Models.Error

		json.NewDecoder(request.Body).Decode(&user)

		if user.Email == "" {
			error.Message = "Email is missing"
			Utils.ResponseWithError(response, http.StatusBadRequest, error)
			return
		}

		if user.Password == "" {
			error.Message = "Password is missing"
			Utils.ResponseWithError(response, http.StatusBadRequest, error)
			return
		}

		password := user.Password

		userRepo := userRepository.UserRepository{}

		user, err := userRepo.Login(DB, user)

		hashedPassword := user.Password
		if err != nil {
			if err == sql.ErrNoRows {
				error.Message = "The user dosen't exists"
				Utils.ResponseWithError(response, http.StatusBadRequest, error)
				return
			}

			log.Fatal(err)
		}

		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

		if err != nil {
			error.Message = "The password is invalid"
			Utils.ResponseWithError(response, http.StatusBadRequest, error)
			return
		}

		token, err := Utils.GenerateToken(user)

		if err != nil {
			log.Fatal(err)
		}

		jwt.Token = token

		Utils.ResponseWithSucces(response, jwt)

	}
}

func (c Controller) Profile() http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {

		response.Write([]byte("Success"))

		return
	}

}

func (c Controller) TokenVerifyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		var errorObj Models.Error
		authHeader := request.Header.Get("Authorization")

		bearerToken := strings.Split(authHeader, " ")

		if len(bearerToken) == 2 {
			authToken := bearerToken[1]
			token, error := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {

				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There error")
				}
				return []byte(os.Getenv("APP_SECRET")), nil
			})

			if error != nil {
				errorObj.Message = error.Error()

				Utils.ResponseWithError(response, http.StatusBadRequest, errorObj)
				return
			}

			if token.Valid {
				next.ServeHTTP(response, request)
			} else {
				errorObj.Message = error.Error()
				Utils.ResponseWithError(response, http.StatusBadRequest, errorObj)
				return
			}

		} else {
			errorObj.Message = "Invalid token"
			Utils.ResponseWithError(response, http.StatusUnauthorized, errorObj)
			return
		}

	})

}
