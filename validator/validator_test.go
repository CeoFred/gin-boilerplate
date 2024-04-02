package validator

import (
	"testing"
)

type User struct {
	Username string `validate:"required,alpha"`
}

func TestValidate(t *testing.T) {
	validUser := User{Username: "JohnDoe"}
	invalidUser := User{Username: ""}

	t.Run("Valid User", func(t *testing.T) {
		err := Validate(validUser)
		if err != nil {
			t.Errorf("Expected no error for valid user, got: %v", err)
		}
	})

	t.Run("Invalid User", func(t *testing.T) {
		err := Validate(invalidUser)
		if err == nil {
			t.Error("Expected error for invalid user, got nil")
		} else {
			expectedError := "Username is required"
			if err.Error() != expectedError {
				t.Errorf("Expected error message '%s', got '%s'", expectedError, err.Error())
			}
		}
	})
}
