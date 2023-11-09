package endpoints

import (
	"errors"
	"strings"

	"github.com/gilperopiola/go-rest-example-small/api/common"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *handler) UpdateUser(c *gin.Context) {
	HandleRequest(c, h.makeUpdateUserRequest, h.updateUser)
}

func (h *handler) makeUpdateUserRequest(c *gin.Context) (req *common.UpdateUserRequest, err error) {

	if err = c.ShouldBindJSON(&req); err != nil {
		return &common.UpdateUserRequest{}, common.Wrap(err.Error(), common.ErrBindingRequest)
	}

	req.UserID = c.GetInt(contextUserIDKey)

	if req.UserID == 0 || (req.Email == "" && req.Username == "" && req.FirstName == nil && req.LastName == nil) {
		return &common.UpdateUserRequest{}, common.ErrAllFieldsRequired
	}

	if req.Email != "" && !validEmailRegex.MatchString(req.Email) {
		return &common.UpdateUserRequest{}, common.ErrInvalidEmailFormat
	}

	if req.Username != "" {
		if len(req.Username) < usernameMinLength || len(req.Username) > usernameMaxLength {
			return &common.UpdateUserRequest{}, common.ErrInvalidUsernameLength(usernameMinLength, usernameMaxLength)
		}
	}

	return req, nil
}

func (h *handler) updateUser(c *gin.Context, request *common.UpdateUserRequest) (common.UpdateUserResponse, error) {
	user := request.ToUserModel()

	// Get user
	query := "(id = ? OR username = ? OR email = ?) AND deleted = false"
	if err := h.db.Preload("Details").Where(query, user.ID, user.Username, user.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.UpdateUserResponse{}, common.Wrap(err.Error(), common.ErrUserNotFound)
		}
		return common.UpdateUserResponse{}, common.Wrap(err.Error(), common.ErrGettingUser)
	}

	// Overwrite fields that aren't empty
	if request.Username != "" {
		user.Username = request.Username
	}
	if request.Email != "" {
		user.Email = request.Email
	}
	if request.FirstName != nil {
		user.Details.FirstName = *request.FirstName
	}
	if request.LastName != nil {
		user.Details.LastName = *request.LastName
	}

	// Update user
	if err := h.db.Omit("Details").Save(&user).Error; err != nil {
		if strings.Contains(err.Error(), "Error 1062") { // Duplicate entry for key
			return common.UpdateUserResponse{}, common.Wrap(err.Error(), common.ErrUsernameOrEmailAlreadyInUse)
		}
		return common.UpdateUserResponse{}, common.Wrap(err.Error(), common.ErrUpdatingUser)
	}

	// Update user details
	if user.Details.ID != 0 {
		if err := h.db.Save(&user.Details).Error; err != nil {
			return common.UpdateUserResponse{}, common.Wrap(err.Error(), common.ErrUpdatingUserDetail)
		}
	}

	return common.UpdateUserResponse{User: user.ToResponseModel()}, nil
}
