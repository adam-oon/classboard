/*
	Package userclass provides SQL query for user_classes table.
*/
package userclass

import (
	classroommodel "classboard/models/classroom"
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
)

type UserClass struct {
	User_id      int
	Classroom_id string
}

type UserClassModel struct {
	Db *sql.DB
}

// JoinClass takes in user_id and classroom_id to insert student's class data into DB.
func (model UserClassModel) JoinClass(user_id int, classroom_id string) error {
	_, err := model.Db.Exec("INSERT INTO user_classes (user_id, classroom_id) VALUES (?,?)", user_id, classroom_id)

	if err != nil {
		if driverErr, ok := err.(*mysql.MySQLError); ok {
			if driverErr.Number == 1062 { //duplicate primary key
				return errors.New("Class already joined!")
			} else if driverErr.Number == 1452 { //cannot update due to invalid classroom_id/user_id
				return errors.New("Invalid Classroom ID")
			}
		}
		return err
	}
	return nil
}

// GetClassroomStudent retrieve user_ids that join the specified classroom_id from DB.
func (model UserClassModel) GetClassroomStudent(classroom_id string) ([]int, error) {
	rows, err := model.Db.Query("SELECT user_id FROM user_classes WHERE classroom_id = ?", classroom_id)

	if err != nil {
		return nil, err
	}
	var user_ids []int
	for rows.Next() {
		var user_id int
		err := rows.Scan(&user_id)
		if err != nil {
			return nil, err
		}
		user_ids = append(user_ids, user_id)
	}

	return user_ids, nil

}

// GetJoinedClass retrieve []Classroom from DB based on user_id given.
func (model UserClassModel) GetJoinedClass(user_id int) ([]classroommodel.Classroom, error) {
	rows, err := model.Db.Query("SELECT classrooms.* FROM user_classes LEFT JOIN classrooms ON user_classes.classroom_id = classrooms.id WHERE user_classes.user_id =  ?", user_id)

	if err != nil {
		return nil, err
	}
	var userClasses []classroommodel.Classroom
	for rows.Next() {
		var userClass classroommodel.Classroom
		err := rows.Scan(&userClass.Id, &userClass.User_id, &userClass.Code, &userClass.Title)
		if err != nil {
			return nil, err
		}
		userClasses = append(userClasses, userClass)
	}

	return userClasses, nil
}

// IsBelongToClassroom checks specified user_id is join under specified classroom_id.
func (model UserClassModel) IsBelongToClassroom(user_id int, classroom_id string) bool {
	rows, err := model.Db.Query("SELECT COUNT(user_id) as totalUserID FROM user_classes WHERE user_id =  ? AND classroom_id = ?", user_id, classroom_id)
	if err != nil {
		return false
	}
	var count int
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return false
		}
	}

	if count > 0 {
		return true
	} else {
		return false
	}
}
