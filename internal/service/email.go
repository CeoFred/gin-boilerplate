package service

import (
	"fmt"
	"log"
	"time"

	"github.com/CeoFred/gin-boilerplate/constants"
	"github.com/CeoFred/gin-boilerplate/internal/helpers"
	"github.com/CeoFred/gin-boilerplate/internal/otp"
	"github.com/CeoFred/gin-boilerplate/sendgrid"
)

type EmailServicer interface {
	SendVerificationEmail(name, email, url string)
	SendForgotPasswordEmail(name, email string)
}

type EmailService struct {
	Client *sendgrid.Client
}

var (
	constant = constants.New()
)

func NewEmailService() EmailServicer {
	client := sendgrid.NewClient(constant.SendGridApiKey)
	return &EmailService{Client: client}
}

func (s *EmailService) Send(name, email, subject, content string) error {
	to := sendgrid.EmailAddress{
		Name:  name,
		Email: email,
	}

	return s.Client.Send(&to, constant.SendFromEmail, constant.SendFromName, subject, content)

}

func (s *EmailService) SendForgotPasswordEmail(name, email string) {
	//TODO: send otp email to user
}

// Sends account verification email
func (s *EmailService) SendVerificationEmail(name, email, url string) {

	otpToken, err := otp.OTPManage.GenerateOTP(email, time.Minute*10)
	type OTP struct {
		Otp  string
		Name string
		Url  string
	}

	if err != nil {
		log.Printf("Error sending email: %v", err.Error())
	}

	verificationUrl := fmt.Sprintf("%s/api/v1/auth/verify/%s/%s", url, email, otpToken)

	messageBody, err := helpers.ParseTemplateFile("verify_account.html", OTP{Otp: otpToken, Name: name, Url: verificationUrl})

	if err != nil {
		log.Printf("Error sending email: %v", err.Error())
	}

	err = s.Send(name, email, "Verify your accounf", messageBody)

	if err != nil {
		log.Printf("Error sending email: %v", err.Error())
	}
}
