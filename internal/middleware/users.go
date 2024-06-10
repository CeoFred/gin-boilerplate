package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"

	"github.com/CeoFred/gin-boilerplate/constants"
	"github.com/CeoFred/gin-boilerplate/internal/helpers"
	"github.com/CeoFred/gin-boilerplate/internal/models"
	"github.com/CeoFred/gin-boilerplate/internal/repository"

	"strings"
)

type AppError struct {
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}
func NewError(message string) *AppError {
	return &AppError{
		Message: message,
	}
}

var (
	constant = constants.New()
)

func OnlyAdmin(db *gorm.DB, u *repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, _ := c.Get("claims")
		claimsData, ok := claims.(*helpers.AuthTokenJwtClaim)
		if !ok {
			helpers.ReturnJSON(c, "Something went wrong", nil, http.StatusUnauthorized)
			c.Abort()

			return
		}

		user, _, err := u.FindByCondition("id", claimsData.UserId)
		if err != nil {
			helpers.ReturnError(c, "Something went wrong", err, http.StatusUnauthorized)
			c.Abort()
			return
		}
		if user.Role != models.AdminRole {
			helpers.ReturnJSON(c, "Unauthorized access to resource", nil, http.StatusUnauthorized)
			c.Abort()
			return
		}

		c.Set("role", user.Role)
		c.Next()
	}
}

func JWTMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the JWT token from the Authorization header
		authHeader := c.GetHeader("Authorization")

		apiQuery := c.Query("access_token")

		if apiQuery != "" {
			authHeader = apiQuery
		}

		if authHeader == "" {
			helpers.ReturnJSON(c, "Missing Authorization Header or Access Token", nil, http.StatusUnauthorized)
			c.Abort()
			return
		}

		// Extract the token from the "Bearer <jwt>" format
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			helpers.ReturnJSON(c, "Invalid Authorization Header or Access Token", nil, http.StatusUnauthorized)
			c.Abort()

			return
		}

		// Parse and validate the JWT token
		token, err := jwt.ParseWithClaims(tokenString, &helpers.AuthTokenJwtClaim{}, func(token *jwt.Token) (interface{}, error) {
			// Provide the same JWT secret key used for signing the tokens
			return []byte(constant.JWTSecretKey), nil
		})
		if err != nil || !token.Valid {
			helpers.ReturnError(c, "Expired Authorization or Access Token", err, http.StatusUnauthorized)
			c.Abort()

			return
		}

		// Extract the claims from the token
		claims, ok := token.Claims.(*helpers.AuthTokenJwtClaim)
		if !ok {
			helpers.ReturnJSON(c, "Invalid claims", nil, http.StatusUnauthorized)
			c.Abort()

			return
		}

		// Attach the claims to the request context for further use
		c.Set("claims", claims)

		_, found, err := repository.NewUserRepository(db).FindByCondition("email", claims.Email)

		if err != nil {
			helpers.ReturnError(c, "Something went wrong", err, http.StatusUnauthorized)
			c.Abort()

			return
		}

		if !found {
			helpers.ReturnJSON(c, "Unauthorized access to resource", nil, http.StatusUnauthorized)
			c.Abort()

			return
		}

		// Proceed to the next middleware or route handler
		c.Next()
	}
}
