package main

import (
	"classboard/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

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
	params := mux.Vars(req)
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
	classroom := models.GetClassroom(params["classroom_id"])
	owner_id := classroom.User_id

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

			delimiter := "|"
			choiceSlice := strings.Split(questionJSON.Choice, delimiter)
			var sanitizeChoiceSlice []string
			var hasSolution bool
			for _, v := range choiceSlice {
				v = strings.TrimSpace(v)
				if v != "" {
					sanitizeChoiceSlice = append(sanitizeChoiceSlice, v)
				}
				if strings.TrimSpace(questionJSON.Solution) == v {
					hasSolution = true
				}
			}
			questionJSON.Choice = strings.Join(sanitizeChoiceSlice, delimiter)

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

	userModel := models.UserModel{
		Db: db,
	}
	user := userModel.GetUser(user_id)

	var classrooms []models.Classroom
	switch user.Type {
	case "lecturer":
		classrooms = models.GetClassroomsByUserId(user_id)
	case "student":
		classrooms = models.GetJoinedClass(user_id)
	}

	data := struct {
		User      models.User
		Classroom []models.Classroom
	}{
		user,
		classrooms,
	}

	fatalErr := tpl.ExecuteTemplate(res, "dashboard.gohtml", data)
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

func joinClassroomPage(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	fatalErr := tpl.ExecuteTemplate(res, "classroom_join.gohtml", nil)
	if fatalErr != nil {
		log.Println(fatalErr)
	}
}

func classroomQuestionPage(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	params := mux.Vars(req)
	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		//error//
	}

	sessionModel := models.SessionModel{
		Db: db,
	}
	user_id := sessionModel.GetUserID(myCookie.Value)
	userModel := models.UserModel{
		Db: db,
	}
	user := userModel.GetUser(user_id)
	classroom := models.GetClassroom(params["classroom_id"])

	if user.Type == "lecturer" {
		// check lecturer is the classroom owner
		if user_id == classroom.User_id {
			questions := models.GetQuestionsByClassroomId(params["classroom_id"])
			data := struct {
				Classroom models.Classroom
				Questions []models.Question
			}{
				classroom,
				questions,
			}
			fatalErr := tpl.ExecuteTemplate(res, "question.gohtml", data)
			if fatalErr != nil {
				log.Println(fatalErr)
			}
			return
		}

	} else if user.Type == "student" {
		// check student is joined the class
		isStudentClass := models.IsBelongToClassroom(user_id, params["classroom_id"])
		if isStudentClass {
			questions := models.GetQuestionsByClassroomId(params["classroom_id"])
			data := struct {
				Classroom models.Classroom
				Questions []models.Question
			}{
				classroom,
				questions,
			}
			fatalErr := tpl.ExecuteTemplate(res, "question_student.gohtml", data)
			if fatalErr != nil {
				log.Println(fatalErr)
			}
			return
		}
	}

	fatalErr := tpl.ExecuteTemplate(res, "403.gohtml", nil)
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

func joinClassHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
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
	userModel := models.UserModel{
		Db: db,
	}
	user := userModel.GetUser(user_id)

	if user.Type == "lecturer" {
		res.WriteHeader(http.StatusForbidden)
		json.NewEncoder(res).Encode(ResMessage{ResponseText: "Forbidden Action"})
		return
	}

	if req.Method == http.MethodPost {
		classroom_id := params["classroom_id"]

		err = models.JoinClass(user_id, classroom_id)
		if err != nil {
			res.WriteHeader(http.StatusUnprocessableEntity)
			json.NewEncoder(res).Encode(ResMessage{ResponseText: err.Error()})
			return
		} else {
			res.WriteHeader(http.StatusCreated)
			json.NewEncoder(res).Encode(ResMessage{ResponseText: "Class joined!"})
		}
	}
}

func questionHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	user := getUser(req)

	classroom := models.GetClassroom(params["classroom_id"])
	owner_id := classroom.User_id

	if user.Id != owner_id {
		res.WriteHeader(http.StatusForbidden)
		json.NewEncoder(res).Encode(ResMessage{ResponseText: "Forbidden Action"})
		return
	}

	switch req.Method {
	case "DELETE":
		question_id, err := strconv.Atoi(params["question_id"])
		if err != nil {
			res.WriteHeader(http.StatusForbidden)
			json.NewEncoder(res).Encode(ResMessage{ResponseText: "Forbidden Action"})
			return
		}

		// delete student answer first, then delete question
		err = models.DeleteAnswer(question_id)
		if err != nil {
			res.WriteHeader(http.StatusForbidden)
			json.NewEncoder(res).Encode(ResMessage{ResponseText: "Forbidden Action"})
			return
		}

		err = models.DeleteQuestion(question_id)
		if err != nil {
			res.WriteHeader(http.StatusForbidden)
			json.NewEncoder(res).Encode(ResMessage{ResponseText: "Forbidden Action"})
			return
		}

		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(ResMessage{ResponseText: "Question Deleted"})

	}

}

func questionPage(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	params := mux.Vars(req)
	user := getUser(req)
	// check question_id is belongs to classroom_

	question_id, _ := strconv.Atoi(params["question_id"])
	question := models.GetQuestion(question_id)

	// if question not belongs to the class
	if question.Classroom_id != params["classroom_id"] {
		fatalErr := tpl.ExecuteTemplate(res, "403.gohtml", nil)
		if fatalErr != nil {
			Warning.Println(fatalErr)
		}
		return
	}

	var template string
	if user.Type == "lecturer" {
		template = "question_detail.gohtml"
	} else if user.Type == "student" {
		// answer, err := models.GetAnswer(question_id, user.Id)
		template = "answer_question.gohtml"
	}

	fatalErr := tpl.ExecuteTemplate(res, template, question)
	if fatalErr != nil {
		Warning.Println(fatalErr)
	}
}

func answerHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	// params := mux.Vars(req)

	// get question solution
	user := getUser(req)

	// if user_id != owner_id {
	// 	res.WriteHeader(http.StatusForbidden)
	// 	json.NewEncoder(res).Encode(ResMessage{ResponseText: "Forbidden Action"})
	// 	return
	// }

	if req.Method == http.MethodPost && req.Header.Get("Content-type") == "application/json" {
		reqBody, err := ioutil.ReadAll(req.Body)
		type AnswerJSON struct {
			Question_id int
			Answer      string
		}
		var answerJSON AnswerJSON
		if err == nil {
			err := json.Unmarshal(reqBody, &answerJSON)
			fmt.Println(err)
			if err != nil {
				res.WriteHeader(http.StatusUnprocessableEntity)
				json.NewEncoder(res).Encode(ResMessage{ResponseText: "Sorry the answer info is incomplete"})
				return
			}

			question := models.GetQuestion(answerJSON.Question_id)
			var isCorrect bool
			if question.Solution == answerJSON.Answer {
				isCorrect = true
			}
			type Answer struct {
				Question_id int
				User_id     int
				Answer      string
				Is_correct  bool
			}
			answer := models.Answer{
				Question_id: answerJSON.Question_id,
				User_id:     user.Id,
				Answer:      answerJSON.Answer,
				Is_correct:  isCorrect,
			}

			err = models.SaveAnswer(answer)
			if err != nil {
				res.WriteHeader(http.StatusUnprocessableEntity)
				json.NewEncoder(res).Encode(ResMessage{ResponseText: err.Error()})
				return
			} else {
				var result string
				if isCorrect {
					result = "Correct Answer!"
				} else {
					result = "Incorrect Answer! Correct Answer is " + question.Solution
				}
				res.WriteHeader(http.StatusOK)
				json.NewEncoder(res).Encode(ResMessage{ResponseText: result})
				return
			}
		}
	}
}

func reportPage(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	// question id, user_id, username,classroom_id
	params := mux.Vars(req)
	user := getUser(req)

	classroom := models.GetClassroom(params["classroom_id"])
	owner_id := classroom.User_id

	if user.Id != owner_id {
		res.WriteHeader(http.StatusForbidden)
		json.NewEncoder(res).Encode(ResMessage{ResponseText: "Forbidden Action"})
		return
	}
	// get class's student
	students := models.GetClassroomStudent(params["classroom_id"])
	fmt.Println(students)
	// get class question
	questions := models.GetQuestionsByClassroomId(params["classroom_id"])
	fmt.Println(questions)
	// get answer from student

	fatalErr := tpl.ExecuteTemplate(res, "classroom_report.gohtml", nil)
	if fatalErr != nil {
		log.Println(fatalErr)
	}
}
