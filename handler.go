package main

import (
	"classboard/dictionary"
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

func classroomHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)

	myCookie, err := req.Cookie("myCookie")
	if err != nil {
		res.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(res).Encode(ResMessage{ResponseText: "Sorry the classroom info is incomplete"})
		return
	}

	user_id := models.GetUserID(myCookie.Value)
	id, _ := uuid.NewV4()

	user := getUser(req)

	if isStudent(user.Type) {
		res.WriteHeader(http.StatusForbidden)
		json.NewEncoder(res).Encode(ResMessage{ResponseText: "403 Forbidden Action"})
		return
	}

	// post
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

			} else {
				res.WriteHeader(http.StatusCreated)
				json.NewEncoder(res).Encode(ResMessage{ResponseText: "Classroom created!"})
			}
		}
		return
	}

	// put
	if req.Method == http.MethodPut && req.Header.Get("Content-type") == "application/json" {
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

			classroom := models.GetClassroom(params["classroom_id"])
			classroom.Title = classroomJSON.Title
			classroom.Code = classroomJSON.Code

			err = models.UpdateClassroom(classroom)
			if err != nil {
				res.WriteHeader(http.StatusUnprocessableEntity)
				json.NewEncoder(res).Encode(ResMessage{ResponseText: err.Error()})
				return
			} else {
				res.WriteHeader(http.StatusCreated)
				json.NewEncoder(res).Encode(ResMessage{ResponseText: "Classroom updated!"})
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

	user_id := models.GetUserID(myCookie.Value)
	classroom := models.GetClassroom(params["classroom_id"])
	owner_id := classroom.User_id

	if user_id != owner_id {
		res.WriteHeader(http.StatusForbidden)
		json.NewEncoder(res).Encode(ResMessage{ResponseText: "403 Forbidden Action"})
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

	user_id := models.GetUserID(myCookie.Value)
	user := models.GetUser(user_id)

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

	var template string
	user := getUser(req)

	if isLecturer(user.Type) {
		template = "classroom_add.gohtml"
	} else if isStudent(user.Type) {
		template = "403.gohtml"
	}

	fatalErr := tpl.ExecuteTemplate(res, template, nil)
	if fatalErr != nil {
		log.Println(fatalErr)
	}
}

func editClassroomPage(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	params := mux.Vars(req)
	user := getUser(req)
	classroom := models.GetClassroom(params["classroom_id"])

	var template string
	if isLecturer(user.Type) && user.Id == classroom.User_id {
		template = "classroom_edit.gohtml"
	} else {
		template = "403.gohtml"
	}

	fatalErr := tpl.ExecuteTemplate(res, template, classroom)
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

	user_id := models.GetUserID(myCookie.Value)
	user := models.GetUser(user_id)
	classroom := models.GetClassroom(params["classroom_id"])

	if isLecturer(user.Type) {
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

	} else if isStudent(user.Type) {
		// check student is joined the class
		isStudentClass := models.IsBelongToClassroom(user_id, params["classroom_id"])
		if isStudentClass {
			questions := models.GetQuestionsByClassroomId(params["classroom_id"])

			studentAnswers := make(map[int]int)

			for _, v := range questions {
				// isCorrect is placeholder to determine answer status
				// -1 = incorrect answer, 0 = no answer,1 = correct answer
				var isCorrect int
				answer, _ := models.GetAnswer(v.Id, user.Id)
				if answer != nil {
					if answer.Is_correct {
						isCorrect = 1
					} else {
						isCorrect = -1
					}
				}
				studentAnswers[v.Id] = isCorrect
			}

			data := struct {
				Classroom models.Classroom
				Questions []models.Question
				Answers   map[int]int
			}{
				classroom,
				questions,
				studentAnswers,
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

	user_id := models.GetUserID(myCookie.Value)
	user := models.GetUser(user_id)

	if isLecturer(user.Type) {
		res.WriteHeader(http.StatusForbidden)
		json.NewEncoder(res).Encode(ResMessage{ResponseText: "403 Forbidden Action"})
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
		json.NewEncoder(res).Encode(ResMessage{ResponseText: "403 Forbidden Action"})
		return
	}

	switch req.Method {
	case "DELETE":
		question_id, err := strconv.Atoi(params["question_id"])
		if err != nil {
			res.WriteHeader(http.StatusForbidden)
			json.NewEncoder(res).Encode(ResMessage{ResponseText: "403 Forbidden Action"})
			return
		}

		// delete student answer first, then delete question
		err = models.DeleteAnswer(question_id)
		if err != nil {
			res.WriteHeader(http.StatusForbidden)
			json.NewEncoder(res).Encode(ResMessage{ResponseText: "403 Forbidden Action"})
			return
		}

		err = models.DeleteQuestion(question_id)
		if err != nil {
			res.WriteHeader(http.StatusForbidden)
			json.NewEncoder(res).Encode(ResMessage{ResponseText: "403 Forbidden Action"})
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

	if isLecturer(user.Type) {
		fatalErr := tpl.ExecuteTemplate(res, "question_detail.gohtml", question)
		if fatalErr != nil {
			Warning.Println(fatalErr)
		}
	} else if isStudent(user.Type) {
		answer, err := models.GetAnswer(question_id, user.Id)
		if err != nil {
			Warning.Println(err)
		}

		var fatalErr error
		if answer != nil { // already answered
			data := struct {
				Question   models.Question
				IsAnswered bool
				Answer     models.Answer
			}{
				question,
				true,
				*answer,
			}
			fatalErr = tpl.ExecuteTemplate(res, "answer_question.gohtml", data)
		} else {
			data := struct { // not yet answer
				Question   models.Question
				IsAnswered bool
				Answer     models.Answer
			}{
				question,
				false,
				models.Answer{},
			}
			fatalErr = tpl.ExecuteTemplate(res, "answer_question.gohtml", data)
		}

		if fatalErr != nil {
			Warning.Println(fatalErr)
		}
	}
}

func answerHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	// params := mux.Vars(req)

	// get question solution
	user := getUser(req)

	// if user_id != owner_id {
	// 	res.WriteHeader(http.StatusForbidden)
	// 	json.NewEncoder(res).Encode(ResMessage{ResponseText: "403 Forbidden Action"})
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

func summaryPage(res http.ResponseWriter, req *http.Request) {
	if !alreadyLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	params := mux.Vars(req)
	user := getUser(req)

	classroom := models.GetClassroom(params["classroom_id"])
	owner_id := classroom.User_id

	if user.Id != owner_id {
		res.WriteHeader(http.StatusForbidden)
		json.NewEncoder(res).Encode(ResMessage{ResponseText: "403 Forbidden Action"})
		return
	}

	var answered, correct int

	// get class's student
	students := models.GetClassroomStudent(params["classroom_id"])
	// get class question
	questions := models.GetQuestionsByClassroomId(params["classroom_id"])
	// get answer from student and put into dictionary type
	var dict *dictionary.Dictionary = &dictionary.Dictionary{}
	for _, user_id := range students {
		user := models.GetUser(user_id)
		var rm *dictionary.ResultMap = &dictionary.ResultMap{}
		for _, question := range questions {
			var result int
			answer, err := models.GetAnswer(question.Id, user.Id)
			if err != nil {
				Warning.Println(err)
			} else if answer == nil && err == nil {
				result = 0
			} else if answer != nil {
				if answer.Is_correct {
					result = 1
					answered++
					correct++
				} else {
					result = -1
					answered++
				}
			}
			rm.SetValue(question.Id, result)
		}
		user_identifier := dictionary.NameKey(user.Name + "(" + user.Username + ")")
		dict.SetResultMap(user_identifier, rm)
	}

	// summary
	var summary Summary
	totalQuestion := dict.GetSize() * len(questions)
	summary.StudentTotal = dict.GetSize()
	summary.QuestionTotal = len(questions)
	summary.Participation = CalculateRatio(totalQuestion, answered)
	summary.Correctness = CalculateRatio(totalQuestion, correct)

	data := struct {
		Questions []models.Question
		Result    *dictionary.Dictionary
		Classroom models.Classroom
		Summary   Summary
	}{
		questions,
		dict,
		classroom,
		summary,
	}

	fatalErr := tpl.ExecuteTemplate(res, "classroom_summary.gohtml", data)
	if fatalErr != nil {
		log.Println(fatalErr)
	}
}
