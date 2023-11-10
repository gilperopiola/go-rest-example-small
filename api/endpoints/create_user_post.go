package endpoints

import (
	"github.com/gilperopiola/go-rest-example-small/api/common"

	"github.com/gin-gonic/gin"
)

func (h *handler) CreateUserPost(c *gin.Context) {
	HandleRequest(c, h.makeCreateUserPostRequest, h.createUserPost)
}

func (h *handler) makeCreateUserPostRequest(c *gin.Context) (req common.CreateUserPostRequest, err error) {

	if err = c.ShouldBindJSON(&req); err != nil {
		return common.CreateUserPostRequest{}, common.Wrap(err.Error(), common.ErrBindingRequest)
	}

	req.UserID = c.GetInt(contextUserIDKey)
	if req.UserID == 0 || req.Title == "" {
		return common.CreateUserPostRequest{}, common.ErrAllFieldsRequired
	}

	return req, nil
}

func (h *handler) createUserPost(c *gin.Context, request common.CreateUserPostRequest) (common.CreateUserPostResponse, error) {
	userPost := request.ToUserPostModel()

	if err := h.db.Create(&userPost).Error; err != nil {
		return common.CreateUserPostResponse{}, common.Wrap(err.Error(), common.ErrCreatingUserPost)
	}

	return common.CreateUserPostResponse{UserPost: userPost.ToResponseModel()}, nil
}
