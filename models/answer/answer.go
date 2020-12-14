package answer

import (
	"classboard/config"
	"database/sql"
	"errors"
)

type Answer struct {
	Question_id int
	User_id     int
	Answer      string
	Is_correct  bool
}

type AnswerModel struct {
	Db *sql.DB
}

func SaveAnswer(answer Answer) error {
	db := config.GetMySQLDB()
	defer db.Close()

	_, err := db.Exec("INSERT INTO answers (question_id, user_id, answer,is_correct) VALUES (?,?,?,?)",
		answer.Question_id, answer.User_id, answer.Answer, answer.Is_correct)

	if err != nil {
		return errors.New("Failed to save answer")
	}
	return nil
}

func GetAnswer(question_id int, user_id int) (*Answer, error) {
	db := config.GetMySQLDB()
	defer db.Close()

	var answer Answer
	row := db.QueryRow("SELECT * from answers WHERE question_id = ? AND user_id = ?", question_id, user_id)
	err := row.Scan(&answer.Question_id, &answer.User_id, &answer.Answer, &answer.Is_correct)
	switch err {
	case sql.ErrNoRows: // no result
		return nil, nil
	case nil: // result found
		pointer_answer := &answer
		return pointer_answer, nil
	default: // sql error
		return nil, err
	}
}

func DeleteAnswer(question_id int) error {
	db := config.GetMySQLDB()
	defer db.Close()

	_, err := db.Exec("DELETE FROM answers WHERE question_id =?", question_id)
	if err != nil {
		return err
	}
	return nil
}
