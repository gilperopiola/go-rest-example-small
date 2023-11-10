package endpoints

import (
	"strconv"

	"github.com/gilperopiola/go-rest-example-small/api/common"

	"github.com/gin-gonic/gin"
)

func (h *handler) SearchUsers(c *gin.Context) {
	HandleRequest(c, h.makeSearchUsersRequest, h.searchUsers)
}

func (h *handler) makeSearchUsersRequest(c *gin.Context) (req common.SearchUsersRequest, err error) {

	defaultPage := "0"
	defaultPerPage := "10"

	req.Username = c.Query("username")

	req.Page, err = strconv.Atoi(c.DefaultQuery("page", defaultPage))
	if err != nil {
		return common.SearchUsersRequest{}, common.ErrInvalidValue("page")
	}

	req.PerPage, err = strconv.Atoi(c.DefaultQuery("per_page", defaultPerPage))
	if err != nil {
		return common.SearchUsersRequest{}, common.ErrInvalidValue("per_page")
	}

	if req.Page < 0 || req.PerPage <= 0 {
		return common.SearchUsersRequest{}, common.ErrAllFieldsRequired
	}

	return req, nil
}

func (h *handler) searchUsers(c *gin.Context, request common.SearchUsersRequest) (common.SearchUsersResponse, error) {
	var (
		users   common.Users
		page    = request.Page
		perPage = request.PerPage
	)

	query := h.db.Preload("Details").Where("username LIKE ?", "%"+request.Username+"%")
	if err := query.Offset(page * perPage).Limit(perPage).Find(&users).Error; err != nil {
		return common.SearchUsersResponse{}, common.Wrap(err.Error(), common.ErrSearchingUsers)
	}

	return common.SearchUsersResponse{
		Users:   users.ToResponseModel(),
		Page:    page,
		PerPage: perPage,
	}, nil
}
