package main

import (
	"classboard/models"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

func classroomHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	// params := mux.Vars(req)

	// myCookie, err := req.Cookie("myCookie")
	// if err != nil {
	// 	res.WriteHeader(http.StatusUnprocessableEntity)
	// 	json.NewEncoder(res).Encode(ResMessage{ResponseText: "Invalid User"})
	// 	return
	// }

	// switch req.Method {
	// case "DELETE":
	// 	// delete student answer, delete question, delete student classroom, delete classroom
	// 	sessionModel := models.SessionModel{
	// 		Db: db,
	// 	}
	// 	user_id := sessionModel.GetUserID(myCookie.Value)

	// }

}

func addClassroomHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(res).Encode(ResMessage{ResponseText: "Sorry the classroom info is incomplete"})
		return
	}

	sessionModel := models.SessionModel{
		Db: db,
	}
	user_id := sessionModel.GetUserID(myCookie.Value)
	id, _ := uuid.NewV4()

	if req.Method == http.MethodPost && req.Header.Get("Content-type") == "application/json" {
		reqBody, err := ioutil.ReadAll(req.Body)
		type ClassroomJSON struct {
			Title string
			Code  string
		}
		var classroomJSON ClassroomJSON
		if err == nil {
			err := json.Unmarshal(reqBody, &classroomJSON)
			if err != nil {
				res.WriteHeader(http.StatusUnprocessableEntity)
				json.NewEncoder(res).Encode(ResMessage{ResponseText: "Sorry the classroom info is incomplete"})
				return
			}

			classroom := models.Classroom{
				Id:      id.String(),
				User_id: user_id,
				Title:   strings.TrimSpace(classroomJSON.Title),
				Code:    strings.TrimSpace(classroomJSON.Code),
			}

			err = models.SaveClassroom(classroom)
			if err != nil {
				res.WriteHeader(http.StatusUnprocessableEntity)
				json.NewEncoder(res).Encode(ResMessage{ResponseText: err.Error()})
				return
			} else {
				res.WriteHeader(http.StatusCreated)
				json.NewEncoder(res).Encode(ResMessage{ResponseText: "Classroom created!"})
			}
		}
	}
}

func addQuestionHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req) //
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(res).Encode(ResMessage{ResponseText: "Sorry the question info is incomplete"})
		return
	}

	sessionModel := models.SessionModel{
		Db: db,
	}
	user_id := sessionModel.GetUserID(myCookie.Value)
	owner_id := models.GetClassroomOwner(params["classroom_id"])

	if user_id != owner_id {
		res.WriteHeader(http.StatusForbidden)
		json.NewEncoder(res).Encode(ResMessage{ResponseText: "Forbidden Action"})
		return
	}

	if req.Method == http.MethodPost && req.Header.Get("Content-type") == "application/json" {
		reqBody, err := ioutil.ReadAll(req.Body)
		type QuestionJSON struct {
			Question string
			Type     string
			Choice   string
			Solution string
		}
		var questionJSON QuestionJSON
		if err == nil {
			err := json.Unmarshal(reqBody, &questionJSON)
			if err != nil {
				res.WriteHeader(http.StatusUnprocessableEntity)
				json.NewEncoder(res).Encode(ResMessage{ResponseText: "Sorry the question info is incomplete"})
				return
			}

			// validate question
			// question := strings.Replace(questionJSON.Question, ",", " ", -1)
			// question = strings.TrimSpace(question)
			// questionJSON.Question = strings.Replace(question, " ", ",", -1)
			questionSlice := strings.Split(questionJSON.Question, ",")
			var sanitizeQuestionSlice []string
			var hasSolution bool
			for _, v := range questionSlice {
				if v != " " {
					sanitizeQuestionSlice = append(sanitizeQuestionSlice, v)
				}

				if questionJSON.Solution == v {
					hasSolution = true
				}
			}

			questionJSON.Question = strings.Join(sanitizeQuestionSlice, ",")

			if !hasSolution {
				res.WriteHeader(http.StatusUnprocessableEntity)
				json.NewEncoder(res).Encode(ResMessage{ResponseText: "Solution is not found in choice!!"})
				return
			}

			// validate solution
			questionInput := models.QuestionInput{
				Classroom_id: strings.TrimSpace(params["classroom_id"]),
				Question:     strings.TrimSpace(questionJSON.Question),
				Type:         strings.TrimSpace(questionJSON.Type),
				Choice:       strings.TrimSpace(questionJSON.Choice),
				Solution:     strings.TrimSpace(questionJSON.Solution),
			}

			err = models.SaveQuestion(questionInput)
			if err != nil {
				res.WriteHeader(http.StatusUnprocessableEntity)
				json.NewEncoder(res).Encode(ResMessage{ResponseText: err.Error()})
				return
			} else {
				res.WriteHeader(http.StatusCreated)
				json.NewEncoder(res).Encode(ResMessage{ResponseText: "Question created!"})
			}
		}
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

// func checkUserType(req *http.Request) bool {
// 	myCookie, err := req.Cookie("myCookie")
// 	if err != nil {
// 		return false
// 	}
// 	sessionModel := models.SessionModel{
// 		Db: db,
// 	}
// 	ok := sessionModel.CheckSession(myCookie.Value)
// 	return ok
// }

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

	// lecturer
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(res).Encode(ResMessage{ResponseText: "Sorry the classroom info is incomplete"})
		return
	}

	sessionModel := models.SessionModel{
		Db: db,
	}
	lecturer_id := sessionModel.GetUserID(myCookie.Value)
	classrooms := models.GetClassroomsByUserId(lecturer_id)
	//student

	fatalErr := tpl.ExecuteTemplate(res, "dashboard.gohtml", classrooms)
	if fatalErr != nil {
		log.Println(fatalErr)
	}
}

func addClassroomPage(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	fatalErr := tpl.ExecuteTemplate(res, "classroom_add.gohtml", nil)
	if fatalErr != nil {
		log.Println(fatalErr)
	}
}

func classroomQuestionPage(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	// lecturer
	// myCookie, err := req.Cookie("myCookie")
	// if err != nil {
	// 	res.WriteHeader(http.StatusUnprocessableEntity)
	// 	json.NewEncoder(res).Encode(ResMessage{ResponseText: "Sorry the classroom info is incomplete"})
	// 	return
	// }

	// sessionModel := models.SessionModel{
	// 	Db: db,
	// }
	// lecturer_id := sessionModel.GetUserID(myCookie.Value)
	// classrooms := models.GetClassroomsByUserId(lecturer_id)
	// //student

	fatalErr := tpl.ExecuteTemplate(res, "question.gohtml", nil)
	if fatalErr != nil {
		log.Println(fatalErr)
	}
}

func addQuestionPage(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	params := mux.Vars(req)
	fatalErr := tpl.ExecuteTemplate(res, "question_add.gohtml", params["classroom_id"])
	if fatalErr != nil {
		log.Println(fatalErr)
	}
}
