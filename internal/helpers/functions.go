package helpers

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"

	// "mime/multipart"
	"io"
	"os"
	"path/filepath"
	"time"

	"crypto/rand"
	"math/big"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/admin"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gin-gonic/gin"
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

func RandRange(min, max int) int {
	// Generate random bytes
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	// Convert bytes to integer within range
	n := int(uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24)
	return min + n%(max-min+1)
}

func CalculatePercentageChange(current, previous int) float64 {
	if previous == 0 {
		return 0
	}
	return (float64(current-previous) / float64(previous)) * 100
}

func GetFile(assetParam admin.AssetParams) (*admin.AssetResult, error) {
	env := constants.New()
	url := fmt.Sprintf("cloudinary://%s:%s@%s", env.CloudinaryAPIKey, env.CloudinaryApiSecret, env.CloudinaryName)
	cld, err := cloudinary.NewFromURL(url)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()

	res, err := cld.Admin.Asset(ctx, assetParam)
	if err != nil {
		return nil, err
	}

	return res, nil
}
func UploadFile(file io.Reader, filename string) (*uploader.UploadResult, error) {
	env := constants.New()
	url := fmt.Sprintf("cloudinary://%s:%s@%s", env.CloudinaryAPIKey, env.CloudinaryApiSecret, env.CloudinaryName)
	cld, err := cloudinary.NewFromURL(url)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	resp, err := cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID: filename,
		Folder:   "Motion365_User_Profile_Photo",
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func BaseURL(c *gin.Context) string {

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
	client := sendgrid.NewClient(constant.SendGridApiKey)
	err = client.Send(&to, constant.SendFromEmail, fromName, subject, messageBody)

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

func GenerateRandomNumber(length int) (int, error) {
	if length <= 0 {
		return 0, fmt.Errorf("length must be greater than 0")
	}

	min := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(length-1)), nil)
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(length)), nil)
	diff := new(big.Int).Sub(max, min)

	n, err := rand.Int(rand.Reader, diff)
	if err != nil {
		return 0, err
	}

	return int(n.Add(n, min).Int64()), nil
}

func ReturnJSON(c *gin.Context, message string, data interface{}, statusCode int) {
	c.Status(statusCode)
	c.JSON(statusCode, gin.H{
		"status":  statusCode <= 201,
		"message": message,
		"data":    data,
	})
}

func ReturnError(c *gin.Context, message string, err error, status int) {
	c.JSON(status, gin.H{
		"message": message,
		"error":   err.Error(),
		"status":  false,
	})
	log.Println("error: ", err.Error())
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 4)
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
			IssuedAt:  time.Now().Unix(),
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
