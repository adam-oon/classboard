package userclass

import (
	"classboard/config"
	"database/sql"
	"log"
	"testing"

	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
)

var db *sql.DB
var model UserClassModel
var mockLecturerId int
var mockStudentId int
var classroomId string
var isCreated bool

func Setup() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalln("Failed to load .env file", err)
	}

	var dbErr error
	db, dbErr = config.GetMySQLDB()
	if dbErr != nil {
		log.Fatalln(dbErr)
	}

	model = UserClassModel{
		Db: db,
	}

	// mock user data
	if !isCreated {
		var user_ids []int
		var mockUsers = []struct {
			username, password, usertype, name string
		}{
			{"userclasstest1", "1@2b3C4d", "lecturer", "userclasstest1 Lecturer"},
			{"userclasstest2", "1@2b3C4d", "student", "userclasstest2 Student"},
		}
		for _, mockUser := range mockUsers {
			res, err := model.Db.Exec("INSERT INTO users ( username, password, type, name) VALUES (?, ?, ?,?)",
				mockUser.username, mockUser.password, mockUser.usertype, mockUser.name)
			if err != nil {
				log.Fatalln(err)
			}
			user_id, _ := res.LastInsertId()
			user_ids = append(user_ids, int(user_id))
		}

		mockLecturerId = user_ids[0]
		mockStudentId = user_ids[1]

		// mock classroom data
		id, _ := uuid.NewV4()
		classroomId = id.String()
		var mockClassroom = struct {
			id, title, code string
			user_id         int
		}{
			classroomId, "classroomTitle", "CT100", mockLecturerId,
		}

		_, err := model.Db.Exec("INSERT INTO classrooms (id, user_id, code, title) VALUES (?, ?, ?, ?)",
			mockClassroom.id, mockClassroom.user_id, mockClassroom.code, mockClassroom.title)
		if err != nil {
			log.Fatalln(err)
		}

		isCreated = true
	}

}

func TearDown() {
	db.Close()
}

func TestJoinClass(t *testing.T) {
	Setup()
	defer TearDown()

	err := model.JoinClass(mockStudentId, classroomId)
	if err != nil {
		t.Error(err)
	}

	// check db through query
	rows, _ := model.Db.Query("SELECT * FROM user_classes WHERE classroom_id = ? AND user_id = ?", classroomId, mockStudentId)
	defer rows.Close()

	var result UserClass
	for rows.Next() {
		rows.Scan(&result.User_id, &result.Classroom_id)
	}

	if result.User_id != mockStudentId || result.Classroom_id != classroomId {
		t.Error("Join classroom failed")
	}
}

func TestGetClassroomStudent(t *testing.T) {
	Setup()
	defer TearDown()

	user_ids, err := model.GetClassroomStudent(classroomId)
	if err != nil {
		t.Error(err)
	}

	var isFound bool
	for _, user_id := range user_ids {
		if user_id == mockStudentId {
			isFound = true
		}
	}

	if !isFound {
		t.Errorf("GetClassroomStudent failed. Expecting %d in %v", mockStudentId, user_ids)
	}
}

func TestGetJoinedClass(t *testing.T) {
	Setup()
	defer TearDown()

	classrooms, err := model.GetJoinedClass(mockStudentId)
	if err != nil {
		t.Error(err)
	}

	var isFound bool
	for _, classroom := range classrooms {
		if classroom.Id == classroomId {
			isFound = true
		}
	}

	if !isFound {
		t.Errorf("GetJoinedClass failed. Expecting %s in %v", classroomId, classrooms)
	}
}
func TestIsBelongToClassroom(t *testing.T) {
	Setup()
	defer TearDown()

	var isBelong bool
	id, _ := uuid.NewV4()
	newClassroomId := id.String()
	var mockClassroom = struct {
		id, title, code string
		user_id         int
	}{
		newClassroomId, "newClassroomTitle", "NCT100", mockLecturerId,
	}

	_, err := model.Db.Exec("INSERT INTO classrooms (id, user_id, code, title) VALUES (?, ?, ?, ?)",
		mockClassroom.id, mockClassroom.user_id, mockClassroom.code, mockClassroom.title)
	if err != nil {
		log.Fatalln(err)
	}

	// before join the classroom
	isBelong = model.IsBelongToClassroom(mockStudentId, newClassroomId)
	if isBelong == true {
		t.Error("IsBelongToClassroom failed. Expecting false before join")
	}

	//join the class
	err = model.JoinClass(mockStudentId, newClassroomId)
	if err != nil {
		t.Error(err)
	}

	// after join classroom
	isBelong = model.IsBelongToClassroom(mockStudentId, newClassroomId)
	if isBelong == false {
		t.Error("IsBelongToClassroom failed. Expecting true after join")
	}
}
