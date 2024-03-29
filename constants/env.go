package constants

import (
	"log"
	"os"
	// "regexp"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                   string
	Env                    string
	ProjectID              string
	DbHost                 string
	DbUser                 string
	DbPassword             string
	DbName                 string
	DbPort                 string
	JWTSecretKey           string
	GoogleClientID         string
	GoogleClientSecret     string
	GithubClientID         string
	GithubClientSecret     string
	OAuthRedirectBaseURL   string
	ClientOauthRedirectURL string
	SendGridApiKey         string
	CloudinaryAPIKey       string
	CloudinaryApiSecret    string
	CloudinaryName         string
	ClientUrl              string
	APIToolkitKey          string
}

func init() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error loading .env file")
	}

}

func New() *Config {
	return &Config{
		DbHost:                 getEnv("POSTGRES_HOST", ""),
		DbUser:                 getEnv("POSTGRES_USER", ""),
		DbPassword:             getEnv("POSTGRES_PASSWORD", ""),
		DbName:                 getEnv("POSTGRES_NAME", ""),
		DbPort:                 getEnv("POSTGRES_PORT", ""),
		Port:                   getEnv("PORT", ""),
		JWTSecretKey:           getEnv("JWT_SCECRET", ""),
		GoogleClientID:         getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret:     getEnv("GOOGLE_CLIENT_SECRET", ""),
		GithubClientID:         getEnv("GITHUB_CLIENT_ID", ""),
		GithubClientSecret:     getEnv("GITHUB_CLIENT_SECRET", ""),
		OAuthRedirectBaseURL:   getEnv("OAUTH_REDIRECT_BASE_URL", ""),
		ClientOauthRedirectURL: getEnv("CLIENT_OAUTH_REDIRECT_URL", ""),
		SendGridApiKey:         getEnv("SENDGRID_API_KEY", ""),
		CloudinaryAPIKey:       getEnv("CLOUDINARY_API_KEY", ""),
		CloudinaryApiSecret:    getEnv("CLOUDINARY_API_SECRET", ""),
		CloudinaryName:         getEnv("CLOUDINARY_NAME", ""),
		ClientUrl:              getEnv("CLIENT_WEBAPP_URL", ""),
		APIToolkitKey:          getEnv("API_TOOLKIT_KEY", ""),
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
