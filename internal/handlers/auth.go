package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/CeoFred/gin-boilerplate/constants"
	"github.com/CeoFred/gin-boilerplate/internal/helpers"
	"github.com/CeoFred/gin-boilerplate/internal/models"
	"github.com/CeoFred/gin-boilerplate/internal/otp"
	"github.com/CeoFred/gin-boilerplate/internal/repository"
	"github.com/CeoFred/gin-boilerplate/sendgrid"
	"github.com/gofrs/uuid"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userRepository *repository.UserRepository
}

func NewAuthHandler(
	userRepo *repository.UserRepository,

) *AuthHandler {
	return &AuthHandler{
		userRepository: userRepo,
	}
}

var (
	constant = constants.New()
)

type ErrorResponse struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
}

type SuccessResponse struct {
	Message int  `json:"message"`
	Success bool `json:"success"`
}

type RegisterResponse struct {
	Success bool                 `json:"success"`
	Message string               `json:"message"`
	Data    RegisterResponseData `json:"data"`
}

type RegisterResponseData struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type AuthenticateUser struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type InputCreateUser struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// LoginResponse represents the response data structure for the login API.
type LoginResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Data    LoginResponseData `json:"data"`
}

// LoginResponseData represents the data section of the login response.
type LoginResponseData struct {
	JWT string `json:"jwt"`
}

type UpdateAccountInformation struct {
	Country        string `json:"country" validate:"required"`
	Manager        string `json:"manager" validate:"required"`
	PhoneNumber    string `json:"phone_number" validate:"required,numeric"`
	CompanyWebsite string `json:"company_website" validate:"required,url"`
}

// ? ForgotPasswordInput struct
type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required"`
}

// ? ResetPasswordInput struct
type ResetPasswordInput struct {
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
}

type OtpVerifyInput struct {
	Token string `json:"token" binding:"required"`
	Email string `json:"email" binding:"required"`
}

// Authenticate authenticates a user and generates a JWT token.
//
// @Summary Authenticate User
// @Description Authenticate a user by validating their email and password.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body AuthenticateUser true "User credentials (email and password)"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/login [post]
func (a *AuthHandler) Authenticate(c *gin.Context) {

	var input AuthenticateUser
	validatedReqBody, exists := c.Get("validatedRequestBody")

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve validated request body"})
		return
	}

	input, ok := validatedReqBody.(AuthenticateUser)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert request body to AuthenticateUser"})
		return
	}

	user, userExist, err := a.userRepository.FindByCondition("email = ?", strings.ToLower(input.Email))
	if err != nil {
		helpers.Dispatch500Error(c, err)
		return
	}

	if !userExist {
		helpers.Dispatch400Error(c, "invalid account credentials", nil)
		return
	}

	if !user.EmailVerified {
		helpers.Dispatch400Error(c, "account not verified", nil)
		return
	}

	hashedPassword := []byte(user.Password)
	plainPassword := []byte(input.Password)
	err = bcrypt.CompareHashAndPassword(hashedPassword, plainPassword)

	if err != nil {
		helpers.Dispatch400Error(c, "email and password does not match", nil)
		return
	}

	timeNow, err := helpers.TimeNow("Africa/Lagos")
	if err != nil {
		helpers.Dispatch400Error(c, err.Error(), err)
		return
	}

	user.LastLogin = timeNow
	user.IP = c.ClientIP()

	_, err = a.userRepository.UpdateUserByCondition("email", user.Email, user)
	if err != nil {
		helpers.Dispatch400Error(c, err.Error(), err)
		return
	}

	jwtToken, err := helpers.GenerateToken(constant.JWTSecretKey, user.Email, user.FirstName, (user.ID).String())
	if err != nil {
		helpers.Dispatch500Error(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "authenticated successfully",
		"data": map[string]string{
			"jwt": jwtToken,
		},
	})
}

func (a *AuthHandler) findUserOrError(email string) (user *models.User, err error) {
	user, userExist, err := a.userRepository.FindByCondition("email = ?", email)
	if err != nil {
		return nil, err
	}
	if !userExist {
		return nil, helpers.NewError("user not found")
	}
	return user, nil
}

