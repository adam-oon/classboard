/*
	Package config provide app configuration
*/
package config

import (
	"database/sql"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

//GetMySQLDB return sql.DB connection instance
func GetMySQLDB() (*sql.DB, error) {
	dbDriver := "mysql"
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbIP := os.Getenv("DB_IP")
	dbSchema := os.Getenv("DB_SCHEMA")
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+dbIP+")/"+dbSchema)
	if err != nil {
		return nil, err
	}
	return db, nil
}
