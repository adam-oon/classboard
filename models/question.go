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

	_, err := db.Exec("INSERT INTO questions (classroom_id, question, type,choice, solution) VALUES (?,?,?,?,?)", question.Classroom_id, question.Question, question.Type, question.Choice, question.Solution)

	if err != nil {
		fmt.Println(err)
		return errors.New("Failed to insert question")
	}
	return nil
}

func GetQuestionsByClassroomId(classroom_id string) []Question {
	db := config.GetMySQLDB()
	defer db.Close()

	var questions []Question
	rows, err := db.Query("SELECT * from questions WHERE classroom_id = ?", classroom_id)
	defer rows.Close()
	if err != nil {
		panic(err.Error())
	}

	for rows.Next() {
		var question Question
		err := rows.Scan(&question.Id, &question.Classroom_id, &question.Question, &question.Type, &question.Choice, &question.Solution)
		if err != nil {
			panic(err.Error())
		}
		questions = append(questions, question)
	}

	return questions
}
