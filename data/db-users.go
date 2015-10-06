package data

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type User struct {
	Userid int

	Username, Email, Role, Joined, Salt, Hash string
}

func GetSalt(username string) (string, error) {
	var salt string
	err := db.QueryRow("SELECT salt FROM USERS WHERE username=$1", username).Scan(&salt)
	if err != nil {
		return "", err
	}
	return salt, nil
}

func DoHashesMatch(username, provided string) bool {
	var actual string
	err := db.QueryRow("SELECT hash FROM USERS WHERE username=$1", username).Scan(&actual)
	if err != nil {
		panic(err)
	}
	return actual == provided
}

func AddNewUser(user *User) sql.Result {
	result, err := db.Exec(
		"INSERT INTO USERS (username, email, role, salt, hash) VALUES($1, $2, $3, $4, $5)",
		user.Username, user.Email, user.Role, user.Salt, user.Hash)
	if err != nil {
		panic(err)
	}
	return result
}

func RemoveUser(userid int) sql.Result {
	result, err := db.Exec("DELETE FROM USERS WHERE userid=?", userid)
	if err != nil {
		panic(err)
	}
	return result
}

func GetUserFromId(id int) (*User, error) {
	var username, email, role string
	var joined time.Time
	err := db.QueryRow("SELECT username, email, role, joined from users WHERE userid=$1;", id).Scan(&username, &email, &role, &joined)
	if err != nil {
		return nil, err
	}
	return &User{id, username, email, role, fmt.Sprintf("%d-%02d-%02d", joined.Year(), joined.Month(), joined.Day()), "", ""}, nil
}

func GetUserId(username string) int {
	var userid int
	err := db.QueryRow("SELECT userid from users WHERE username=$1;", username).Scan(&userid)
	if err != nil {
		return -1
	}
	return userid
}

func GetAllUsers(page int) ([]*User, error) {
	var count int
	_ = db.QueryRow("SELECT COUNT(*) FROM USERS").Scan(&count)
	rows, err := db.Query("SELECT userid, username, email, role, joined FROM USERS ORDER BY userid offset $1 rows fetch next 50 rows only;", page)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	var users = make([]*User, count)
	i := 0
	for rows.Next() {
		var userid int
		var username, email, role string
		var joined time.Time
		err = rows.Scan(&userid, &username, &email, &role, &joined)
		if err != nil {
			return nil, err
		}
		users[i] = &User{userid, username, email, role, fmt.Sprintf("%d-%02d-%02d", joined.Year(), joined.Month(), joined.Day()), "", ""}
		i++
	}
	return users, err
}
