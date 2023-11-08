package endpoints

import (
	"github.com/gilperopiola/go-rest-example/api/common"
	"github.com/gin-gonic/gin"
)

func (h *handler) SearchUsers(c *gin.Context) {
	HandleRequest(c, h.makeSignupRequest, h.signup)
}

func (h *handler) searchUsers(c *gin.Context, request *common.SignupRequest) (common.SignupResponse, error) {
	return common.SignupResponse{}, nil
}

func (h *handler) makeSearchUsersRequest(c *gin.Context) (req *common.SignupRequest, err error) {

	if err = c.ShouldBindJSON(&req); err != nil {
		return &common.SignupRequest{}, common.Wrap(err.Error(), common.ErrBindingRequest)
	}

	if err = validateUsernameEmailAndPassword(req.Username, req.Email, req.Password); err != nil {
		return &common.SignupRequest{}, common.Wrap("makeSignupRequest", err)
	}

	if req.Password != req.RepeatPassword {
		return &common.SignupRequest{}, common.Wrap("makeSignupRequest", common.ErrPasswordsDontMatch)
	}

	return req, nil
}
