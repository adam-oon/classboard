package models

import (
	"classboard/config"
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
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

	_, err := db.Exec("INSERT INTO user_classes (user_id, classroom_id) VALUES (?,?)", user_id, classroom_id)

	if err != nil {
		if driverErr, ok := err.(*mysql.MySQLError); ok {
			if driverErr.Number == 1062 { //duplicate primary key
				return errors.New("Class already joined!")
			} else if driverErr.Number == 1452 { //cannot update due to invalid classroom_id/user_id
				return errors.New("Invalid Classroom ID")
			}
		}
		return err
	}
	return nil
}

func GetJoinedClass(user_id int) []Classroom {
	db := config.GetMySQLDB()
	defer db.Close()

	rows, err := db.Query("SELECT classrooms.* FROM user_classes LEFT JOIN classrooms ON user_classes.classroom_id = classrooms.id WHERE user_classes.user_id =  ?", user_id)

	if err != nil {
		// return errors.New("Failed to get classes")
	}
	var userClasses []Classroom
	for rows.Next() {
		var userClass Classroom
		err := rows.Scan(&userClass.Id, &userClass.User_id, &userClass.Code, &userClass.Title)
		if err != nil {
			panic(err.Error())
		}
		userClasses = append(userClasses, userClass)
	}

	return userClasses
}
