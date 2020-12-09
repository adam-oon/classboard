package models

import (
	"classboard/config"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
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

type ResMessage struct {
	ResponseText string
	ID           string
}

func GetClassroomsByUserId(user_id int) []Classroom {
	db := config.GetMySQLDB()
	defer db.Close()

	var classrooms []Classroom
	rows, err := db.Query("SELECT * from classrooms WHERE user_id = ?", user_id)
	defer rows.Close()
	if err != nil {
		panic(err.Error())
	}

	for rows.Next() {
		var classroom Classroom
		err := rows.Scan(&classroom.Id, &classroom.User_id, &classroom.Code, &classroom.Title)
		if err != nil {
			panic(err.Error())
		}
		classrooms = append(classrooms, classroom)
	}

	return classrooms
}

func GetClassroomOwner(classroom_id string) int {
	var user_id int
	db := config.GetMySQLDB()
	defer db.Close()

	rows, err := db.Query("SELECT user_id from classrooms WHERE id = ?", classroom_id)
	defer rows.Close()
	if err != nil {
		panic(err.Error())
	}

	for rows.Next() {
		err := rows.Scan(&user_id)
		if err != nil {
			panic(err.Error())
		}
	}

	return user_id
}

func SaveClassroom(classroom Classroom) error {
	db := config.GetMySQLDB()
	defer db.Close()

	query := fmt.Sprintf("INSERT INTO classrooms (id, user_id, title, code) VALUES ('%s', %d, '%s', '%s')",
		classroom.Id, classroom.User_id, classroom.Title, classroom.Code)
	_, err := db.Exec(query)

	if err != nil {
		fmt.Println(err)
		return errors.New("Failed to insert classroom")
	}
	return nil
}

// func (classroomModel ClassroomModel) AddClassroom(username string) User {
// 	rows, err := classroomModel.Db.Query("SELECT * FROM "+os.Getenv("DB_SCHEMA")+".users WHERE username = ?", username)
// 	defer rows.Close()
// 	if err != nil {
// 		return User{}
// 	}

// 	var user User
// 	for rows.Next() {
// 		rows.Scan(&user.Id, &user.Username, &user.Password, &user.Type, &user.Name)
// 		break
// 	}
// 	return user
// }

// func (classroomModel ClassroomModel) SaveClassroom(res http.ResponseWriter, req *http.Request) {
// 	res.Header().Set("Content-Type", "application/json")

// 	if req.Method == http.MethodPost && req.Header.Get("Content-type") == "application/json" {

// 		id, _ := uuid.NewV4()
// 		var newClassroom AddClassroom
// 		reqBody, err := ioutil.ReadAll(req.Body)

// 		if err == nil {
// 			err := json.Unmarshal(reqBody, &newClassroom)
// 			if err != nil {
// 				res.WriteHeader(http.StatusUnprocessableEntity)
// 				json.NewEncoder(res).Encode(ResMessage{ResponseText: "Sorry the classroom info is incomplete"})
// 				return
// 			}

// 			myCookie, err := req.Cookie("myCookie")
// 			if err != nil {
// 				res.WriteHeader(http.StatusUnprocessableEntity)
// 				json.NewEncoder(res).Encode(ResMessage{ResponseText: "Sorry the classroom info is incomplete"})
// 				return
// 			}
// 			sessionModel := SessionModel{
// 				Db: db,
// 			}

// 			// 	// input validation
// 			// 	ok := checkValidCourseID(newCourse.Course_id)
// 			// 	if !ok {
// 			// 		res.WriteHeader(http.StatusUnprocessableEntity)
// 			// 		json.NewEncoder(res).Encode(ResMessage{ResponseText: "Sorry the Course ID cannot has spaces"})
// 			// 		return
// 			// 	}
// 			// 	sanitizeCourseInput(&newCourse)
// 			// 	if newCourse.Course_id == "" || newCourse.Title == "" || newCourse.Description == "" || newCourse.Lecturer == "" {
// 			// 		res.WriteHeader(http.StatusUnprocessableEntity)
// 			// 		json.NewEncoder(res).Encode(ResMessage{ResponseText: "Sorry the course info is incomplete"})
// 			// 		return
// 			// 	}

// 			// 	var id string
// 			// 	// check course existence
// 			// 	rows, err := db.Query("SELECT course_id FROM "+os.Getenv("DB_SCHEMA")+".courses WHERE course_id = ?", newCourse.Course_id)
// 			// 	defer rows.Close()
// 			// 	if err != nil {
// 			// 		res.WriteHeader(http.StatusUnprocessableEntity)
// 			// 		json.NewEncoder(res).Encode(ResMessage{ResponseText: "Sorry the course cannot be added"})
// 			// 		return
// 			// 	}

// 			// 	for rows.Next() {
// 			// 		err := rows.Scan(&id)
// 			// 		if err != nil {
// 			// 			res.WriteHeader(http.StatusUnprocessableEntity)
// 			// 			json.NewEncoder(res).Encode(ResMessage{ResponseText: "Sorry the course cannot be added"})
// 			// 			return
// 			// 		}
// 			// 		break
// 			// 	}

// 			// 	if id == "" {
// 			// 		query := fmt.Sprintf("INSERT INTO "+os.Getenv("DB_SCHEMA")+".courses (course_id, title, description, lecturer, fee) VALUES ('%s', '%s', '%s', '%s', %.2f)",
// 			// 			newCourse.Course_id, newCourse.Title, newCourse.Description, newCourse.Lecturer, newCourse.Fee)
// 			// 		_, err := db.Exec(query)

// 			// 		if err != nil {
// 			// 			panic(err.Error())
// 			// 		}

// 			// 		res.WriteHeader(http.StatusCreated)
// 			// 		json.NewEncoder(res).Encode(ResMessage{ResponseText: "Course is created"})
// 			// 		return

// 			// 	} else {
// 			// 		res.WriteHeader(http.StatusConflict)
// 			// 		json.NewEncoder(res).Encode(ResMessage{ResponseText: "Sorry the course is already exist"})
// 			// 		return
// 			// 	}
// 			// } else {
// 			// 	res.WriteHeader(http.StatusUnprocessableEntity)
// 			// 	json.NewEncoder(res).Encode(ResMessage{ResponseText: "Sorry the course cannot be added"})
// 		}
// 	}
// }

func DeleteClassroom(res http.ResponseWriter, req *http.Request) {

}
