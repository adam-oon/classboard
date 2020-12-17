package main

import (
	"classboard/helper"
	sessionmodel "classboard/models/session"
	usermodel "classboard/models/user"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRegistration struct {
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

/*
	http handlers
*/
func registerHandler(res http.ResponseWriter, req *http.Request) {
	if isLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	res.Header().Set("Content-Type", "application/json")

	userModel := usermodel.UserModel{
		Db: db,
	}

	// process form submission
	if req.Method == http.MethodPost { //POST
		// get form values
		var newUser UserRegistration
		reqBody, err := ioutil.ReadAll(req.Body)

		if err == nil {
			err := json.Unmarshal(reqBody, &newUser)
			if err != nil {
				res.WriteHeader(http.StatusUnprocessableEntity)
				json.NewEncoder(res).Encode(ResMessage{ResponseText: "Sorry the user info is incomplete"})
				return
			}

			// input sanitization & validation
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

			// check user existence
			var count int
			count, err = userModel.CheckUserByUsername(newUser.Username)

			if count == 0 {
				bPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.MinCost)
				if err != nil {
					res.WriteHeader(http.StatusInternalServerError)
					json.NewEncoder(res).Encode(ResMessage{ResponseText: "Internal Server Error. Please contact system administrator!"})
					return
				}
				newUser.Password = string(bPassword)

				// store new user info
				err = userModel.SaveUser(newUser.Username, newUser.Password, newUser.Type, newUser.Name)
				if err != nil {
					Warning.Println(err)
					res.WriteHeader(http.StatusUnprocessableEntity)
					json.NewEncoder(res).Encode(ResMessage{ResponseText: "Sorry the user cannot be added"})
					return
				}

				res.WriteHeader(http.StatusCreated)
				json.NewEncoder(res).Encode(ResMessage{ResponseText: "Register Successful!"})

			} else {
				res.WriteHeader(http.StatusConflict)
				json.NewEncoder(res).Encode(ResMessage{ResponseText: "Sorry the username is already taken"})
				return
			}
		} else {
			Warning.Println(err)
			res.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(res).Encode(ResMessage{ResponseText: "Sorry the user cannot be added"})
			return
		}
	}
}

func loginHandler(res http.ResponseWriter, req *http.Request) {
	if isLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

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

		userModel := usermodel.UserModel{
			Db: db,
		}
		user, err := userModel.GetUserByUsername(userLogin.Username)
		if err != nil {
			Info.Println(err)
		}
		if user == (usermodel.User{}) {
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
		sessionModel := sessionmodel.SessionModel{
			Db: db,
		}
		sessionModel.DeleteSessionByUserId(user.Id)

		// create session
		id, _ := uuid.NewV4()
		myCookie := &http.Cookie{
			Name:     "myCookie",
			Value:    id.String(),
			Expires:  time.Now().Add(2 * time.Hour), // 2 hrs
			HttpOnly: true,                          // cannot be used by javascript
			Secure:   true,                          // send through HTTPS only
			Domain:   "localhost",
		}
		http.SetCookie(res, myCookie)

		session := sessionmodel.Session{
			Session_id: myCookie.Value,
			User_id:    user.Id,
		}
		err = sessionModel.SaveSession(session)
		if err != nil {
			Warning.Println(err)
			res.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(res).Encode(ResMessage{ResponseText: "Username and/or password do not match"})
			return
		}

		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(ResMessage{ResponseText: "Welcome User!"})
	}
}

func logoutHandler(res http.ResponseWriter, req *http.Request) {
	if !isLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	myCookie, _ := req.Cookie("myCookie")
	// delete the session from DB
	sessionModel := sessionmodel.SessionModel{
		Db: db,
	}
	sessionModel.DeleteSessionBySessionId(myCookie.Value)
	// remove the cookie
	myCookie = &http.Cookie{
		Name:   "myCookie",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(res, myCookie)
	http.Redirect(res, req, "/", http.StatusSeeOther)
}

/*
	auth functions
*/
func getSessionUser(req *http.Request) usermodel.User {
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		Error.Println(err)
	}
	userModel := usermodel.UserModel{
		Db: db,
	}
	sessionModel := sessionmodel.SessionModel{
		Db: db,
	}
	user_id, err := sessionModel.GetUserID(myCookie.Value)
	if err != nil {
		Info.Println(err)
	}
	user, err := userModel.GetUser(user_id)
	if err != nil {
		Info.Println(err)
	}
	return user
}

func isLecturer(user_type string) bool {
	return user_type == "lecturer"
}

func isStudent(user_type string) bool {
	return user_type == "student"
}

func isLoggedIn(req *http.Request) bool {
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		return false
	}

	sessionModel := sessionmodel.SessionModel{
		Db: db,
	}
	ok := sessionModel.CheckSession(myCookie.Value)
	return ok
}

func sanitizeUserInput(user *UserRegistration) {
	user.Username = strings.TrimSpace(user.Username)
	user.Type = strings.TrimSpace(user.Type)
	user.Name = strings.TrimSpace(user.Name)
}
