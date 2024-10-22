package routes

import (
	deps "github.com/CeoFred/gin-boilerplate/internal/bootstrap"
	"github.com/CeoFred/gin-boilerplate/internal/handlers"
	"github.com/CeoFred/gin-boilerplate/internal/middleware"
	"github.com/CeoFred/gin-boilerplate/internal/validators"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(router *gin.RouterGroup, d *deps.AppDependencies) {
	userRouter := router.Group("user")

	handler := handlers.NewUserHandler(d)

	userRouter.GET("/profile", middleware.JWTMiddleware(d.DatabaseService), handler.UserProfile)
	userRouter.PUT("/", middleware.JWTMiddleware(d.DatabaseService), validators.ValidateUpdateUserProfile, handler.UpdateUserProfile)
}
