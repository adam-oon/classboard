package models

import (
	"database/sql"
	"fmt"
)

type User struct {
	Id       int
	Username string
	Password string
	Type     string
	Name     string
}

type UserModel struct {
	Db *sql.DB
}

func (userModel UserModel) GetUserByUsername(username string) User {
	rows, err := userModel.Db.Query("SELECT * FROM users WHERE username = ?", username)
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

func (userModel UserModel) GetUser(user_id int) User {
	rows, err := userModel.Db.Query("SELECT * FROM users WHERE id = ?", user_id)
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
