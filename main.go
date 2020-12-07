package main

import (
	"classboard/models"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var tpl *template.Template
var db *sql.DB

const (
	APPURL string = "http://localhost:8080"
	APIURL string = "http://localhost:8080/api/v1"
)

func init() {
	tpl = template.Must(template.ParseGlob("views/*.gohtml"))

	// initialize env file
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalln("Failed to load .env file", err)
	}

	// for docker's database schema & table creation
	// connection := os.Getenv("DB_ACCOUNT") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_IP") + ":" + os.Getenv("DB_PORT") + ")/"
	// db, dbErr := sql.Open("mysql", connection)
	// defer db.Close()
	// if dbErr != nil {
	// 	panic(dbErr.Error())
	// }

	// _, err := db.Exec("CREATE DATABASE " + os.Getenv("DB_SCHEMA"))
	// if driverErr, ok := err.(*mysql.MySQLError); ok {
	// 	if driverErr.Number != 1007 { // skip if schema already exist
	// 		panic(err.Error())
	// 	}
	// } else {
	// 	log.Println("Database schema created successfully")
	// }

	// _, err = db.Exec("USE " + os.Getenv("DB_SCHEMA"))
	// if err != nil {
	// 	panic(err.Error())
	// }

	// createQuery := "CREATE TABLE `" + os.Getenv("DB_SCHEMA") + "`.`courses` (`course_id` VARCHAR(255) NOT NULL, `title` VARCHAR(255) NOT NULL, `description` MEDIUMTEXT NOT NULL, `lecturer` VARCHAR(255) NOT NULL, `fee` DECIMAL(12,2) NOT NULL, PRIMARY KEY (`course_id`));"
	// _, err = db.Exec(createQuery)
	// if driverErr, ok := err.(*mysql.MySQLError); ok {
	// 	if driverErr.Number != 1050 { // skip if table already exist
	// 		panic(err.Error())
	// 	}
	// } else {
	// 	log.Println("Database table created successfully")
	// }
}

func main() {
	router := mux.NewRouter()
	// page
	router.HandleFunc("/", indexPage)
	router.HandleFunc("/register", registerPage).Methods("GET")
	router.HandleFunc("/register", register).Methods("POST")
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/logout", logout)
	router.HandleFunc("/dashboard", dashboardPage)
	// router.HandleFunc("/add", addCoursePage)
	// router.HandleFunc("/update/{courseid}", updateCoursePage)
	// router.HandleFunc("/get/{courseid}", getCoursePage)
	// router.HandleFunc("/delete", deleteCoursePage)
	// router.HandleFunc("/courses", coursesPage)

	// // api
	// router.HandleFunc("/api/v1/courses", getCourses).Methods("GET")
	// router.HandleFunc("/api/v1/courses/{courseid}", course).Methods("GET", "PUT", "POST", "DELETE")

	var dbErr error
	connection := os.Getenv("DB_ACCOUNT") + ":" + os.Getenv("DB_PASSWORD") + "@tcp(" + os.Getenv("DB_IP") + ":" + os.Getenv("DB_PORT") + ")/" + os.Getenv("DB_SCHEMA")
	fmt.Println(connection)
	db, dbErr = sql.Open("mysql", connection)
	defer db.Close()
	if dbErr != nil {
		panic(dbErr.Error())
	}

	// if fatalErr := http.ListenAndServe(":8080", router); fatalErr != nil {
	// 	log.Fatalln(fatalErr)
	// }

	if fatalErr := http.ListenAndServeTLS(":8080", os.Getenv("CERTIFICATE"), os.Getenv("PRIVATE_KEY"), router); fatalErr != nil {
		log.Fatalln(fatalErr)
	}
}

func alreadyLoggedIn(req *http.Request) bool {
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		return false
	}
	sessionModel := models.SessionModel{
		Db: db,
	}
	ok := sessionModel.CheckSession(myCookie.Value)
	return ok
}

func indexPage(res http.ResponseWriter, req *http.Request) {
	// return to dashboard if login
	if alreadyLoggedIn(req) {
		http.Redirect(res, req, "/dashboard", http.StatusSeeOther)
		return
	}
	fatalErr := tpl.ExecuteTemplate(res, "index.gohtml", nil)
	if fatalErr != nil {
		log.Println(fatalErr)
	}
}

