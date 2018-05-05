package models

import (
	"fmt"
	"testing"
)

func TestDBConfig(t *testing.T){
	fmt.Println("Running TestDBConfig...")

	dbConfig := DBConfig{"user", "password", "host", "db", "5432"}
	actual := dbConfig.ConnectionString()
	expectation := "postgresql://user:password@host:5432/db?sslmode=disable"

	if expectation != actual {
		t.Errorf(
			"Generation of db connection string from DBConfig object failed.\nWas expecting: %s\nReceieved: %s",
			expectation, actual,
		)
	}
}
