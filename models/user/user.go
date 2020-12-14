package user

import (
	"database/sql"
	"fmt"
)

type User struct {
	Id       int
	Username string
	Password string
	Type     string
	Name     string
}

type UserModel struct {
	Db *sql.DB
}

func (model UserModel) CheckUserByUsername(username string) (int, error) {
	var count int
	rows, err := model.Db.Query("SELECT COUNT(username) as totalUsername FROM users WHERE username = ?", username)
	defer rows.Close()
	if err != nil {
		return count, err
	}

	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return count, err
		}
		break
	}
	return count, nil
}

func (model UserModel) GetUserByUsername(username string) (User, error) {
	rows, err := model.Db.Query("SELECT * FROM users WHERE username = ?", username)
	defer rows.Close()
	if err != nil {
		return User{}, err
	}

	var user User
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Username, &user.Password, &user.Type, &user.Name)
		if err != nil {
			return User{}, err
		}
		break
	}
	return user, nil
}

func (model UserModel) GetUser(user_id int) (User, error) {
	rows, err := model.Db.Query("SELECT * FROM users WHERE id = ?", user_id)
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		return User{}, err
	}

	var user User
	for rows.Next() {
		err := rows.Scan(&user.Id, &user.Username, &user.Password, &user.Type, &user.Name)
		if err != nil {
			return User{}, err
		}
		break
	}
	return user, nil
}

func (model UserModel) SaveUser(username, password, usertype, name string) error {
	_, err := model.Db.Exec("INSERT INTO users ( username, password, type, name) VALUES (?, ?, ?,?)",
		username, password, usertype, name)
	if err != nil {
		return err
	}

	return nil
}
