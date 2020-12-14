package user

import (
	"classboard/config"
	"database/sql"
	"log"
	"testing"

	"github.com/joho/godotenv"
)

var db *sql.DB
var model UserModel

func Setup() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalln("Failed to load .env file", err)
	}

	var dbErr error
	db, dbErr = config.GetMySQLDB()
	if dbErr != nil {
		log.Fatalln(dbErr)
	}

	model = UserModel{
		Db: db,
	}
}

func TearDown() {
	db.Close()
}

func TestSaveUser(t *testing.T) {
	Setup()
	defer TearDown()

	var multipleTest = []struct {
		username, password, usertype, name string
	}{
		{"billgates", "1@2b3C4d", "lecturer", "Bill Gates"},
		{"adam91", "5@6b7C8d", "student", "Adam Oon"},
	}

	for _, td := range multipleTest {
		t.Run("TestSaveUser", func(t *testing.T) {
			err := model.SaveUser(td.username, td.password, td.usertype, td.name)
			if err != nil {
				t.Errorf("Error when saving %+v\n", td)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	Setup()
	defer TearDown()

	var mockUser = struct {
		username, password, usertype, name string
	}{
		"billgates1", "1@2b3C4d", "lecturer", "Bill Gates",
	}
	res, err := model.Db.Exec("INSERT INTO users ( username, password, type, name) VALUES (?, ?, ?,?)",
		mockUser.username, mockUser.password, mockUser.usertype, mockUser.name)
	if err != nil {
		t.Error(err)
	}
	user_id, _ := res.LastInsertId()

	user, err := model.GetUser(int(user_id))
	if err != nil {
		t.Error(err)
	}

	if mockUser.username != user.Username {
		t.Errorf("Get incorrect user %s and %s\n", mockUser.username, user.Username)
	}
}

func TestGetUserByUsername(t *testing.T) {
	Setup()
	defer TearDown()

	var mockUser = struct {
		username, password, usertype, name string
	}{
		"billgates2", "1@2b3C4d", "lecturer", "Bill Gates",
	}
	_, err := model.Db.Exec("INSERT INTO users ( username, password, type, name) VALUES (?, ?, ?,?)",
		mockUser.username, mockUser.password, mockUser.usertype, mockUser.name)
	if err != nil {
		t.Error(err)
	}

	user, err := model.GetUserByUsername(mockUser.username)
	if err != nil {
		t.Error(err)
	}
	if mockUser.username != user.Username {
		t.Errorf("Get incorrect user %s and %s\n", mockUser.username, user.Username)
	}
}

func TestCheckUserByUsername(t *testing.T) {
	Setup()
	defer TearDown()

	// get 0 if no username found
	count, err := model.CheckUserByUsername("mock User.username")
	if err != nil {
		t.Error(err)
	}

	if count != 0 {
		t.Error("Should be 0 if no username found")
	}

	// get 1 if username found
	var mockUser = struct {
		username, password, usertype, name string
	}{
		"billgates4", "1@2b3C4d", "lecturer", "Bill Gates",
	}
	_, err = model.Db.Exec("INSERT INTO users ( username, password, type, name) VALUES (?, ?, ?,?)",
		mockUser.username, mockUser.password, mockUser.usertype, mockUser.name)
	if err != nil {
		t.Error(err)
	}

	count, err = model.CheckUserByUsername(mockUser.username)
	if err != nil {
		t.Error(err)
	}

	if count != 1 || count > 1 {
		t.Errorf("Should be 1 if username found")
	}
}
