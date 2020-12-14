package main

import (
	"classboard/dictionary"
	answermodel "classboard/models/answer"
	classroommodel "classboard/models/classroom"
	questionmodel "classboard/models/question"
	summarymodel "classboard/models/summary"
	usermodel "classboard/models/user"
	userclassmodel "classboard/models/userclass"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

// JSON response struct
type ResMessage struct {
	ResponseText string
}

func classroomHandler(res http.ResponseWriter, req *http.Request) {
	if !isLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)

	id, _ := uuid.NewV4()
	user := getSessionUser(req)

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

			classroom := classroommodel.Classroom{
				Id:      id.String(),
				User_id: user.Id,
				Title:   strings.TrimSpace(classroomJSON.Title),
				Code:    strings.TrimSpace(classroomJSON.Code),
			}

			classroomModel := classroommodel.ClassroomModel{
				Db: db,
			}
			err = classroomModel.SaveClassroom(classroom)
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

			classroomModel := classroommodel.ClassroomModel{
				Db: db,
			}
			classroom, err := classroomModel.GetClassroom(params["classroom_id"])
			classroom.Title = classroomJSON.Title
			classroom.Code = classroomJSON.Code

			err = classroomModel.UpdateClassroom(classroom)
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
	if !isLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	user := getSessionUser(req)

	if isStudent(user.Type) {
		res.WriteHeader(http.StatusForbidden)
		json.NewEncoder(res).Encode(ResMessage{ResponseText: "403 Forbidden Action"})
		return
	}

	classroomModel := classroommodel.ClassroomModel{
		Db: db,
	}

	classroom, err := classroomModel.GetClassroom(params["classroom_id"])
	if err != nil {
		Info.Println(err)
	}
	owner_id := classroom.User_id

	if user.Id != owner_id {
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
			questionInput := questionmodel.QuestionInput{
				Classroom_id: strings.TrimSpace(params["classroom_id"]),
				Question:     strings.TrimSpace(questionJSON.Question),
				Type:         strings.TrimSpace(questionJSON.Type),
				Choice:       strings.TrimSpace(questionJSON.Choice),
				Solution:     strings.TrimSpace(questionJSON.Solution),
			}

			questionModel := questionmodel.QuestionModel{
				Db: db,
			}
			err = questionModel.SaveQuestion(questionInput)
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
	if isLoggedIn(req) {
		http.Redirect(res, req, "/dashboard", http.StatusSeeOther)
		return
	}

	templateErr := tpl.ExecuteTemplate(res, "index.gohtml", nil)
	if templateErr != nil {
		Warning.Println(templateErr)
	}
}

func registerPage(res http.ResponseWriter, req *http.Request) {
	// return to dashboard if login
	if isLoggedIn(req) {
		http.Redirect(res, req, "/dashboard", http.StatusSeeOther)
		return
	}
	templateErr := tpl.ExecuteTemplate(res, "register.gohtml", nil)
	if templateErr != nil {
		Warning.Println(templateErr)
	}
}

func dashboardPage(res http.ResponseWriter, req *http.Request) {
	if !isLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	user := getSessionUser(req)

	var classrooms []classroommodel.Classroom
	var err error
	switch user.Type {
	case "lecturer":
		classroomModel := classroommodel.ClassroomModel{
			Db: db,
		}
		classrooms, err = classroomModel.GetClassroomsByUserId(user.Id)

	case "student":
		userclassModel := userclassmodel.UserClassModel{
			Db: db,
		}
		classrooms, err = userclassModel.GetJoinedClass(user.Id)
	}

	if err != nil {
		Info.Println(err)
	}

	data := struct {
		User      usermodel.User
		Classroom []classroommodel.Classroom
	}{
		user,
		classrooms,
	}

	templateErr := tpl.ExecuteTemplate(res, "dashboard.gohtml", data)
	if templateErr != nil {
		Warning.Println(templateErr)
	}
}

func addClassroomPage(res http.ResponseWriter, req *http.Request) {
	if !isLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	var template string
	user := getSessionUser(req)

	if isLecturer(user.Type) {
		template = "classroom_add.gohtml"
	} else if isStudent(user.Type) {
		template = "403.gohtml"
	}

	templateErr := tpl.ExecuteTemplate(res, template, nil)
	if templateErr != nil {
		Warning.Println(templateErr)
	}
}

func editClassroomPage(res http.ResponseWriter, req *http.Request) {
	if !isLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	params := mux.Vars(req)
	user := getSessionUser(req)

	classroomModel := classroommodel.ClassroomModel{
		Db: db,
	}
	classroom, err := classroomModel.GetClassroom(params["classroom_id"])
	if err != nil {
		Info.Println(err)
	}

	var template string
	if isLecturer(user.Type) && user.Id == classroom.User_id {
		template = "classroom_edit.gohtml"
	} else {
		template = "403.gohtml"
	}

	templateErr := tpl.ExecuteTemplate(res, template, classroom)
	if templateErr != nil {
		Warning.Println(templateErr)
	}
}

func joinClassroomPage(res http.ResponseWriter, req *http.Request) {
	if !isLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	user := getSessionUser(req)

	var template string
	if isLecturer(user.Type) {
		template = "403.gohtml"
	} else {
		template = "classroom_join.gohtml"
	}

	templateErr := tpl.ExecuteTemplate(res, template, nil)
	if templateErr != nil {
		Warning.Println(templateErr)
	}
}

func classroomQuestionPage(res http.ResponseWriter, req *http.Request) {
	if !isLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	params := mux.Vars(req)
	user := getSessionUser(req)

	classroomModel := classroommodel.ClassroomModel{
		Db: db,
	}
	questionModel := questionmodel.QuestionModel{
		Db: db,
	}
	userclassModel := userclassmodel.UserClassModel{
		Db: db,
	}
	answerModel := answermodel.AnswerModel{
		Db: db,
	}

	classroom, err := classroomModel.GetClassroom(params["classroom_id"])
	if err != nil {
		Info.Println(err)
	}

	if isLecturer(user.Type) {
		// check lecturer is the classroom owner
		if user.Id == classroom.User_id {
			questions, err := questionModel.GetQuestionsByClassroomId(params["classroom_id"])
			if err != nil {
				Info.Println(err)
			}
			data := struct {
				Classroom classroommodel.Classroom
				Questions []questionmodel.Question
			}{
				classroom,
				questions,
			}
			templateErr := tpl.ExecuteTemplate(res, "question.gohtml", data)
			if templateErr != nil {
				Warning.Println(templateErr)
			}
			return
		}

	} else if isStudent(user.Type) {
		// check student is joined the class
		isStudentClass := userclassModel.IsBelongToClassroom(user.Id, params["classroom_id"])
		if isStudentClass {
			questions, err := questionModel.GetQuestionsByClassroomId(params["classroom_id"])
			if err != nil {
				Info.Println(err)
			}

			studentAnswers := make(map[int]int)

			for _, v := range questions {
				// isCorrect is placeholder to determine answer status
				// -1 = incorrect answer, 0 = no answer,1 = correct answer
				var isCorrect int
				answer, _ := answerModel.GetAnswer(v.Id, user.Id)
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
				Classroom classroommodel.Classroom
				Questions []questionmodel.Question
				Answers   map[int]int
			}{
				classroom,
				questions,
				studentAnswers,
			}
			templateErr := tpl.ExecuteTemplate(res, "question_student.gohtml", data)
			if templateErr != nil {
				Warning.Println(templateErr)
			}
			return
		}
	}

	templateErr := tpl.ExecuteTemplate(res, "403.gohtml", nil)
	if templateErr != nil {
		Warning.Println(templateErr)
	}

}

func addQuestionPage(res http.ResponseWriter, req *http.Request) {
	if !isLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}
	params := mux.Vars(req)
	user := getSessionUser(req)

	var template string
	if isLecturer(user.Type) {
		template = "question_add.gohtml"
	} else if isStudent(user.Type) {
		template = "403.gohtml"
	}

	templateErr := tpl.ExecuteTemplate(res, template, params["classroom_id"])
	if templateErr != nil {
		Warning.Println(templateErr)
	}
}

