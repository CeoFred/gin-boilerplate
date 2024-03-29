package middleware

import (
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
			c.JSON(500, NewError("Internal Server Error"))
			c.Abort()
			return
		}

		user, _, err := u.FindByCondition("user_id", claimsData.UserId)
		if err != nil {
			helpers.Dispatch500Error(c, err)
			c.Abort()
			return
		}
		if user.Role != models.AdminRole {
			c.JSON(500, NewError("Unauthorized access"))
			c.Abort()
			return
		}
		c.Next()
	}
}

func JWTMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the JWT token from the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, NewError("Unauthorized"))
			c.Abort()
			return
		}

		// Extract the token from the "Bearer <jwt>" format
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			c.JSON(401, NewError("Unauthorized"))
			c.Abort()
			return
		}

		// Parse and validate the JWT token
		token, err := jwt.ParseWithClaims(tokenString, &helpers.AuthTokenJwtClaim{}, func(token *jwt.Token) (interface{}, error) {
			// Provide the same JWT secret key used for signing the tokens
			return []byte(constant.JWTSecretKey), nil
		})
		if err != nil || !token.Valid {
			c.JSON(403, NewError("Unauthorized"))
			c.Abort()
			return
		}

		// Extract the claims from the token
		claims, ok := token.Claims.(*helpers.AuthTokenJwtClaim)
		if !ok {
			c.JSON(403, NewError("Unauthorized"))
			c.Abort()
			return
		}

		// Attach the claims to the request context for further use
		c.Set("claims", claims)

		_, found, err := repository.NewUserRepository(db).FindByCondition("email", claims.Email)

		if err != nil {
			c.JSON(401, NewError("Unauthorized"))
			c.Abort()
			return
		}

		if !found {
			c.JSON(401, NewError("Unauthorized"))
			c.Abort()
			return
		}

		// Proceed to the next middleware or route handler
		c.Next()
	}
}
