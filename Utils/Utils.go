package Utils

import (
	"encoding/json"
	"go-sand/Models"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/dgrijalva/jwt-go"
)

func ResponseWithError(response http.ResponseWriter, status int, error Models.Error) {

	response.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(response).Encode(error)
}

func ResponseWithSucces(response http.ResponseWriter, data interface{}) {
	response.Header().Set("Content-type", "application/json")
	response.WriteHeader(http.StatusOK)
	json.NewEncoder(response).Encode(data)
}

func GenerateToken(user Models.User) (string, error) {

	// Create the JWT key used to create the signature
	var jwtKey = []byte(os.Getenv("APP_SECRET"))

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