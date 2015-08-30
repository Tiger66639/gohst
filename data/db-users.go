package data

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type User struct {
	Userid int

	Username, Email, Role, Salt, Hash string
}

var db *sql.DB

func GetSalt(username string) (string, error) {
	var salt string
	err := db.QueryRow("SELECT salt FROM AUTH WHERE username=$1", username).Scan(&salt)
	if err != nil {
		return "", err
	}
	return salt, nil
}

func DoHashesMatch(username, provided string) bool {
	var actual string
	err := db.QueryRow("SELECT hash FROM AUTH WHERE username=$1", username).Scan(&actual)
	if err != nil {
		panic(err)
	}
	log.Printf("provided: %s\ncorrect:  %s", provided, actual)
	return actual == provided
}

func AddNewUser(user *User) sql.Result {
	result, err := db.Exec(
		"INSERT INTO AUTH (username, email, role, salt, hash) VALUES($1, $2, $3, $4, $5)",
		user.Username, user.Email, user.Role, user.Salt, user.Hash)
	if err != nil {
		panic(err)
	}
	return result
}

func RemoveUser(userid int) sql.Result {
	result, err := db.Exec("DELETE FROM AUTH WHERE userid=?", userid)
	if err != nil {
		panic(err)
	}
	return result
}

func Connect(connection string) {
	var err error
	db, err = sql.Open("postgres", connection)
	if err != nil {
		panic(err)
	}
}
