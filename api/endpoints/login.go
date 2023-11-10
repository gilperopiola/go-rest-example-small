package endpoints

import (
	"errors"

	"github.com/gilperopiola/go-rest-example-small/api/common"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *handler) Login(c *gin.Context) {
	HandleRequest(c, h.makeLoginRequest, h.login)
}

func (h *handler) makeLoginRequest(c *gin.Context) (req common.LoginRequest, err error) {

	if err = c.ShouldBindJSON(&req); err != nil {
		return common.LoginRequest{}, common.Wrap(err.Error(), common.ErrBindingRequest)
	}

	if req.UsernameOrEmail == "" || req.Password == "" {
		return common.LoginRequest{}, common.Wrap("makeLoginRequest", common.ErrAllFieldsRequired)
	}

	return req, nil
}

func (h *handler) login(c *gin.Context, request common.LoginRequest) (common.LoginResponse, error) {
	user := request.ToUserModel()

	// Get user
	query := "(username = ? OR email = ?) AND deleted = false"
	if err := h.db.Where(query, user.Username, user.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common.LoginResponse{}, common.Wrap(err.Error(), common.ErrUserNotFound)
		}
		return common.LoginResponse{}, common.Wrap(err.Error(), common.ErrGettingUser)
	}

	// Check password
	if user.Password != common.Hash(request.Password, h.config.HashSalt) {
		return common.LoginResponse{}, common.Wrap("login: user.Password != common.Hash", common.ErrWrongPassword)
	}

	// Generate token
	tokenString, err := h.auth.GenerateToken(user.ID, user.Username, user.Email, user.GetRole())
	if err != nil {
		return common.LoginResponse{}, common.Wrap("login: auth.GenerateToken", common.ErrUnauthorized)
	}

	return common.LoginResponse{Token: tokenString}, nil
}
