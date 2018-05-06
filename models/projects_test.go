package models

import (
	"fmt"
	"testing"
)

func TestDBConfig(t *testing.T) {
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

type productsQueryTestCases struct {
	tag      string
	id       string
	expected string
	isValid  bool
}

func TestProductsQuery(t *testing.T) {
	fmt.Println("Running TestProductsQuery...")

	testCases := []productsQueryTestCases{
		{
			// implicit id=nil
			tag:      "Case 1",
			expected: "SELECT p.id, p.name, a.id, p.created_at FROM projects p JOIN assets a ON a.project_id=p.id WHERE a.category=1;",
			isValid:  false,
		},
		{
			tag:      "Case 2",
			id:       "-1",
			expected: "SELECT p.id, p.name, a.id, p.created_at FROM projects p JOIN assets a ON a.project_id=p.id WHERE a.category=1;",
			isValid:  false,
		},
		{
			tag:      "Case 3",
			id:       "99",
			expected: "SELECT p.id, p.name, a.id, p.created_at FROM projects p JOIN assets a ON a.project_id=p.id WHERE a.category=1 AND p.id=99;",
			isValid:  true,
		},
	}

	for _, testCase := range testCases {
		fmt.Println("Testing ProductsQuery for id: ", testCase.id)
		query := NewProjectsQuery(testCase.id, "", "")
		err := query.Validate()
		if testCase.isValid != (err == nil) {
			t.Errorf("%s: Validation failed with error:\n%s", testCase.tag, err)
		}
		actual := query.Build()
		if testCase.expected != actual {
			t.Errorf("%s: Evaluation returned an unexpected output.\nExpected: %s\nActual: %s", testCase.tag, testCase.expected, actual)
		}
	}
}