func joinClassHandler(res http.ResponseWriter, req *http.Request) {
	if !isLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)

	user := getSessionUser(req)

	if isLecturer(user.Type) {
		res.WriteHeader(http.StatusForbidden)
		json.NewEncoder(res).Encode(ResMessage{ResponseText: "403 Forbidden Action"})
		return
	}

	if req.Method == http.MethodPost {
		classroom_id := params["classroom_id"]

		userclassModel := userclassmodel.UserClassModel{
			Db: db,
		}
		err := userclassModel.JoinClass(user.Id, classroom_id)
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
	if !isLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	params := mux.Vars(req)
	user := getSessionUser(req)

	if isStudent(user.Type) {
		res.WriteHeader(http.StatusForbidden)
		json.NewEncoder(res).Encode(ResMessage{ResponseText: "403 Forbidden Action"})
		return
	}

	classroomModel := classroommodel.ClassroomModel{
		Db: db,
	}
	classroom, err := classroomModel.GetClassroom(params["classroom_id"])
	if err != nil {
		Info.Println(err)
	}
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
		answerModel := answermodel.AnswerModel{
			Db: db,
		}
		questionModel := questionmodel.QuestionModel{
			Db: db,
		}

		err = answerModel.DeleteAnswer(question_id)
		if err != nil {
			res.WriteHeader(http.StatusForbidden)
			json.NewEncoder(res).Encode(ResMessage{ResponseText: "403 Forbidden Action"})
			return
		}

		err = questionModel.DeleteQuestion(question_id)
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
	if !isLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	params := mux.Vars(req)
	user := getSessionUser(req)

	questionModel := questionmodel.QuestionModel{
		Db: db,
	}
	answerModel := answermodel.AnswerModel{
		Db: db,
	}

	// check question_id is belongs to classroom_
	question_id, _ := strconv.Atoi(params["question_id"])
	question, err := questionModel.GetQuestion(question_id)
	if err != nil {
		Info.Println(err)
	}

	// if question not belongs to the class
	if question.Classroom_id != params["classroom_id"] {
		templateErr := tpl.ExecuteTemplate(res, "403.gohtml", nil)
		if templateErr != nil {
			Warning.Println(templateErr)
		}
		return
	}

	if isLecturer(user.Type) {
		templateErr := tpl.ExecuteTemplate(res, "question_detail.gohtml", question)
		if templateErr != nil {
			Warning.Println(templateErr)
		}
	} else if isStudent(user.Type) {
		answer, err := answerModel.GetAnswer(question_id, user.Id)
		if err != nil {
			Warning.Println(err)
		}

		var templateErr error
		if answer != nil { // already answered
			data := struct {
				Question   questionmodel.Question
				IsAnswered bool
				Answer     answermodel.Answer
			}{
				question,
				true,
				*answer,
			}
			templateErr = tpl.ExecuteTemplate(res, "answer_question.gohtml", data)
		} else {
			data := struct { // not yet answer
				Question   questionmodel.Question
				IsAnswered bool
				Answer     answermodel.Answer
			}{
				question,
				false,
				answermodel.Answer{},
			}
			templateErr = tpl.ExecuteTemplate(res, "answer_question.gohtml", data)
		}

		if templateErr != nil {
			Warning.Println(templateErr)
		}
	}
}

