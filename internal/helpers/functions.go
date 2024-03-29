package helpers

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"crypto/rand"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	"github.com/CeoFred/gin-boilerplate/constants"
	"github.com/CeoFred/gin-boilerplate/sendgrid"
)

const (
	letterBytes  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" // Characters to choose from
	randomLength = 10                                                               // Length of the random string
)

var (
	constant = constants.New()
)

func SendEmail(subject, fromName, template, name, email string, templateData interface{}) {
	messageBody, err := ParseTemplateFile(template+".html", templateData)

	if err != nil {
		log.Printf("Error sending email: %v", err.Error())
		return
	}
	to := sendgrid.EmailAddress{
		Name:  name,
		Email: email,
	}
	client := sendgrid.NewClient(constant.SendGridApiKey, "hello@bonpay.finance", fromName, subject, messageBody)
	err = client.Send(&to)

	if err != nil {
		log.Printf("Error sending email: %v", err.Error())
	}
}

func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	for i, b := range bytes {
		bytes[i] = letterBytes[b%byte(len(letterBytes))]
	}

	return string(bytes), nil
}

func Dispatch200OK(c *gin.Context, message string, data interface{}) {
	c.Status(http.StatusOK)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

func Dispatch201Created(c *gin.Context, message string, data interface{}) {
	c.Status(http.StatusCreated)
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": message,
		"data":    data,
	})
}

// 500 - internal server error
func Dispatch500Error(c *gin.Context, err error) {
	c.Status(http.StatusInternalServerError)
	c.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"message": fmt.Sprintf("%v", err),
		"data":    nil,
	})
}

// 400 - bad request
func Dispatch400Error(c *gin.Context, msg string, err any) {
	c.Status(http.StatusBadRequest)
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"message": msg,
		"data":    err,
	})
}

// 404 - not found
func Dispatch404Error(c *gin.Context, msg string, err any) {
	c.Status(http.StatusNotFound)
	c.JSON(http.StatusOK, gin.H{
		"success": false,
		"message": msg,
		"data":    err,
	})
}

func SchemaError(c *gin.Context, err error) {
	var errors []*IError
	for _, err := range err.(validator.ValidationErrors) {
		var el IError
		el.Field = err.Field()
		el.Tag = err.Tag()
		el.Value = err.Param()
		errors = append(errors, &el)
	}
	Dispatch400Error(c, "invalid body schema", errors)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// GenerateToken generates a jwt token
func GenerateToken(JWTSecretKey, email, name, userid string) (signedToken string, err error) {
	claims := &AuthTokenJwtClaim{
		Email:  email,
		Name:   name,
		UserId: userid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * 24 * 20).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err = token.SignedString([]byte(JWTSecretKey))
	if err != nil {
		return
	}
	return
}

func ParseTemplateFile(filename string, mapping interface{}) (string, error) {
	absolutePath, err := filepath.Abs("templates/email/" + filename)
	if err != nil {
		return "", err
	}

	content, err := os.ReadFile(filepath.Clean(absolutePath))
	if err != nil {
		return "", err
	}

	temp, err := template.New("emailTemplate").Parse(string(content))
	if err != nil {
		return "", err
	}
	messageBody := new(bytes.Buffer)
	err = temp.Execute(messageBody, mapping)
	if err != nil {
		return "", err
	}

	return messageBody.String(), nil
}

func TimeNow(timezone string) (string, error) {

	location, err := time.LoadLocation(timezone)
	if err != nil {
		return "", err
	}

	currentTime := time.Now().In(location)
	return currentTime.String(), nil
}

type AppError struct {
	message string
}

func (e AppError) Error() string {
	return e.message
}

func NewError(message string) *AppError {
	return &AppError{message: message}
}

func GetBaseURL(c *gin.Context) string {
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
	return baseURL
}

func GetAuthenticatedUser(c *gin.Context) (*AuthTokenJwtClaim, error) {

	var claims *AuthTokenJwtClaim

	user, claims_exists := c.Get("claims")

	if !claims_exists {
		return nil, NewError("Failed to retrieve claims")
	}

	claims, ok := user.(*AuthTokenJwtClaim)

	if !ok {
		return nil, NewError("Failed to convert user claims")
	}

	return claims, nil
}
