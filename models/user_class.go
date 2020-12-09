package models

import (
	"classboard/config"
	"database/sql"
	"errors"
	"fmt"
)

type UserClass struct {
	User_id      int
	Classroom_id string
}

type UserClassModel struct {
	Db *sql.DB
}

func JoinClass(user_id int, classroom_id string) error {
	db := config.GetMySQLDB()
	defer db.Close()

	query := fmt.Sprintf("INSERT INTO user_classes (user_id, classroom_id) VALUES (%d,'%s')",
		user_id, classroom_id)
	_, err := db.Exec(query)

	if err != nil {
		fmt.Println(err)
		return errors.New("Failed to join classes")
	}
	return nil
}
