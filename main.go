package main

import (
	"classboard/config"
	"classboard/helper"
	"database/sql"
	"net/http"
	"os"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var tpl *template.Template
var db *sql.DB

func main() {
	// predefined functions for template
	fMap := template.FuncMap{
		"inc":        helper.Inc,
		"strToSlice": helper.StrToSlice,
	}
	tpl = template.Must(template.New("").Funcs(fMap).ParseGlob("views/*.gohtml"))

	// initialize env file
	if err := godotenv.Load(".env"); err != nil {
		Error.Fatalln("Failed to load .env file", err)
	}

	var dbErr error
	db, dbErr = config.GetMySQLDB()
	if dbErr != nil {
		Error.Fatalln(dbErr)
	}
	defer db.Close()

	router := router()

	if fatalErr := http.ListenAndServeTLS(":8080", os.Getenv("CERTIFICATE"), os.Getenv("PRIVATE_KEY"), router); fatalErr != nil {
		Error.Fatalln(fatalErr)
	}
}
