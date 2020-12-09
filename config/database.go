package config

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func GetMySQLDB() *sql.DB {
	dbDriver := "mysql"
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbIP := os.Getenv("DB_IP")
	dbSchema := os.Getenv("DB_SCHEMA")
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+dbIP+")/"+dbSchema)
	if err != nil {
		panic(err.Error())
	}
	return db
}
