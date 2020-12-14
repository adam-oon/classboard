package question

import (
	"database/sql"
	"errors"
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

type QuestionModel struct {
	Db *sql.DB
}

func (model QuestionModel) SaveQuestion(question QuestionInput) error {
	_, err := model.Db.Exec("INSERT INTO questions (classroom_id, question, type,choice, solution) VALUES (?,?,?,?,?)", question.Classroom_id, question.Question, question.Type, question.Choice, question.Solution)

	if err != nil {
		return err
	}
	return nil
}

func (model QuestionModel) GetQuestionsByClassroomId(classroom_id string) ([]Question, error) {
	var questions []Question
	rows, err := model.Db.Query("SELECT * from questions WHERE classroom_id = ?", classroom_id)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var question Question
		err := rows.Scan(&question.Id, &question.Classroom_id, &question.Question, &question.Type, &question.Choice, &question.Solution)
		if err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}

	return questions, nil
}

func (model QuestionModel) GetQuestion(question_id int) (Question, error) {

	rows, err := model.Db.Query("SELECT * from questions WHERE id = ?", question_id)
	defer rows.Close()
	if err != nil {
		return Question{}, err
	}

	var question Question
	for rows.Next() {
		err := rows.Scan(&question.Id, &question.Classroom_id, &question.Question, &question.Type, &question.Choice, &question.Solution)
		if err != nil {
			return Question{}, err
		}
	}

	return question, nil
}

func (model QuestionModel) DeleteQuestion(question_id int) error {

	result, err := model.Db.Exec("DELETE FROM questions WHERE id =?", question_id)
	if err != nil {
		return err
	}

	if changed, _ := result.RowsAffected(); changed == 1 {
		return nil
	} else {
		return errors.New("Question doesn't exist")
	}
}
