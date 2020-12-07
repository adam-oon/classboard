package models

import (
	"database/sql"
	"os"
)

type Session struct {
	Session_id string
	Username   string
}

type SessionModel struct {
	Db *sql.DB
}

func (sessionModel SessionModel) CheckSession(id string) bool {
	rows, err := sessionModel.Db.Query("SELECT *  FROM "+os.Getenv("DB_SCHEMA")+".sessions WHERE session_id = ?", id)
	defer rows.Close()
	if err != nil {
		return false
	}

	var session_id, username string
	for rows.Next() {
		rows.Scan(&session_id, &username)
	}

	if session_id != "" && username != "" {
		return true
	}
	return false
}

func (sessionModel SessionModel) DeleteSession(username string) bool {
	result, err := sessionModel.Db.Exec("DELETE FROM "+os.Getenv("DB_SCHEMA")+".sessions WHERE username =?", username)
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
	result, err := sessionModel.Db.Exec("DELETE FROM "+os.Getenv("DB_SCHEMA")+".sessions WHERE session_id =?", session_id)
	if err != nil {
		panic(err.Error())
	}

	if changed, _ := result.RowsAffected(); changed == 1 {
		return true
	} else {
		return false
	}
}
