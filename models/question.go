package models

import (
	"classboard/config"
	"errors"
	"fmt"
)

type Question struct {
	Id           int
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

func GetQuestion(question_id int) Question {
	db := config.GetMySQLDB()
	defer db.Close()

	rows, err := db.Query("SELECT * from questions WHERE id = ?", question_id)
	defer rows.Close()
	if err != nil {
		panic(err.Error())
	}

	var question Question
	for rows.Next() {
		err := rows.Scan(&question.Id, &question.Classroom_id, &question.Question, &question.Type, &question.Choice, &question.Solution)
		if err != nil {
			panic(err.Error())
		}
	}

	return question
}

func DeleteQuestion(question_id int) error {
	db := config.GetMySQLDB()
	defer db.Close()

	result, err := db.Exec("DELETE FROM questions WHERE id =?", question_id)
	if err != nil {
		panic(err.Error()) // error//
	}

	if changed, _ := result.RowsAffected(); changed == 1 {
		return nil
	} else {
		return errors.New("Question doesn't exist")
	}
}
