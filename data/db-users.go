package data

import (
	"database/sql"
)

type User struct {
	Userid int

	Login, Username, Email, Role, Salt, Hash string
}

var db *sql.DB

func GetSalt(username string) string {
	var salt string
	err := db.QueryRow("SELECT salt FROM AUTH WHERE username=?", username).Scan(&salt)
	if err != nil {
		panic(err)
	}
	return salt
}

func DoHashesMatch(username, provided string) bool {
	var actual string
	err := db.QueryRow("SELECT hash FROM AUTH WHERE username=?", username).Scan(&actual)
	if err != nil {
		panic(err)
	}
	return actual == provided
}

func AddNewUser(user *User) sql.Result {
	result, err := db.Exec(
		"INSERT INTO AUTH (login, username, email, role, salt, hash) VALUES($1, $2, $3, $4, $5, $6)",
		user.Login, user.Username, user.Email, user.Role, user.Salt, user.Hash)
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

func connect(connection string) {
	var err error
	db, err = sql.Open("postgres", connection)
	if err != nil {
		panic(err)
	}
}
