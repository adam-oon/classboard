package models

import (
	"classboard/config"
	"fmt"
)

type User struct {
	Id       int
	Username string
	Password string
	Type     string
	Name     string
}

func GetUserByUsername(username string) User {
	db := config.GetMySQLDB()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM users WHERE username = ?", username)
	defer rows.Close()
	if err != nil {
		return User{}
	}

	var user User
	for rows.Next() {
		rows.Scan(&user.Id, &user.Username, &user.Password, &user.Type, &user.Name)
		break
	}
	return user
}

func GetUser(user_id int) User {
	db := config.GetMySQLDB()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM users WHERE id = ?", user_id)
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		return User{}
	}

	var user User
	for rows.Next() {
		rows.Scan(&user.Id, &user.Username, &user.Password, &user.Type, &user.Name)
		break
	}
	return user
}
