package endpoints

import (
	"errors"

	"github.com/gilperopiola/go-rest-example-small/api/common"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *handler) ChangePassword(c *gin.Context) {
	HandleRequest(c, h.makeChangePasswordRequest, h.changePassword)
}

func (h *handler) makeChangePasswordRequest(c *gin.Context) (req common.ChangePasswordRequest, err error) {

	if err = c.ShouldBindJSON(&req); err != nil {
		return common.ChangePasswordRequest{}, common.Wrap(err.Error(), common.ErrBindingRequest)
	}

	req.UserID = c.GetInt(contextUserIDKey)
	if req.UserID == 0 {
		return common.ChangePasswordRequest{}, common.ErrAllFieldsRequired
	}

	if req.OldPassword == "" || req.NewPassword == "" || req.RepeatPassword == "" {
		return common.ChangePasswordRequest{}, common.ErrAllFieldsRequired
	}

	if len(req.NewPassword) < passwordMinLength || len(req.NewPassword) > passwordMaxLength {
		return common.ChangePasswordRequest{}, common.ErrInvalidPasswordLength(passwordMinLength, passwordMaxLength)
	}

	if req.NewPassword != req.RepeatPassword {
		return common.ChangePasswordRequest{}, common.ErrPasswordsDontMatch
	}

	return req, nil
}

func (h *handler) changePassword(c *gin.Context, request common.ChangePasswordRequest) (common.ChangePasswordResponse, error) {
	user := request.ToUserModel()

	// Get user
	query := "(id = ? OR username = ? OR email = ?) AND deleted = false"
	if err := h.db.Where(query, user.ID, user.Username, user.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.ChangePasswordResponse{}, common.Wrap(err.Error(), common.ErrUserNotFound)
		}
		return common.ChangePasswordResponse{}, common.Wrap(err.Error(), common.ErrGettingUser)
	}

	// Check if old password matches
	if user.Password != common.Hash(request.OldPassword, h.config.HashSalt) {
		return common.ChangePasswordResponse{}, common.Wrap("changePassword: user.Password != common.Hash", common.ErrWrongPassword)
	}

	// Generate new hashed password
	newPassword := common.Hash(request.NewPassword, h.config.HashSalt)

	// Update password
	if err := h.db.Model(&common.User{}).Where("id = ?", user.ID).Update("password", newPassword).Error; err != nil {
		return common.ChangePasswordResponse{}, common.Wrap(err.Error(), common.ErrUpdatingUser)
	}

	return common.ChangePasswordResponse{User: user.ToResponseModel()}, nil
}
