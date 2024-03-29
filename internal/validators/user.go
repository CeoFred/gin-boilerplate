package validators

import (
	"net/http"

	"github.com/CeoFred/gin-boilerplate/internal/handlers"
	"github.com/CeoFred/gin-boilerplate/internal/helpers"
	"github.com/CeoFred/gin-boilerplate/validator"

	"github.com/gin-gonic/gin"
)

func ValidateAccountResetScheme(c *gin.Context) {
	var body helpers.AccountReset
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		c.Abort()
		return
	}
	if err := validator.Validate(body); err != nil {
		helpers.Dispatch400Error(c, "input valudation failed", err)
		c.Abort()
		return
	}

	c.Next()
}

func ValidateOTPVerifySchema(c *gin.Context) {
	body := new(helpers.OtpVerify)
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid payload"})
		c.Abort()
		return
	}
	if err := validator.Validate(body); err != nil {
		helpers.Dispatch400Error(c, "input valudation failed", err)
		c.Abort()
		return
	}

	c.Next()
}

func ValidateRegisterUserSchema(c *gin.Context) {

	var body handlers.InputCreateUser
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid payload"})
		c.Abort()
		return
	}
	if err := validator.Validate(body); err != nil {
		helpers.Dispatch400Error(c, "input valudation failed", err)
		c.Abort()
		return
	}

	c.Set("validatedRequestBody", body)
	c.Next()
}

func ValidateLoginUser(c *gin.Context) {
	var body handlers.AuthenticateUser

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid payload"})
		c.Abort()
		return
	}
	if err := validator.Validate(body); err != nil {
		helpers.Dispatch400Error(c, "input valudation failed", err)
		c.Abort()
		return
	}

	c.Set("validatedRequestBody", body)
	c.Next()
}

func ValidateCompleteUserProfile(c *gin.Context) {
	body := new(helpers.UpdateAccountInformation)

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid payload"})
		c.Abort()
		return
	}
	if err := validator.Validate(body); err != nil {
		helpers.Dispatch400Error(c, "input valudation failed", err)
		c.Abort()
		return
	}

	c.Set("validatedRequestBody", body)
	c.Next()
}

func ValidateUpdateUserProfile(c *gin.Context) {
	body := new(handlers.UpdateUserProfileInput)

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid payload"})
		c.Abort()
		return
	}
	if err := validator.Validate(body); err != nil {
		helpers.Dispatch400Error(c, "input valudation failed", err)
		c.Abort()
		return
	}

	c.Set("validatedRequestBody", body)
	c.Next()
}

func ValidateResetUserSchema(c *gin.Context) {

	var body handlers.ForgotPasswordInput
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid payload"})
		c.Abort()
		return
	}
	if err := validator.Validate(body); err != nil {
		helpers.Dispatch400Error(c, "input valudation failed", err)
		c.Abort()
		return
	}

	c.Set("validatedRequestBody", body)
	c.Next()
}

func ValidateResetPasswordSchema(c *gin.Context) {

	var body handlers.ResetPasswordInput
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid payload"})
		c.Abort()
		return
	}
	if err := validator.Validate(body); err != nil {
		helpers.Dispatch400Error(c, "input valudation failed", err)
		c.Abort()
		return
	}

	c.Set("validatedRequestBody", body)
	c.Next()
}

func ValidateResetOTPVerifySchema(c *gin.Context) {
	var body handlers.OtpVerifyInput
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "invalid payload"})
		c.Abort()
		return
	}
	if err := validator.Validate(body); err != nil {
		helpers.Dispatch400Error(c, "input valudation failed", err)
		c.Abort()
		return
	}

	c.Set("validatedRequestBody", body)
	c.Next()
}
