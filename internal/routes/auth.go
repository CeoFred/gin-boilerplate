package routes

import (
	"github.com/CeoFred/gin-boilerplate/internal/handlers"
	"github.com/CeoFred/gin-boilerplate/internal/repository"
	"github.com/CeoFred/gin-boilerplate/internal/validators"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAuthRoutes(router *gin.RouterGroup, db *gorm.DB) {
	userRepo := repository.NewUserRepository(db)

	handler := handlers.NewAuthHandler(userRepo)

	authRouter := router.Group("/auth")

	authRouter.POST("/register", validators.ValidateRegisterUserSchema, handler.Register)
	authRouter.POST("/login", validators.ValidateLoginUser, handler.Authenticate)
	authRouter.GET("/verify/:email/:otp", handler.VerifyEmail)

	authRouter.POST("/forgot-password/verify", validators.ValidateResetOTPVerifySchema, handler.VerifyResetOTP)
	authRouter.POST("/forgot-password", validators.ValidateResetUserSchema, handler.ForgotPassword)
	authRouter.POST("/reset-password/confirm/:reset-token", validators.ValidateResetPasswordSchema, handler.ResetPassword)

}
