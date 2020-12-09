package models

import (
	"database/sql"
	"os"
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

func (userModel UserModel) GetUser(username string) User {
	rows, err := userModel.Db.Query("SELECT * FROM "+os.Getenv("DB_SCHEMA")+".users WHERE username = ?", username)
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
