package classroom

import (
	"database/sql"
)

type Classroom struct {
	Id      string
	User_id int
	Title   string
	Code    string
}

type ClassroomModel struct {
	Db *sql.DB
}

func (model ClassroomModel) GetClassroomsByUserId(user_id int) ([]Classroom, error) {
	var classrooms []Classroom
	rows, err := model.Db.Query("SELECT * from classrooms WHERE user_id = ?", user_id)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var classroom Classroom
		err := rows.Scan(&classroom.Id, &classroom.User_id, &classroom.Code, &classroom.Title)
		if err != nil {
			return nil, err
		}
		classrooms = append(classrooms, classroom)
	}

	return classrooms, nil
}

func (model ClassroomModel) GetClassroom(classroom_id string) (Classroom, error) {
	rows, err := model.Db.Query("SELECT * from classrooms WHERE id = ?", classroom_id)
	defer rows.Close()
	if err != nil {
		return Classroom{}, err
	}

	var classroom Classroom
	for rows.Next() {
		err := rows.Scan(&classroom.Id, &classroom.User_id, &classroom.Title, &classroom.Code)
		if err != nil {
			return classroom, err
		}
	}

	return classroom, nil
}

func (model ClassroomModel) SaveClassroom(classroom Classroom) error {
	_, err := model.Db.Exec("INSERT INTO classrooms (id, user_id, title, code) VALUES (?, ?, ?, ?)", classroom.Id, classroom.User_id, classroom.Title, classroom.Code)
	if err != nil {
		return err
	}

	return nil
}

func (model ClassroomModel) UpdateClassroom(classroom Classroom) error {
	_, err := model.Db.Exec("UPDATE classrooms SET title = ?, code = ? WHERE id = ?", classroom.Title, classroom.Code, classroom.Id)
	if err != nil {
		return err
	}

	return nil
}
