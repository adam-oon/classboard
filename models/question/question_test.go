package question

import (
	"classboard/config"
	"database/sql"
	"log"
	"testing"

	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
)

var db *sql.DB
var model QuestionModel
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

	model = QuestionModel{
		Db: db,
	}

	// mock user data
	if !isCreated {
		var user_ids []int
		var mockUsers = []struct {
			username, password, usertype, name string
		}{
			{"questiontest1", "1@2b3C4d", "lecturer", "questiontest1 Lecturer"},
			{"questiontest2", "1@2b3C4d", "student", "questiontest2 Student"},
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

func TestSaveQuestion(t *testing.T) {
	Setup()
	defer TearDown()

	var multipleTest = []QuestionInput{
		{
			Classroom_id: classroomId,
			Question:     "1 + 2 = ?",
			Type:         "multiple",
			Choice:       "3|4|5",
			Solution:     "3",
		},
		{
			Classroom_id: classroomId,
			Question:     "Which of the following is not programming language?",
			Type:         "multiple",
			Choice:       "go|html|",
			Solution:     "html",
		},
	}

	for _, td := range multipleTest {
		t.Run("TestSaveQuestion", func(t *testing.T) {
			err := model.SaveQuestion(td)
			if err != nil {
				t.Errorf("Error when saving %+v\n", td)
			}
		})
	}
}

func TestGetQuestionsByClassroomId(t *testing.T) {
	Setup()
	defer TearDown()

	question := QuestionInput{
		Classroom_id: classroomId,
		Question:     "3 + 4 = ?",
		Type:         "multiple",
		Choice:       "3|4|5",
		Solution:     "3",
	}

	res, err := model.Db.Exec("INSERT INTO questions (classroom_id, question, type, choice,solution) VALUES (?,?,?,?,?)",
		question.Classroom_id, question.Question, question.Type, question.Choice, question.Solution)
	if err != nil {
		t.Error(err)
	}
	question_id, _ := res.LastInsertId()

	questions, err := model.GetQuestionsByClassroomId(classroomId)
	if err != nil {
		t.Error(err)
	}

	var isFound bool
	for _, question := range questions {
		if question.Id == int(question_id) {
			isFound = true
		}
	}

	if !isFound {
		t.Errorf("GetQuestionsByClassroomId failed. Expecting %d ID found in %v", question_id, questions)
	}
}

func TestGetQuestion(t *testing.T) {
	Setup()
	defer TearDown()

	question := QuestionInput{
		Classroom_id: classroomId,
		Question:     "3 + 4 = ?",
		Type:         "multiple",
		Choice:       "3|4|5",
		Solution:     "3",
	}

	res, err := model.Db.Exec("INSERT INTO questions (classroom_id, question, type, choice,solution) VALUES (?,?,?,?,?)",
		question.Classroom_id, question.Question, question.Type, question.Choice, question.Solution)
	if err != nil {
		t.Error(err)
	}
	question_id, _ := res.LastInsertId()

	result, err := model.GetQuestion(int(question_id))
	if err != nil {
		t.Error(err)
	}

	if result.Id != int(question_id) && result.Question != question.Question {
		t.Errorf("GetQuestion failed. Expecting %v match with result %v", question, result)
	}
}

func TestDeleteQuestion(t *testing.T) {
	Setup()
	defer TearDown()

	question := QuestionInput{
		Classroom_id: classroomId,
		Question:     "3 + 4 = ?",
		Type:         "multiple",
		Choice:       "3|4|5",
		Solution:     "3",
	}

	res, err := model.Db.Exec("INSERT INTO questions (classroom_id, question, type, choice,solution) VALUES (?,?,?,?,?)",
		question.Classroom_id, question.Question, question.Type, question.Choice, question.Solution)
	if err != nil {
		t.Error(err)
	}
	question_id, _ := res.LastInsertId()

	err = model.DeleteQuestion(int(question_id))
	if err != nil {
		t.Error(err)
	}
}
