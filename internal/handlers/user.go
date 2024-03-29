package handlers

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/gin-gonic/gin"

	"github.com/CeoFred/gin-boilerplate/constants"
	"github.com/CeoFred/gin-boilerplate/internal/helpers"
	"github.com/CeoFred/gin-boilerplate/internal/models"
	"github.com/CeoFred/gin-boilerplate/internal/repository"
)

var ()

type UserHandler struct {
	userRepository *repository.UserRepository
}

func NewUserHandler(userRepo *repository.UserRepository,
) *UserHandler {
	return &UserHandler{
		userRepository: userRepo,
	}
}

type FileUploadResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type UserProfile struct {
	Email     string    `json:"email"`
	Userid    string    `json:"userid"`
	Name      string    `json:"name"`
	Country   string    `json:"country"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type UpdateUserProfileInput struct {
	PhoneNumber string `json:"phone_number" validate:"required"`
}

// UpdateUserProfile is a route handler that handles updating the user profile
//
// # This endpoint is used to update the user profile
//
// @Summary Update user profile
// @Description Updates some details about the user
// @Tags User
// @Accept json
// @Produce json
// @Param credentials body UpdateUserProfileInput true "update user profile"
// @Security BearerAuth
// @Success 200 {object} SuccessResponse
// @Failure 401 {object} ErrorResponse
// @Router /user [put]
func (u *UserHandler) UpdateUserProfile(c *gin.Context) {
	var input UpdateUserProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		helpers.Dispatch500Error(c, err)
		return
	}

	claimsRaw, exists := c.Get("claims")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	claims, ok := claimsRaw.(*helpers.AuthTokenJwtClaim)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, found, err := u.userRepository.FindByCondition("email = ?", claims.Email)
	if err != nil {
		helpers.Dispatch400Error(c, "something went wrong", err)
		return
	}
	if !found {
		helpers.Dispatch400Error(c, "user not found", nil)
		return
	}

	user.PhoneNumber = input.PhoneNumber

	_, err = u.userRepository.UpdateUserByCondition("email", user.Email, user)
	if err != nil {
		helpers.Dispatch400Error(c, "something went wrong", err)
		return
	}

	helpers.Dispatch200OK(c, "OK", nil)
}

// UserProfile is a route handler that retrieves the user profile of the authenticated user.
//
// This endpoint is used to get the profile information of the authenticated user based on the JWT claims.
//
// @Summary Get user profile
// @Description Retrieves the profile information of the authenticated user
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} User
// @Failure 401 {object} ErrorResponse
// @Router /user/profile [get]
func (u *UserHandler) UserProfile(c *gin.Context) {

	claimsRaw, exists := c.Get("claims")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, ok := claimsRaw.(*helpers.AuthTokenJwtClaim)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	usr, f, err := u.userRepository.FindByCondition("user_id = ?", user.UserId)

	if err != nil {
		helpers.Dispatch400Error(c, "failed to fetch wallet", nil)
		return
	}

	if !f {
		helpers.Dispatch400Error(c, "failed to fetch user", err)
		return
	}

	helpers.Dispatch200OK(c, "success", models.User{
		CreatedAt:     usr.CreatedAt,
		Email:         usr.Email,
		Status:        usr.Status,
		EmailVerified: usr.EmailVerified,
		LastLogin:     usr.LastLogin,
		Role:          usr.Role,
	})

}

func uploadFile(file *multipart.FileHeader) (resp *uploader.UploadResult, err error) {

	env_ := constants.New()
	// Open the uploaded file
	fileOpened, err := file.Open()
	if err != nil {
		// Handle error
		return nil, err
	}

	defer fileOpened.Close()
	url := fmt.Sprintf("cloudinary://%s:%s@%s", env_.CloudinaryAPIKey, env_.CloudinaryApiSecret, env_.CloudinaryName)

	cld, err := cloudinary.NewFromURL(url)

	if err != nil {
		return nil, err
	}
	var ctx = context.Background()
	// Upload the image to Cloudinary
	resp, err = cld.Upload.Upload(ctx, fileOpened, uploader.UploadParams{PublicID: file.Filename,
		Folder: "BonpayFiance",
	})

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// NotFound returns custom 404 page
func NotFound(c *gin.Context) {
	c.Status(404)
	c.File("./static/private/404.html")
}
