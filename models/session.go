package models

import "classboard/config"

type Session struct {
	Session_id string
	User_id    int
}

func CheckSession(id string) bool {
	db := config.GetMySQLDB()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM sessions WHERE session_id = ?", id)
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

func GetUserID(session_id string) int {
	db := config.GetMySQLDB()
	defer db.Close()

	rows, err := db.Query("SELECT user_id FROM sessions WHERE session_id = ?", session_id)
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

func DeleteSession(user_id int) bool {
	db := config.GetMySQLDB()
	defer db.Close()

	result, err := db.Exec("DELETE FROM sessions WHERE user_id =?", user_id)
	if err != nil {
		panic(err.Error())
	}

	if changed, _ := result.RowsAffected(); changed == 1 {
		return true
	} else {
		return false
	}
}

func DeleteSessionByID(session_id string) bool {
	db := config.GetMySQLDB()
	defer db.Close()

	result, err := db.Exec("DELETE FROM sessions WHERE session_id =?", session_id)
	if err != nil {
		panic(err.Error())
	}

	if changed, _ := result.RowsAffected(); changed == 1 {
		return true
	} else {
		return false
	}
}
