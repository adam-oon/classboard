package config

import "testing"

func TestGetMySQLDB(t *testing.T) {
	db, err := GetMySQLDB()
	defer db.Close()

	if err != nil {
		t.Errorf(err.Error())
	}

	err = db.Ping()
	if err != nil {
		t.Errorf(err.Error())
	}
}
