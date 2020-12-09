package models

import (
	"classboard/config"
	"database/sql"
	"errors"
	"fmt"
)

type Question struct {
	Id           string
	Classroom_id string
	Question     string
	Type         string
	Choice       string
	Solution     string
}

type QuestionInput struct {
	Classroom_id string
	Question     string
	Type         string
	Choice       string
	Solution     string
}

type QuestionModel struct {
	Db *sql.DB
}

func SaveQuestion(question QuestionInput) error {
	db := config.GetMySQLDB()
	defer db.Close()

	query := fmt.Sprintf("INSERT INTO questions (classroom_id, question, type,choice, solution) VALUES ('%s','%s','%s','%s','%s')",
		question.Classroom_id, question.Question, question.Type, question.Choice, question.Solution)
	_, err := db.Exec(query)

	if err != nil {
		fmt.Println(err)
		return errors.New("Failed to insert question")
	}
	return nil
}
