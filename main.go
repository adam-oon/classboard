package main

import (
	"classboard/config"
	"classboard/pkg/helper"
	"database/sql"
	"log"
	"net/http"
	"os"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var tpl *template.Template
var db *sql.DB

const (
	APPURL string = "http://localhost:8080"
	APIURL string = "http://localhost:8080/api/v1"
)

func main() {
	// predefined functions for template
	fMap := template.FuncMap{
		"inc": helper.Inc,
	}
	tpl = template.Must(template.New("").Funcs(fMap).ParseGlob("views/*.gohtml"))

	// initialize env file
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalln("Failed to load .env file", err)
	}

	db = config.GetMySQLDB()
	defer db.Close()

	router := router()

	if fatalErr := http.ListenAndServeTLS(":8080", os.Getenv("CERTIFICATE"), os.Getenv("PRIVATE_KEY"), router); fatalErr != nil {
		log.Fatalln(fatalErr)
	}
}
