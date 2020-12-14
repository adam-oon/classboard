package session

import (
	"database/sql"
	"errors"
)

type Session struct {
	Session_id string
	User_id    int
}

type SessionModel struct {
	Db *sql.DB
}

func (model SessionModel) CheckSession(id string) bool {
	rows, err := model.Db.Query("SELECT * FROM sessions WHERE session_id = ?", id)
	defer rows.Close()
	if err != nil {
		return false
	}

	var session_id string
	var user_id int
	for rows.Next() {
		err := rows.Scan(&session_id, &user_id)
		if err != nil {
			return false
		}
	}

	if session_id != "" && user_id != 0 {
		return true
	}
	return false
}

func (model SessionModel) GetUserID(session_id string) (int, error) {
	rows, err := model.Db.Query("SELECT user_id FROM sessions WHERE session_id = ?", session_id)
	defer rows.Close()
	if err != nil {
		return 0, err
	}

	var user_id int
	for rows.Next() {
		err := rows.Scan(&user_id)
		if err != nil {
			return 0, err
		}
	}

	return user_id, nil
}

func (model SessionModel) DeleteSession(user_id int) error {
	result, err := model.Db.Exec("DELETE FROM sessions WHERE user_id =?", user_id)
	if err != nil {
		return err
	}

	if changed, _ := result.RowsAffected(); changed == 1 {
		return nil
	} else {
		return errors.New("Session doesn't exist")
	}
}

func (model SessionModel) DeleteSessionByID(session_id string) error {
	result, err := model.Db.Exec("DELETE FROM sessions WHERE session_id =?", session_id)
	if err != nil {
		return err
	}

	if changed, _ := result.RowsAffected(); changed == 1 {
		return nil
	} else {
		return errors.New("Session doesn't exist")
	}
}

func (model SessionModel) SaveSession(session Session) error {
	_, err := model.Db.Exec("INSERT INTO sessions (session_id, user_id) VALUES (?,?)", session.Session_id, session.User_id)
	if err != nil {
		return err
	}

	return nil
}
