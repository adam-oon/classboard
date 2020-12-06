package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type messageLog struct {
	Message string
}

type User struct {
	Username        string
	Type            string
	Password        string
	ConfirmPassword string
	Name            string
}

type ResMessage struct {
	ResponseText string
	ID           string
}

func sanitizeUserInput(user *User) {
	user.Username = strings.TrimSpace(user.Username)
	user.Type = strings.TrimSpace(user.Type)
	user.Name = strings.TrimSpace(user.Name)
}

func register(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	// if alreadyLoggedIn(req) {
	// 	http.Redirect(res, req, "/", http.StatusSeeOther)
	// 	return
	// }

	// var newUser User
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

			if newUser.Password != newUser.ConfirmPassword {
				res.WriteHeader(http.StatusForbidden)
				json.NewEncoder(res).Encode(ResMessage{ResponseText: "Passwords are not match"})
				return
			}

			var id string
			// check course existence
			rows, err := db.Query("SELECT COUNT(UPPERProductID) FROM "+os.Getenv("DB_SCHEMA")+".courses WHERE course_id = ?", newUser.Course_id)
			defer rows.Close()
			if err != nil {
				res.WriteHeader(http.StatusUnprocessableEntity)
				json.NewEncoder(res).Encode(ResMessage{ResponseText: "Sorry the course cannot be added"})
				return
			}

			for rows.Next() {
				err := rows.Scan(&id)
				if err != nil {
					res.WriteHeader(http.StatusUnprocessableEntity)
					json.NewEncoder(res).Encode(ResMessage{ResponseText: "Sorry the course cannot be added"})
					return
				}
				break
			}

			if id == "" {
				query := fmt.Sprintf("INSERT INTO "+os.Getenv("DB_SCHEMA")+".courses (course_id, title, description, lecturer, fee) VALUES ('%s', '%s', '%s', '%s', %.2f)",
					newUser.Course_id, newUser.Title, newUser.Description, newUser.Lecturer, newUser.Fee)
				_, err := db.Exec(query)

				if err != nil {
					panic(err.Error())
				}

				res.WriteHeader(http.StatusCreated)
				json.NewEncoder(res).Encode(ResMessage{ResponseText: "Course is created"})
				return

			} else {
				res.WriteHeader(http.StatusConflict)
				json.NewEncoder(res).Encode(ResMessage{ResponseText: "Sorry the course is already exist"})
				return
			}
		} else {
			res.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(res).Encode(ResMessage{ResponseText: "Sorry the course cannot be added"})
		}

		if username != "" {
			// check if username exist/ taken
			if _, exist := mapUsers[username]; exist {
				res.WriteHeader(http.StatusForbidden)
				mLog = messageLog{"Username already taken. Please choose another username"}
				fatalErr := tpl.ExecuteTemplate(res, "register.gohtml", mLog)
				if fatalErr != nil {
					log.Fatalln(fatalErr)
				}
				return
			}
			// create session
			id, _ := uuid.NewV4()
			myCookie := &http.Cookie{
				Name:  "myCookie",
				Value: id.String(),
			}
			http.SetCookie(res, myCookie)
			mapSessions[myCookie.Value] = username
			bPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
			if err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				mLog = messageLog{"Internal Server Error. Please contact system administrator!"}
				fatalErr := tpl.ExecuteTemplate(res, "register.gohtml", mLog)
				if fatalErr != nil {
					log.Fatalln(fatalErr)
				}
				return
			}

			// store new user info
			query := fmt.Sprintf("INSERT INTO "+os.Getenv("DB_SCHEMA")+".users (id, username, password, type, name) VALUES ('%s', '%s', '%s', '%s','%s')",
				newUser.ID, newUser.Username, newUser.Password, newUser.Type, newUser.Name)
			_, err2 := db.Exec(query)

			if err2 != nil {
				panic(err.Error())
			}

			res.WriteHeader(http.StatusCreated)
			json.NewEncoder(res).Encode(ResMessage{ResponseText: "Course is created"})
			return
		}

		// redirect to main index
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	fatalErr := tpl.ExecuteTemplate(res, "register.gohtml", mLog)
	if fatalErr != nil {
		log.Fatalln(fatalErr)
	}
}