// Register creates a new user account.
//
// @Summary Register a new user
// @Description Create a new user account with the provided information
// @Tags Authentication
// @Accept json
// @Produce json
// @Param input body InputCreateUser true "User data to create an account"
// @Success 201 {object} RegisterResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/register [post]
func (a *AuthHandler) Register(c *gin.Context) {

	var input InputCreateUser
	validatedReqBody, exists := c.Get("validatedRequestBody")

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve validated request body"})
		return
	}

	input, ok := validatedReqBody.(InputCreateUser)

	if !ok {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert request body to InputCreateUser"})
		return
	}

	_, err := a.findUserOrError(input.Email)

	if err == nil {
		fmt.Println(err)
		helpers.Dispatch400Error(c, "user already registered", nil)
		return
	}

	hash, err := helpers.HashPassword(input.Password)
	if err != nil {
		helpers.Dispatch400Error(c, "error hashing password", err)
		return
	}

	ID, err := uuid.NewV7()
	if err != nil {
		helpers.Dispatch400Error(c, "error generating uuid", err)
		return
	}
	// create record
	user := &models.User{
		Email:         strings.ToLower(input.Email),
		Password:      hash,
		ID:            ID,
		IP:            c.ClientIP(),
		Role:          models.UserRole,
		EmailVerified: false,
		CreatedAt:     time.Now(),
		Status:        string(models.InactiveAccount),
		FirstName:     input.FirstName,
		LastName:      input.LastName,
	}

	if err := a.userRepository.CreateUser(user); err != nil {
		helpers.Dispatch500Error(c, err)
		return
	}

	scheme := "http" // Default scheme
	isLocal := gin.Mode() == gin.DebugMode

	if isLocal {
		// Running in local development mode
		scheme = "http"
	} else {
		// Running in production or other mode
		scheme = "https"
	}

	// Get the host (domain) from the request
	host := c.Request.Host

	// Construct the base URL by combining the scheme and host
	baseURL := fmt.Sprintf("%s://%s", scheme, host)

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "registeristration successful",
		"data": map[string]interface{}{
			"id":             fmt.Sprint(user.ID),
			"email":          user.Email,
			"email_verified": user.EmailVerified,
		},
	})

	go sendVerificationEmail(user.FirstName, user.Email, baseURL)

}

func sendVerificationEmail(name, email, url string) {
	to := sendgrid.EmailAddress{
		Name:  name,
		Email: email,
	}

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
	client := sendgrid.NewClient(constant.SendGridApiKey, "hello@bonpay.finance", "Bonpay Finance", "Verify your email", messageBody)
	err = client.Send(&to)

	if err != nil {
		log.Printf("Error sending email: %v", err.Error())
	}
}

// VerifyEmail is a route handler that verifies the user's email address.
//
// This endpoint is used to verify the user's email address by providing the email and OTP token.
//
// @Summary Verify email address
// @Description Verifies the user's email address
// @Tags Authentication
// @Accept json
// @Produce json
// @Param email path string true "User's email address"
// @Param otp path string true "One-time password (OTP) token"
// @Success 302 {string} string "Redirects to the client URL with a jwt token"
// @Failure 302 {string} string "Redirects to the client URL with an error code"
// @Router /auth/verify/{email}/{otp} [get]
func (a *AuthHandler) VerifyEmail(c *gin.Context) {
	email := c.Param("email")
	token := c.Param("otp")

	user, userExist, err := a.userRepository.FindByCondition("email = ?", email)
	clientUrl := constant.ClientUrl

	if err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("%s/auth?error=500", clientUrl))
		return
	}
	if !userExist {
		c.Redirect(http.StatusFound, fmt.Sprintf("%s/auth?error=402", clientUrl))
		return
	}

	jwtToken, err := helpers.GenerateToken(constant.JWTSecretKey, user.Email, user.FirstName, user.ID.String())

	if err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("%s/auth?error=500", clientUrl))
		return
	}

	if user.EmailVerified {
		c.Redirect(http.StatusFound, fmt.Sprintf("%s/signin?&token=%s", clientUrl, jwtToken))
		return
	}

	valid := otp.OTPManage.VerifyOTP(email, token)

	if !valid {
		c.Redirect(http.StatusFound, fmt.Sprintf("%s/auth?error=401V", clientUrl))
		return
	}

	user.EmailVerified = true
	user.Status = string(models.ActiveAccount)
	user.UpdatedAt = time.Now()

	_, err = a.userRepository.UpdateUserByCondition("email", user.Email, user)

	if err != nil {
		c.Redirect(http.StatusFound, fmt.Sprintf("%s/auth?error=500", clientUrl))
		return
	}

	c.Redirect(http.StatusFound, fmt.Sprintf("%s/signin?token=%s", clientUrl, jwtToken))
}

func sendPasswordResetEmail(name, email string) {
	to := sendgrid.EmailAddress{
		Name:  name,
		Email: email,
	}

	otpToken, err := otp.OTPManage.GenerateOTP(email, time.Minute*10)
	if err != nil {
		log.Printf("Error generating otp: %v", err.Error())
	}

	type OTP struct {
		Otp  string
		Name string
		Url  string
	}

	messageBody, err := helpers.ParseTemplateFile("account_reset.html", OTP{Otp: otpToken, Name: name})

	if err != nil {
		log.Printf("Error sending email: %v", err.Error())
	}
	client := sendgrid.NewClient(constant.SendGridApiKey, "hello@bonpay.finance", "Bonpay Finance", "Reset your password", messageBody)
	err = client.Send(&to)

	if err != nil {
		log.Printf("Error sending email: %v", err.Error())
	}
}

