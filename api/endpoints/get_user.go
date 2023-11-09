package endpoints

import (
	"errors"

	"github.com/gilperopiola/go-rest-example-small/api/common"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *handler) GetUser(c *gin.Context) {
	HandleRequest(c, h.makeGetUserRequest, h.getUser)
}

func (h *handler) makeGetUserRequest(c *gin.Context) (req *common.GetUserRequest, err error) {

	req = &common.GetUserRequest{UserID: c.GetInt(contextUserIDKey)}
	if req.UserID == 0 {
		return &common.GetUserRequest{}, common.ErrAllFieldsRequired
	}

	return req, nil
}

func (h *handler) getUser(c *gin.Context, request *common.GetUserRequest) (common.GetUserResponse, error) {
	user := request.ToUserModel()

	// Get user
	query := "(id = ?)"
	if err := h.db.Preload("Details").Preload("Posts").Where(query, user.ID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.GetUserResponse{}, common.Wrap(err.Error(), common.ErrUserNotFound)
		}
		return common.GetUserResponse{}, common.Wrap(err.Error(), common.ErrGettingUser)
	}

	// If deleted
	if user.Deleted {
		return common.GetUserResponse{}, common.Wrap("getUser: user.Deleted", common.ErrUserAlreadyDeleted)
	}

	return common.GetUserResponse{User: user.ToResponseModel()}, nil
}
