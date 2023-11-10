package endpoints

import (
	"strings"

	"github.com/gilperopiola/go-rest-example-small/api/common"

	"github.com/gin-gonic/gin"
)

func (h *handler) Signup(c *gin.Context) {
	HandleRequest(c, h.makeSignupRequest, h.signup)
}

func (h *handler) makeSignupRequest(c *gin.Context) (req common.SignupRequest, err error) {

	if err = c.ShouldBindJSON(&req); err != nil {
		return common.SignupRequest{}, common.Wrap(err.Error(), common.ErrBindingRequest)
	}

	if err = validateUsernameEmailAndPassword(req.Username, req.Email, req.Password); err != nil {
		return common.SignupRequest{}, common.Wrap("makeSignupRequest", err)
	}

	if req.Password != req.RepeatPassword {
		return common.SignupRequest{}, common.Wrap("makeSignupRequest", common.ErrPasswordsDontMatch)
	}

	return req, nil
}

func (h *handler) signup(c *gin.Context, request common.SignupRequest) (common.SignupResponse, error) {
	user := request.ToUserModel()
	user.HashPassword(h.config.HashSalt)

	// Create user
	if err := h.db.Create(&user).Error; err != nil {
		if strings.Contains(err.Error(), "Error 1062") { // Duplicate entry for key
			return common.SignupResponse{}, common.Wrap(err.Error(), common.ErrUsernameOrEmailAlreadyInUse)
		}
		return common.SignupResponse{}, common.Wrap(err.Error(), common.ErrCreatingUser)
	}

	return common.SignupResponse{User: user.ToResponseModel()}, nil
}
