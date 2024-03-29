// validator/validator_test.go
package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCase struct {
	Name   string
	Input  interface{}
	Errors []string
}

func TestValidate(t *testing.T) {
	tests := []testCase{
		{
			Name: "Valid Case",
			Input: struct {
				Username string `validate:"required"`
			}{
				Username: "john_doe",
			},
			Errors: nil,
		},
		{
			Name: "Invalid Case",
			Input: struct {
				Username string `validate:"required"`
			}{
				Username: "",
			},
			Errors: []string{"Username is required"},
		},
		{
			Name: "Non-Aphanumeric Case",
			Input: struct {
				Username string `validate:"required,alpha"`
			}{
				Username: "--[;]a",
			},
			Errors: []string{"Username should be alphanumeric only"},
		},
		// Add more test cases as needed
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			errors := Validate(test.Input)
			assert.ElementsMatch(t, test.Errors, errors)
		})
	}
}
