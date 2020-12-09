package models

import (
	"database/sql"
)

type Session struct {
	Session_id string
	User_id    int
}

type SessionModel struct {
	Db *sql.DB
}

func (sessionModel SessionModel) CheckSession(id string) bool {
	rows, err := sessionModel.Db.Query("SELECT * FROM sessions WHERE session_id = ?", id)
	defer rows.Close()
	if err != nil {
		return false
	}

	var session_id string
	var user_id int
	for rows.Next() {
		rows.Scan(&session_id, &user_id)
	}

	if session_id != "" && user_id != 0 {
		return true
	}
	return false
}

func (sessionModel SessionModel) GetUserID(session_id string) int {
	rows, err := sessionModel.Db.Query("SELECT user_id FROM sessions WHERE session_id = ?", session_id)
	defer rows.Close()
	if err != nil {
		panic(err.Error())
	}

	var user_id int
	for rows.Next() {
		rows.Scan(&user_id)
	}

	return user_id
}

func (sessionModel SessionModel) DeleteSession(user_id int) bool {
	result, err := sessionModel.Db.Exec("DELETE FROM sessions WHERE user_id =?", user_id)
	if err != nil {
		panic(err.Error())
	}

	if changed, _ := result.RowsAffected(); changed == 1 {
		return true
	} else {
		return false
	}
}

func (sessionModel SessionModel) DeleteSessionByID(session_id string) bool {
	result, err := sessionModel.Db.Exec("DELETE FROM sessions WHERE session_id =?", session_id)
	if err != nil {
		panic(err.Error())
	}

	if changed, _ := result.RowsAffected(); changed == 1 {
		return true
	} else {
		return false
	}
}
