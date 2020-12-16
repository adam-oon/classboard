package answer

import (
	"classboard/config"
	"database/sql"
	"log"
	"testing"

	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
)

var db *sql.DB
var model AnswerModel
var mockLecturerId int
var mockStudentId int
var mockQusetionId int
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

	model = AnswerModel{
		Db: db,
	}

	// mock user data
	if !isCreated {
		var user_ids []int
		var mockUsers = []struct {
			username, password, usertype, name string
		}{
			{"answertest1", "1@2b3C4d", "lecturer", "answertest1 Lecturer"},
			{"answertest2", "1@2b3C4d", "student", "answertest2 Student"},
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
			classroomId, "answerTestTitle", "AT100", mockLecturerId,
		}

		_, err := model.Db.Exec("INSERT INTO classrooms (id, user_id, code, title) VALUES (?, ?, ?, ?)",
			mockClassroom.id, mockClassroom.user_id, mockClassroom.code, mockClassroom.title)
		if err != nil {
			log.Fatalln(err)
		}

		// mock student join classroom
		_, err = model.Db.Exec("INSERT INTO user_classes (user_id, classroom_id) VALUES (?,?)", mockStudentId, classroomId)
		if err != nil {
			log.Fatalln(err)
		}

		// mock question data
		var mockQuestion = struct {
			classroom_id, question, question_type, choice, solution string
		}{
			classroomId, "1 + 2 = ?", "multiple", "3|4|5", "3",
		}
		res, err := model.Db.Exec("INSERT INTO questions (classroom_id, question, type,choice, solution) VALUES (?,?,?,?,?)",
			mockQuestion.classroom_id, mockQuestion.question, mockQuestion.question_type, mockQuestion.choice, mockQuestion.solution)
		if err != nil {
			log.Fatalln(err)
		}
		latestId, _ := res.LastInsertId()
		mockQusetionId = int(latestId)

		isCreated = true
	}

}

func TearDown() {
	db.Close()
}

func TestSaveAnswer(t *testing.T) {
	Setup()
	defer TearDown()

	var answer = Answer{
		Question_id: mockQusetionId,
		User_id:     mockStudentId,
		Answer:      "3",
		Is_correct:  true,
	}

	err := model.SaveAnswer(answer)
	if err != nil {
		t.Errorf("Error when saving %+v : %s\n", answer, err.Error())
	}
}

func TestGetAnswer(t *testing.T) {
	Setup()
	defer TearDown()

	var expectedAns = Answer{
		Question_id: mockQusetionId,
		User_id:     mockStudentId,
		Answer:      "3",
		Is_correct:  true,
	}

	ans, err := model.GetAnswer(mockQusetionId, mockStudentId)
	if err != nil {
		t.Error(err)
	}
	if ans.Question_id != expectedAns.Question_id || ans.User_id != expectedAns.User_id || ans.Answer != expectedAns.Answer || ans.Is_correct != expectedAns.Is_correct {
		t.Errorf("GetAnswer(%d,%d) failed. Expected %v match with Result %v\n", mockQusetionId, mockStudentId, expectedAns, ans)
	}
}

func TestDeleteAnswer(t *testing.T) {
	Setup()
	defer TearDown()

	var expectedAns = Answer{
		Question_id: mockQusetionId,
		User_id:     mockStudentId,
		Answer:      "3",
		Is_correct:  true,
	}

	err := model.DeleteAnswer(expectedAns.Question_id)
	if err != nil {
		t.Error(err)
	}
}
