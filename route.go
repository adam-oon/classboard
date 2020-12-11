package main

import "github.com/gorilla/mux"

func router() *mux.Router {
	router := mux.NewRouter()
	// page
	router.HandleFunc("/", indexPage)
	router.HandleFunc("/register", registerPage).Methods("GET")
	router.HandleFunc("/register", register).Methods("POST")
	router.HandleFunc("/login", login).Methods("POST")
	router.HandleFunc("/logout", logout)
	router.HandleFunc("/dashboard", dashboardPage)
	router.HandleFunc("/classroom/add", addClassroomPage)
	router.HandleFunc("/classroom/join", joinClassroomPage)
	router.HandleFunc("/classroom/{classroom_id}/question", classroomQuestionPage)
	router.HandleFunc("/classroom/{classroom_id}/question/add", addQuestionPage)
	router.HandleFunc("/classroom/{classroom_id}/question/{question_id}", questionPage)
	router.HandleFunc("/classroom/{classroom_id}/report", reportPage)

	// api
	router.HandleFunc("/api/v1/classroom", addClassroomHandler).Methods("POST")
	router.HandleFunc("/api/v1/{classroom_id}/question", addQuestionHandler).Methods("POST")
	router.HandleFunc("/api/v1/{classroom_id}/question/{question_id}", questionHandler).Methods("DELETE")
	router.HandleFunc("/api/v1/user_class/{classroom_id}", joinClassHandler).Methods("POST")
	router.HandleFunc("/api/v1/answer", answerHandler).Methods("POST")

	return router
}