// ForgotPassword is a route handler that sends the reset otp to the user's email address.
//
// This endpoint is used to send the otp to the user's email address by providing the email.
//
// @Summary Sends reset OTP
// @Description Sends the reset OTP to the user's email address
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body ForgotPasswordInput true "Input (email)"
// @Success 200 {string} string "Returns 'success' "
// @Failure 400 {string} string "Returns error message"
// @Router /auth/forgot-password [post]
func (a *AuthHandler) ForgotPassword(c *gin.Context) {
	var input ForgotPasswordInput

	validatedReqBody, exists := c.Get("validatedRequestBody")

	if !exists {
		helpers.Dispatch400Error(c, "Failed to retrieve validated request body", "")
		return
	}

	input, ok := validatedReqBody.(ForgotPasswordInput)

	if !ok {

		helpers.Dispatch400Error(c, "Failed to retrieve validated request body", "")
		return
	}

	userFound, err := a.findUserOrError(input.Email)

	if userFound == nil && err != nil {
		helpers.Dispatch200OK(c, "success", "success")
		return
	}

	var fullName string = userFound.FirstName + userFound.LastName
	var email string = userFound.Email

	if strings.Contains(fullName, " ") {
		fullName = strings.Split(fullName, " ")[1]
	}

	go sendPasswordResetEmail(fullName, email)

	helpers.Dispatch200OK(c, "success", "success")
}

// VerifyResetOTP is a route handler that verifies the user's email address.
//
// This endpoint is used to verify the user's email address by providing the email and OTP token.
//
// @Summary Verify email address
// @Description Verifies the user's email address
// @Tags Authentication
// @Accept json
// @Produce json
// @Param credentials body OtpVerifyInput true "Input (token and email)"
// @Success 200 {string} string "Returns 'success and JWT' "
// @Failure 400 {string} string "Returns error message"
// @Router /auth/forgot-password/verify/ [post]
func (a *AuthHandler) VerifyResetOTP(c *gin.Context) {
	var input OtpVerifyInput

	validatedReqBody, exists := c.Get("validatedRequestBody")

	if !exists {
		helpers.Dispatch400Error(c, "Failed to retrieve validated request body", "")
		return
	}

	input, ok := validatedReqBody.(OtpVerifyInput)

	if !ok {
		helpers.Dispatch400Error(c, "Failed to retrieve validated request body", "")
		return
	}

	user, userExist, err := a.userRepository.FindUserByCondition("email", input.Email)

	if err != nil {
		helpers.Dispatch400Error(c, "Error getting user", err)
		return
	}

	if !userExist {
		helpers.Dispatch404Error(c, "User not found", "")
		return
	}

	valid := otp.OTPManage.VerifyOTP(input.Email, input.Token)

	if !valid {
		c.JSON(http.StatusFound, gin.H{"error": "Invalid OTP"})
		return
	}

	jwtToken, err := helpers.GenerateToken(constant.JWTSecretKey, user.Email, user.FirstName, user.ID.String())

	if err != nil {
		helpers.Dispatch400Error(c, "Error generating token", err)
		return
	}

	helpers.Dispatch200OK(c, "success", jwtToken)

}

// ResetPassword is a route handler for resetting the user's password.
// It requires a valid JWT token and a JSON request body with new credentials.
//
// @Summary Reset the user's password
// @Description Reset the user's password using a JWT token and new credentials.
// @Tags Authentication
// @Accept json
// @Produce json
// @Param reset-token path string true "JWT token for resetting the password"
// @Param credentials body ResetPasswordInput true "New password and password confirmation"
// @Success 200 {string} string "Success: Password reset"
// @Failure 400 {string} string "Error: Invalid input or token"
// @Router /auth/reset-password/confirm/{reset-token} [post]
func (a *AuthHandler) ResetPassword(c *gin.Context) {
	resetToken := c.Params.ByName("reset-token")

	var input ResetPasswordInput

	validatedReqBody, exists := c.Get("validatedRequestBody")

	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve validated request body"})
		return
	}

	input, ok := validatedReqBody.(ResetPasswordInput)

	if !ok {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert request body to InputCreateUser"})
		return
	}

	token, err := jwt.ParseWithClaims(
		resetToken, &helpers.AuthTokenJwtClaim{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(constant.JWTSecretKey), nil
		})

	claims := token.Claims.(*helpers.AuthTokenJwtClaim)

	if err != nil {
		helpers.Dispatch400Error(c, "failed to get claims", nil)
		return
	}

	user, _, err := a.userRepository.FindUserByCondition("user_id", claims.UserId)

	if err != nil {
		helpers.Dispatch400Error(c, "failed to find user", nil)
		return
	}

	if user == nil {
		helpers.Dispatch400Error(c, "user not found", nil)
		return
	}

	if input.Password != input.PasswordConfirm {
		c.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	hashedPassword, _ := helpers.HashPassword(input.Password)

	// Update User in Database
	user.Password = hashedPassword
	user.EmailVerified = true
	user.Status = string(models.ActiveAccount)
	user.UpdatedAt = time.Now()

	_, err = a.userRepository.UpdateUserByCondition("email", user.Email, user)

	if err != nil {
		helpers.Dispatch400Error(c, err.Error(), err)
		return
	}

	c.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	c.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	c.SetCookie("logged_in", "", -1, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Password data updated successfully"})
}
