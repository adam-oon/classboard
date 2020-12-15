package session

import (
	"classboard/config"
	"database/sql"
	"log"
	"testing"

	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
)

var db *sql.DB
var model SessionModel
var mockUserId int
var sessionId string
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

	model = SessionModel{
		Db: db,
	}

	// mock user data
	if !isCreated {
		var mockUser = struct {
			username, password, usertype, name string
		}{
			"sessiontest", "1@2b3C4d", "lecturer", "Session Test",
		}
		res, err := model.Db.Exec("INSERT INTO users ( username, password, type, name) VALUES (?, ?, ?,?)",
			mockUser.username, mockUser.password, mockUser.usertype, mockUser.name)
		if err != nil {
			log.Fatalln(err)
		}
		user_id, _ := res.LastInsertId()
		mockUserId = int(user_id)

		// mock session id

		isCreated = true
	}

	id, _ := uuid.NewV4()
	sessionId = id.String()
	// clear session table
	model.Db.Exec("DELETE FROM sessions WHERE user_id =?", mockUserId)
}

func TearDown() {
	db.Close()
}

func TestSaveSession(t *testing.T) {
	Setup()
	defer TearDown()

	session := Session{
		Session_id: sessionId,
		User_id:    mockUserId,
	}

	err := model.SaveSession(session)
	if err != nil {
		t.Errorf("Error when saving %+v\n", session)
	}

	// check db through query
	rows, _ := model.Db.Query("SELECT * FROM sessions WHERE session_id = ? and user_id = ?", session.Session_id, session.User_id)
	defer rows.Close()

	var sessionResult Session
	for rows.Next() {
		rows.Scan(&sessionResult.Session_id, &sessionResult.User_id)
	}

	if sessionResult.Session_id != session.Session_id || sessionResult.User_id != session.User_id {
		t.Error("Session save failed")
	}
}

func TestCheckSession(t *testing.T) {
	Setup()
	defer TearDown()

	session := Session{
		Session_id: sessionId,
		User_id:    mockUserId,
	}

	_, err := model.Db.Exec("INSERT INTO sessions ( session_id, user_id) VALUES (?,?)",
		session.Session_id, session.User_id)
	if err != nil {
		t.Error(err)
	}

	exist := model.CheckSession(sessionId)
	if !exist {
		t.Errorf("Should be exist after insert session data\n")
	}
}

func TestGetUserID(t *testing.T) {
	Setup()
	defer TearDown()

	user_id, err := model.GetUserID(sessionId)
	if err != nil {
		t.Error(err)
	}

	if user_id != 0 {
		t.Error("Should be 0 if no session found\n")
	}

	session := Session{
		Session_id: sessionId,
		User_id:    mockUserId,
	}

	_, err = model.Db.Exec("INSERT INTO sessions ( session_id, user_id) VALUES (?,?)",
		session.Session_id, session.User_id)
	if err != nil {
		t.Error(err)
	}

	user_id, err = model.GetUserID(sessionId)
	if user_id != session.User_id {
		t.Errorf("Returned ID %d is not equal to Expected ID %d\n", user_id, session.User_id)
	}
}

func TestDeleteSessionByUserId(t *testing.T) {
	Setup()
	defer TearDown()

	session := Session{
		Session_id: sessionId,
		User_id:    mockUserId,
	}

	_, err := model.Db.Exec("INSERT INTO sessions ( session_id, user_id) VALUES (?,?)",
		session.Session_id, session.User_id)
	if err != nil {
		t.Error(err)
	}

	err = model.DeleteSessionByUserId(mockUserId)
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteSessionBySessionId(t *testing.T) {
	Setup()
	defer TearDown()

	session := Session{
		Session_id: sessionId,
		User_id:    mockUserId,
	}

	_, err := model.Db.Exec("INSERT INTO sessions ( session_id, user_id) VALUES (?,?)",
		session.Session_id, session.User_id)
	if err != nil {
		t.Error(err)
	}

	err = model.DeleteSessionBySessionId(sessionId)
	if err != nil {
		t.Error(err)
	}
}
