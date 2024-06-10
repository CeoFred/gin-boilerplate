package models

import (
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
)

type AccountRole string
type AccountStatus string

const (
	UserRole  AccountRole = "user"
	AdminRole AccountRole = "admin"
)

const (
	ActiveAccount    AccountStatus = "Active"
	SuspendedAccount AccountStatus = "Suspended"
	InactiveAccount  AccountStatus = "Inactive"
)

type User struct {
	Email         string      `json:"email"`
	Password      string      `json:"password,omitempty"`
	LastLogin     string      `json:"last_login"`
	IP            string      `json:"ip"`
	ID            uuid.UUID   `json:"id" validate:"required"`
	Role          AccountRole `json:"role" validate:"required"`
	EmailVerified bool        `json:"email_verified" validate:"required"`
	Country       string      `json:"country"`
	PhoneNumber   string      `json:"phone_number"`
	FirstName     string      `json:"first_name" validate:"required"`
	LastName      string      `json:"last_name" validate:"required"`
	Status        string      `json:"status"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

func (u *User) MarshalJSON() ([]byte, error) {
	type Alias User
	return json.Marshal(&struct {
		Password string `json:"password,omitempty"`
		*Alias
	}{
		Password: "",
		Alias:    (*Alias)(u),
	})
}
