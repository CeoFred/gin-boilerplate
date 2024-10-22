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
	bootstrap "github.com/CeoFred/gin-boilerplate/internal/bootstrap"
	"github.com/CeoFred/gin-boilerplate/internal/helpers"
)

type UserHandler struct {
	deps *bootstrap.AppDependencies
}

func NewUserHandler(deps *bootstrap.AppDependencies,
) *UserHandler {
	return &UserHandler{
		deps: deps,
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

	validatedReqBody, exists := c.Get("validatedRequestBody")

	if !exists {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf(helpers.INVALID_REQUEST_BODY), http.StatusBadRequest)
		return
	}

	input, ok := validatedReqBody.(UpdateUserProfileInput)
	if !ok {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf(helpers.REQUEST_BODY_PARSE_ERROR), http.StatusBadRequest)
		return
	}

	claims, err := helpers.GetAuthenticatedUser(c)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	user, found, err := u.deps.UserRepo.FindByCondition("email = ?", claims.Email)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	if !found {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf("user not found"), http.StatusNotFound)
		return
	}

	user.PhoneNumber = input.PhoneNumber

	_, err = u.deps.UserRepo.Save(user)
	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	helpers.ReturnJSON(c, "Profile updated successfully", user, http.StatusOK)
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
// @Success 200 {object} models.User
// @Failure 401 {object} ErrorResponse
// @Router /user/profile [get]
func (u *UserHandler) UserProfile(c *gin.Context) {

	claimsRaw, exists := c.Get("claims")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	authClaims, ok := claimsRaw.(*helpers.AuthTokenJwtClaim)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, f, err := u.deps.UserRepo.FindByCondition("user_id = ?", authClaims.UserId)

	if err != nil {
		helpers.ReturnError(c, "Something went wrong", err, http.StatusInternalServerError)
		return
	}

	if !f {
		helpers.ReturnError(c, "Something went wrong", fmt.Errorf("account not found"), http.StatusNotFound)
		return
	}

	helpers.ReturnJSON(c, "Profile retrieved", user, http.StatusOK)

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
