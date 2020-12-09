package models

import (
	"database/sql"
)

type Answer struct {
	Question_id int
	User_id     int
	Answer      string
	Marked      bool
}

type AnswerModel struct {
	Db *sql.DB
}
