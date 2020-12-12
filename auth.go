package main

import (
	"classboard/models"
	"classboard/pkg/helper"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username        string
	Type            string
	Password        string
	ConfirmPassword string
	Name            string
}

type UserLogin struct {
	Username string
	Password string
}

type ResMessage struct {
	ResponseText string
	ID           string
}

type messageLog struct {
	Message string
}

func register(res http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	res.Header().Set("Content-Type", "application/json")

	var mLog messageLog
	// process form submission
	if req.Method == http.MethodPost { //POST
		// get form values

		var newUser User
		reqBody, err := ioutil.ReadAll(req.Body)

		if err == nil {
			err := json.Unmarshal(reqBody, &newUser)
			if err != nil {
				res.WriteHeader(http.StatusUnprocessableEntity)
				json.NewEncoder(res).Encode(ResMessage{ResponseText: "Sorry the user info is incomplete"})
				return
			}

			// input validation
			sanitizeUserInput(&newUser)
			if newUser.Username == "" || newUser.Type == "" || newUser.Password == "" || newUser.ConfirmPassword == "" || newUser.Name == "" {
				res.WriteHeader(http.StatusUnprocessableEntity)
				json.NewEncoder(res).Encode(ResMessage{ResponseText: "Sorry the user info is incomplete"})
				return
			}

			err = helper.CheckPasswordStrength(newUser.Password)
			if err != nil {
				res.WriteHeader(http.StatusUnprocessableEntity)
				json.NewEncoder(res).Encode(ResMessage{ResponseText: err.Error()})
				return
			}

			if newUser.Password != newUser.ConfirmPassword {
				res.WriteHeader(http.StatusUnprocessableEntity)
				json.NewEncoder(res).Encode(ResMessage{ResponseText: "Passwords are not match"})
				return
			}

			var count int
			// check user existence
			rows, err := db.Query("SELECT COUNT(username) as totalUsername FROM users WHERE username = ?", newUser.Username)
			defer rows.Close()
			if err != nil {
				res.WriteHeader(http.StatusUnprocessableEntity)
				json.NewEncoder(res).Encode(ResMessage{ResponseText: "Sorry the user cannot be added"})
				return
			}

			for rows.Next() {
				err := rows.Scan(&count)
				if err != nil {
					res.WriteHeader(http.StatusUnprocessableEntity)
					json.NewEncoder(res).Encode(ResMessage{ResponseText: "Sorry the user has been taken. Please choose another one"})
					return
				}
				break
			}

			if count == 0 {
				bPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.MinCost)
				if err != nil {
					res.WriteHeader(http.StatusInternalServerError)
					mLog = messageLog{"Internal Server Error. Please contact system administrator!"}
					fatalErr := tpl.ExecuteTemplate(res, "register.gohtml", mLog)
					if fatalErr != nil {
						log.Fatalln(fatalErr)
					}
					return
				}

				newUser.Password = string(bPassword)

				var errExec error
				// store new user info
				_, errExec = db.Exec("INSERT INTO users ( username, password, type, name) VALUES (?, ?, ?,?)", newUser.Username, newUser.Password, newUser.Type, newUser.Name)

				if errExec != nil {
					panic(errExec.Error())
				}

				res.WriteHeader(http.StatusCreated)
				json.NewEncoder(res).Encode(ResMessage{ResponseText: "Register Successful!"})

			} else {
				res.WriteHeader(http.StatusConflict)
				json.NewEncoder(res).Encode(ResMessage{ResponseText: "Sorry the username is already taken"})
				return

			}
		} else {
			res.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(res).Encode(ResMessage{ResponseText: "Sorry the user cannot be added"})
			return
		}

	}
}

func login(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	var userLogin UserLogin
	reqBody, err := ioutil.ReadAll(req.Body)

	if err == nil {
		err := json.Unmarshal(reqBody, &userLogin)
		if err != nil {
			res.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(res).Encode(ResMessage{ResponseText: "Username and/or password do not match"})
			return
		}

		user := models.GetUserByUsername(userLogin.Username)
		if user == (models.User{}) {
			res.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(res).Encode(ResMessage{ResponseText: "Username and/or password do not match"})
			return
		}

		// verify user
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userLogin.Password))
		if err != nil {
			res.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(res).Encode(ResMessage{ResponseText: "Username and/or password do not match"})
			return
		}

		// check user session and delete it
		models.DeleteSession(user.Id)

		// create session
		id, _ := uuid.NewV4()
		myCookie := &http.Cookie{
			Name:     "myCookie",
			Value:    id.String(),
			Expires:  time.Now().Add(2 * time.Hour),
			HttpOnly: true,
			Domain:   "localhost",
			Secure:   true,
		}
		http.SetCookie(res, myCookie)

		_, errExec := db.Exec("INSERT INTO sessions (session_id, user_id) VALUES (?,?)", myCookie.Value, user.Id)
		if errExec != nil {
			panic(errExec.Error())
		}
		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(ResMessage{ResponseText: "Course doesn't exist"})
	}
}

func logout(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	myCookie, _ := req.Cookie("myCookie")
	// delete the session
	models.DeleteSessionByID(myCookie.Value)
	// remove the cookie
	myCookie = &http.Cookie{
		Name:   "myCookie",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(res, myCookie)
	http.Redirect(res, req, "/", http.StatusSeeOther)
}

func sanitizeUserInput(user *User) {
	user.Username = strings.TrimSpace(user.Username)
	user.Type = strings.TrimSpace(user.Type)
	user.Name = strings.TrimSpace(user.Name)
}

func getUser(req *http.Request) models.User {
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		//error//
	}

	user_id := models.GetUserID(myCookie.Value)
	user := models.GetUser(user_id)
	return user
}

func isLecturer(user_type string) bool {
	return user_type == "lecturer"
}

func isStudent(user_type string) bool {
	return user_type == "student"
}

func alreadyLoggedIn(req *http.Request) bool {
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		return false
	}

	ok := models.CheckSession(myCookie.Value)
	return ok
}