func answerHandler(res http.ResponseWriter, req *http.Request) {
	if !isLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	user := getSessionUser(req)

	if isLecturer(user.Type) {
		res.WriteHeader(http.StatusForbidden)
		json.NewEncoder(res).Encode(ResMessage{ResponseText: "403 Forbidden Action"})
		return
	}

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

			questionModel := questionmodel.QuestionModel{
				Db: db,
			}
			answerModel := answermodel.AnswerModel{
				Db: db,
			}

			question, err := questionModel.GetQuestion(answerJSON.Question_id)
			if err != nil {
				Info.Println(err)
			}

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
			answer := answermodel.Answer{
				Question_id: answerJSON.Question_id,
				User_id:     user.Id,
				Answer:      answerJSON.Answer,
				Is_correct:  isCorrect,
			}

			err = answerModel.SaveAnswer(answer)
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
	if !isLoggedIn(req) {
		http.Redirect(res, req, "/", http.StatusSeeOther)
		return
	}

	// middleware
	params := mux.Vars(req)
	user := getSessionUser(req)

	classroomModel := classroommodel.ClassroomModel{
		Db: db,
	}
	userclassModel := userclassmodel.UserClassModel{
		Db: db,
	}
	questionModel := questionmodel.QuestionModel{
		Db: db,
	}
	userModel := usermodel.UserModel{
		Db: db,
	}
	answerModel := answermodel.AnswerModel{
		Db: db,
	}

	classroom, err := classroomModel.GetClassroom(params["classroom_id"])
	if err != nil {
		Info.Println(err)
	}
	owner_id := classroom.User_id

	if user.Id != owner_id || isStudent(user.Type) {
		templateErr := tpl.ExecuteTemplate(res, "403.gohtml", nil)
		if templateErr != nil {
			Warning.Println(templateErr)
		}
		return
	}

	var answered, correct int

	// get class's student
	students, err := userclassModel.GetClassroomStudent(params["classroom_id"])
	if err != nil {
		Info.Println(err)
	}
	// get class question
	questions, err := questionModel.GetQuestionsByClassroomId(params["classroom_id"])
	if err != nil {
		Info.Println(err)
	}
	// get answer from student and put into dictionary type
	var dict *dictionary.Dictionary = &dictionary.Dictionary{}
	for _, user_id := range students {
		user, err := userModel.GetUser(user_id)
		if err != nil {
			Info.Println(err)
		}
		var rm *dictionary.ResultMap = &dictionary.ResultMap{}
		for _, question := range questions {
			var result int
			answer, err := answerModel.GetAnswer(question.Id, user.Id)
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
	var summary summarymodel.Summary
	totalQuestion := dict.GetSize() * len(questions)
	summary.StudentTotal = dict.GetSize()
	summary.QuestionTotal = len(questions)
	summary.Participation = summarymodel.CalculateRatio(totalQuestion, answered)
	summary.Correctness = summarymodel.CalculateRatio(totalQuestion, correct)

	data := struct {
		Questions []questionmodel.Question
		Result    *dictionary.Dictionary
		Classroom classroommodel.Classroom
		Summary   summarymodel.Summary
	}{
		questions,
		dict,
		classroom,
		summary,
	}

	templateErr := tpl.ExecuteTemplate(res, "classroom_summary.gohtml", data)
	if templateErr != nil {
		Warning.Println(templateErr)
	}
}
