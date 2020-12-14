package classroom

import (
	"classboard/config"
	"database/sql"
	"log"
	"testing"

	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
)

var db *sql.DB
var model ClassroomModel
var mockUserId int
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

	model = ClassroomModel{
		Db: db,
	}

	// mock user data
	if !isCreated {
		var mockUser = struct {
			username, password, usertype, name string
		}{
			"classroomtest1", "1@2b3C4d", "lecturer", "Classroom Test",
		}
		res, err := model.Db.Exec("INSERT INTO users ( username, password, type, name) VALUES (?, ?, ?,?)",
			mockUser.username, mockUser.password, mockUser.usertype, mockUser.name)
		if err != nil {
			log.Fatalln(err)
		}
		user_id, _ := res.LastInsertId()
		mockUserId = int(user_id)

		// mock classroom id
		id, _ := uuid.NewV4()
		classroomId = id.String()

		isCreated = true
	}

}

func TearDown() {
	db.Close()
}

func TestSaveClassroom(t *testing.T) {
	Setup()
	defer TearDown()

	var (
		newTitle string = "Save Classroom Test"
		newCode  string = "SCT100"
	)

	classroom := Classroom{
		Id:      classroomId,
		User_id: mockUserId,
		Title:   "Save Classroom Test",
		Code:    "SCT100",
	}

	err := model.SaveClassroom(classroom)
	if err != nil {
		t.Errorf("Error when saving %+v\n", classroom)
	}

	// check db through query
	rows, _ := model.Db.Query("SELECT * FROM classrooms WHERE id = ?", classroomId)
	defer rows.Close()

	var classroomResult Classroom
	for rows.Next() {
		rows.Scan(&classroomResult.Id, &classroomResult.User_id, &classroomResult.Code, &classroomResult.Title)
	}

	if classroomResult.Code != newCode || classroomResult.Title != newTitle {
		t.Error("Classroom save failed")
	}
}

func TestUpdateClassroom(t *testing.T) {
	Setup()
	defer TearDown()

	var (
		newTitle string = "Save Classroom Test1"
		newCode  string = "SCT101"
	)
	classroom := Classroom{
		Id:      classroomId,
		User_id: mockUserId,
		Title:   newTitle,
		Code:    newCode,
	}

	err := model.UpdateClassroom(classroom)
	if err != nil {
		t.Errorf("Error when saving %+v\n", classroom)
	}

	// check db through query
	rows, _ := model.Db.Query("SELECT * FROM classrooms WHERE id = ?", classroomId)
	defer rows.Close()

	var classroomResult Classroom
	for rows.Next() {
		rows.Scan(&classroomResult.Id, &classroomResult.User_id, &classroomResult.Code, &classroomResult.Title)
	}

	if classroomResult.Code != newCode || classroomResult.Title != newTitle {
		t.Error("Classroom update failed")
	}
}

func TestGetClassroom(t *testing.T) {
	Setup()
	defer TearDown()

	classroom, err := model.GetClassroom(classroomId)
	if err != nil {
		t.Error(err)
	}

	if mockUserId != classroom.User_id {
		t.Errorf("Get incorrect classroom %+v with %s\n", classroom, classroomId)
	}
}

func TestGetClassroomsByUserId(t *testing.T) {
	Setup()
	defer TearDown()

	classrooms, err := model.GetClassroomsByUserId(mockUserId)
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
		t.Errorf("Get incorrect classroom %+v with user_id :%s\n", classrooms, classroomId)
	}
}
