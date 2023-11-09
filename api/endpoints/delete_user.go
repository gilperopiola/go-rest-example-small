package endpoints

import (
	"errors"

	"github.com/gilperopiola/go-rest-example-small/api/common"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *handler) DeleteUser(c *gin.Context) {
	HandleRequest(c, h.makeDeleteUserRequest, h.deleteUser)
}

func (h *handler) makeDeleteUserRequest(c *gin.Context) (req *common.DeleteUserRequest, err error) {
	req = &common.DeleteUserRequest{UserID: c.GetInt(contextUserIDKey)}
	if req.UserID == 0 {
		return &common.DeleteUserRequest{}, common.ErrAllFieldsRequired
	}

	return req, nil
}

func (h *handler) deleteUser(c *gin.Context, request *common.DeleteUserRequest) (common.DeleteUserResponse, error) {
	user := request.ToUserModel()

	// Get user
	query := "(id = ?)"
	if err := h.db.Where(query, user.ID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.DeleteUserResponse{}, common.Wrap(err.Error(), common.ErrUserNotFound)
		}
		return common.DeleteUserResponse{}, common.Wrap(err.Error(), common.ErrGettingUser)
	}

	// If already deleted
	if user.Deleted {
		return common.DeleteUserResponse{}, common.Wrap("deleteUser: user.Deleted", common.ErrUserAlreadyDeleted)
	}

	// Delete user
	if err := h.db.Model(&user).Update("deleted", true).Error; err != nil {
		return common.DeleteUserResponse{}, common.Wrap(err.Error(), common.ErrDeletingUser)
	}

	return common.DeleteUserResponse{User: user.ToResponseModel()}, nil
}
