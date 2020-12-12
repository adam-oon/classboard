package models

import (
	"classboard/config"
	"database/sql"
	"errors"
	"fmt"
)

type Classroom struct {
	Id      string
	User_id int
	Title   string
	Code    string
}

type ClassroomModel struct {
	Db *sql.DB
}

type ResMessage struct {
	ResponseText string
	ID           string
}

func GetClassroomsByUserId(user_id int) []Classroom {
	db := config.GetMySQLDB()
	defer db.Close()

	var classrooms []Classroom
	rows, err := db.Query("SELECT * from classrooms WHERE user_id = ?", user_id)
	defer rows.Close()
	if err != nil {
		panic(err.Error())
	}

	for rows.Next() {
		var classroom Classroom
		err := rows.Scan(&classroom.Id, &classroom.User_id, &classroom.Code, &classroom.Title)
		if err != nil {
			panic(err.Error())
		}
		classrooms = append(classrooms, classroom)
	}

	return classrooms
}

func GetClassroom(classroom_id string) Classroom {
	db := config.GetMySQLDB()
	defer db.Close()

	rows, err := db.Query("SELECT * from classrooms WHERE id = ?", classroom_id)
	defer rows.Close()
	if err != nil {
		panic(err.Error())
	}

	var classroom Classroom
	for rows.Next() {
		err := rows.Scan(&classroom.Id, &classroom.User_id, &classroom.Title, &classroom.Code)
		if err != nil {
			panic(err.Error())
		}
	}

	return classroom
}

func SaveClassroom(classroom Classroom) error {
	db := config.GetMySQLDB()
	defer db.Close()

	_, err := db.Exec("INSERT INTO classrooms (id, user_id, title, code) VALUES (?, ?, ?, ?)", classroom.Id, classroom.User_id, classroom.Title, classroom.Code)

	if err != nil {
		fmt.Println(err)
		return errors.New("Failed to insert classroom")
	}
	return nil
}

func UpdateClassroom(classroom Classroom) error {
	db := config.GetMySQLDB()
	defer db.Close()

	_, err := db.Exec("UPDATE classrooms SET title = ?, code = ? WHERE id = ?", classroom.Title, classroom.Code, classroom.Id)

	if err != nil {
		fmt.Println(err)
		return errors.New("Failed to update classroom")
	}
	return nil
}
