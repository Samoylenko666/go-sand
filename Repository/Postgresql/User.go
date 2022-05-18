package userRepository

import (
	"database/sql"
	"go-sand/Models"
	"log"
)

type UserRepository struct{}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (u UserRepository) Register(DB *sql.DB, user Models.User) Models.User {

	stmt := "insert into users (email, password, name) values($1, $2, $3) RETURNING id;"

	err := DB.QueryRow(stmt, user.Email, user.Password, user.Name).Scan(&user.Id)

	logFatal(err)
	user.Password = ""

	return user
}

func (u UserRepository) Login(DB *sql.DB, user Models.User) (Models.User, error) {

	stmt := "select * from users where email = $1"

	row := DB.QueryRow(stmt, user.Email)

	err := row.Scan(&user.Id, &user.Email, &user.Password, &user.Name)

	if err != nil {
		return user, err
	}

	return user, nil
}
