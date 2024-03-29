package routes

import (
	"github.com/CeoFred/gin-boilerplate/internal/handlers"
	"github.com/CeoFred/gin-boilerplate/internal/middleware"
	"github.com/CeoFred/gin-boilerplate/internal/repository"
	"github.com/CeoFred/gin-boilerplate/internal/validators"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

func RegisterUserRoutes(router *gin.RouterGroup, db *gorm.DB) {
	userRouter := router.Group("user")

	handler := handlers.NewUserHandler(repository.NewUserRepository(db))

	userRouter.GET("/profile", middleware.JWTMiddleware(db), handler.UserProfile)
	userRouter.PUT("/", middleware.JWTMiddleware(db), validators.ValidateUpdateUserProfile, handler.UpdateUserProfile)
}
