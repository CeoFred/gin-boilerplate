package helpers

import (
	"github.com/golang-jwt/jwt"
)

type UpdateAccountInformation struct {
	Country        string `json:"country" validate:"required"`
	Manager        string `json:"manager" validate:"required"`
	PhoneNumber    string `json:"phone_number" validate:"required,numeric"`
	CompanyWebsite string `json:"company_website" validate:"required,url"`
}
type ListProjectOffsets struct {
	Projects []string `json:"projects" validate:"required"`
}

type OtpVerify struct {
	Token string `json:"token" validate:"required,len=5"`
	Email string `json:"email" validate:"required,email"`
}

type AccountReset struct {
	Email string `json:"email" validate:"required,email"`
}

type IError struct {
	Field string
	Tag   string
	Value string
}

type AuthTokenJwtClaim struct {
	Email  string
	Name   string
	UserId string
	jwt.StandardClaims
}
type ProjectType string

type AccountStatus int