func registerPage(res http.ResponseWriter, req *http.Request) {
	if alreadyLoggedIn(req) {
		http.Redirect(res, req, "/dashboard", http.StatusSeeOther)
		return
	}
	fatalErr := tpl.ExecuteTemplate(res, "register.gohtml", nil)
	if fatalErr != nil {
		log.Println(fatalErr)
	}
}

func dashboardPage(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	fatalErr := tpl.ExecuteTemplate(res, "dashboard.gohtml", nil)
	if fatalErr != nil {
		log.Println(fatalErr)
	}
}

// func loginPage(res http.ResponseWriter, req *http.Request) {
// 	Trace.Println("Login Page")
// 	var mLog messageLog
// 	if alreadyLoggedIn(req) {
// 		http.Redirect(res, req, "/", http.StatusSeeOther)
// 		return
// 	}

// 	// security : constant response time return
// 	time.Sleep(3 * time.Second)

// 	if req.Method == http.MethodPost { //POST
// 		username := strings.TrimSpace(req.FormValue("username"))
// 		password := req.FormValue("password")
// 		// check if user exist with username
// 		myUser, exist := mapUsers[username]
// 		if !exist {
// 			res.WriteHeader(http.StatusForbidden)
// 			Info.Println("Username and/or password do not match")
// 			mLog = messageLog{"Username and/or password do not match"}
// 			fatalErr := tpl.ExecuteTemplate(res, "index.gohtml", mLog)
// 			if fatalErr != nil {
// 				Warning.Println(fatalErr)
// 			}
// 			return
// 		}
// 		// Matching of password entered
// 		err := bcrypt.CompareHashAndPassword(myUser.Password, []byte(password))
// 		if err != nil {
// 			res.WriteHeader(http.StatusForbidden)
// 			Info.Println("Username and/or password do not match")
// 			mLog = messageLog{"Username and/or password do not match"}
// 			fatalErr := tpl.ExecuteTemplate(res, "index.gohtml", mLog)
// 			if fatalErr != nil {
// 				Warning.Println(fatalErr)
// 			}
// 			return
// 		}

// 		expireTime := time.Now().Add(2 * time.Hour)
// 		// Create and sign the JWT
// 		claims := &Claims{
// 			Username: username,
// 			StandardClaims: jwt.StandardClaims{
// 				ExpiresAt: expireTime.Unix(),
// 			},
// 		}
// 		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 		tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
// 		if err != nil {
// 			Warning.Println("Error during token creation")
// 		}
// 		myCookie := &http.Cookie{
// 			Name:     "myCookie",
// 			Value:    tokenString,
// 			Expires:  expireTime,
// 			HttpOnly: true,
// 			Domain:   "localhost",
// 			Secure:   true,
// 		}
// 		http.SetCookie(res, myCookie)
// 		mapSessions[myCookie.Value] = username
// 		Trace.Println(username + " is login")
// 		http.Redirect(res, req, "/", http.StatusSeeOther)
// 		return
// 	}
// 	fatalErr := tpl.ExecuteTemplate(res, "index.gohtml", nil)
// 	if fatalErr != nil {
// 		Warning.Println(fatalErr)
// 	}
// }

// func logoutPage(res http.ResponseWriter, req *http.Request) {
// 	Trace.Println("Logout Page")
// 	if !alreadyLoggedIn(req) {
// 		http.Redirect(res, req, "/", http.StatusSeeOther)
// 		return
// 	}
// 	myCookie, err := req.Cookie("myCookie")
// 	username := mapSessions[myCookie.Value]
// 	if err != nil {
// 		Info.Println("cookie cannot be found")
// 	}
// 	// delete the session
// 	delete(mapSessions, myCookie.Value)
// 	// remove the cookie
// 	myCookie = &http.Cookie{
// 		Name:   "myCookie",
// 		Value:  "",
// 		MaxAge: -1,
// 	}
// 	http.SetCookie(res, myCookie)
// 	Trace.Println(username + " is logout")
// 	http.Redirect(res, req, "/", http.StatusSeeOther)
// }

