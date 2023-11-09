package endpoints

import (
	"strings"

	"github.com/gilperopiola/go-rest-example-small/api/common"

	"github.com/gin-gonic/gin"
)

func (h *handler) CreateUser(c *gin.Context) {
	HandleRequest(c, h.makeCreateUserRequest, h.createUser)
}

func (h *handler) makeCreateUserRequest(c *gin.Context) (req *common.CreateUserRequest, err error) {

	if err = c.ShouldBindJSON(&req); err != nil {
		return &common.CreateUserRequest{}, common.Wrap(err.Error(), common.ErrBindingRequest)
	}

	if err = validateUsernameEmailAndPassword(req.Username, req.Email, req.Password); err != nil {
		return &common.CreateUserRequest{}, common.Wrap("makeCreateUserRequest", err)
	}

	return req, nil
}

func (h *handler) createUser(c *gin.Context, request *common.CreateUserRequest) (common.CreateUserResponse, error) {
	user := request.ToUserModel()
	user.Password = common.Hash(user.Password, h.config.Auth.HashSalt) // TODO this can be inside of the .ToUserModel fn?

	// Create user
	if err := h.db.Create(&user).Error; err != nil {
		errStr := err.Error()
		if strings.Contains(errStr, "Error 1062") { // Duplicate entry for key
			return common.CreateUserResponse{}, common.Wrap(errStr, common.ErrUsernameOrEmailAlreadyInUse)
		}
		return common.CreateUserResponse{}, common.Wrap(errStr, common.ErrCreatingUser)
	}

	return common.CreateUserResponse{User: user.ToResponseModel()}, nil
}