// func signupPage(res http.ResponseWriter, req *http.Request) {
// 	Trace.Println("SignUp Page")
// 	if alreadyLoggedIn(req) {
// 		http.Redirect(res, req, "/", http.StatusSeeOther)
// 		return
// 	}
// 	var myUser *user
// 	var mLog messageLog
// 	// process form submission
// 	if req.Method == http.MethodPost { //POST
// 		// get form values
// 		username := strings.TrimSpace(req.FormValue("username"))
// 		password := req.FormValue("password")
// 		confirmPassword := req.FormValue("confirmPassword")
// 		firstname := strings.TrimSpace(req.FormValue("firstname"))
// 		lastname := strings.TrimSpace(req.FormValue("lastname"))

// 		if username == "" || password == "" || confirmPassword == "" || firstname == "" || lastname == "" {
// 			res.WriteHeader(http.StatusForbidden)
// 			Info.Println("Please fill up all details")
// 			mLog = messageLog{"Please fill up all details"}
// 			fatalErr := tpl.ExecuteTemplate(res, "signup.gohtml", mLog)
// 			if fatalErr != nil {
// 				Warning.Println(fatalErr)
// 			}
// 			return
// 		}

// 		if password != confirmPassword {
// 			res.WriteHeader(http.StatusForbidden)
// 			Info.Println("Passwords are not match. Please make sure you key in the same password")
// 			mLog = messageLog{"Passwords are not match. Please make sure you key in the same password"}
// 			fatalErr := tpl.ExecuteTemplate(res, "signup.gohtml", mLog)
// 			if fatalErr != nil {
// 				Warning.Println(fatalErr)
// 			}
// 			return
// 		}

// 		// check password strength
// 		passwordErr := helper.CheckPasswordStrength(password)
// 		if passwordErr != nil {
// 			res.WriteHeader(http.StatusForbidden)
// 			Info.Println(passwordErr.Error())
// 			mLog = messageLog{passwordErr.Error()}
// 			fatalErr := tpl.ExecuteTemplate(res, "signup.gohtml", mLog)
// 			if fatalErr != nil {
// 				Warning.Println(fatalErr)
// 			}
// 			return
// 		}

// 		if username != "" {
// 			// check if username exist/ taken
// 			if _, exist := mapUsers[username]; exist {
// 				res.WriteHeader(http.StatusForbidden)
// 				Info.Println("Username already taken. Please choose another username")
// 				mLog = messageLog{"Username already taken. Please choose another username"}
// 				fatalErr := tpl.ExecuteTemplate(res, "signup.gohtml", mLog)
// 				if fatalErr != nil {
// 					Warning.Println(fatalErr)
// 				}
// 				return
// 			}

// 			expireTime := time.Now().Add(2 * time.Hour)
// 			// Create and sign the JWT
// 			claims := &Claims{
// 				Username: username,
// 				StandardClaims: jwt.StandardClaims{
// 					ExpiresAt: expireTime.Unix(),
// 				},
// 			}
// 			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 			tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
// 			if err != nil {
// 				Warning.Println("Error during token creation")
// 			}
// 			myCookie := &http.Cookie{
// 				Name:     "myCookie",
// 				Value:    tokenString,
// 				Expires:  expireTime,
// 				HttpOnly: true,
// 				Domain:   "localhost",
// 				Secure:   true,
// 			}
// 			http.SetCookie(res, myCookie)

// 			mapSessions[myCookie.Value] = username
// 			bPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
// 			if err != nil {
// 				res.WriteHeader(http.StatusInternalServerError)
// 				Info.Println("Internal Server Error. Please contact system administrator!")
// 				mLog = messageLog{"Internal Server Error. Please contact system administrator!"}
// 				fatalErr := tpl.ExecuteTemplate(res, "signup.gohtml", mLog)
// 				if fatalErr != nil {
// 					Warning.Println(fatalErr)
// 				}
// 				return
// 			}

// 			// store new user info
// 			myUser = &user{username, bPassword, "user", firstname, lastname}
// 			mapUsers[username] = myUser
// 			Trace.Println(username + " is signed up")
// 		}
// 		// redirect to main index
// 		http.Redirect(res, req, "/", http.StatusSeeOther)
// 		return
// 	}
// 	fatalErr := tpl.ExecuteTemplate(res, "signup.gohtml", mLog)
// 	if fatalErr != nil {
// 		Warning.Println(fatalErr)
// 	}
// }
