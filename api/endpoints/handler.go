package endpoints

import (
	"net/http"
	"regexp"

	"github.com/gilperopiola/go-rest-example-small/api/common"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler interface {
	HealthCheck(c *gin.Context)
	Signup(c *gin.Context)
	Login(c *gin.Context)
	CreateUser(c *gin.Context)
	GetUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	SearchUsers(c *gin.Context)
	ChangePassword(c *gin.Context)
	CreateUserPost(c *gin.Context)
}

type handler struct {
	config *common.Config
	db     *gorm.DB
	auth   *common.Auth
}

func NewHandler(config *common.Config, db *gorm.DB, auth *common.Auth) *handler {
	return &handler{
		db:     db,
		config: config,
		auth:   auth,
	}
}

func HandleRequest[req common.AllRequests, resp common.AllResponses](c *gin.Context, makeRequestFn func(*gin.Context) (req, error), serviceCallFn func(*gin.Context, req) (resp, error)) {

	// Build, validate and get request
	request, err := makeRequestFn(c)
	if err != nil {
		c.Error(err)
		return
	}

	// Call service with that request
	response, err := serviceCallFn(c, request)
	if err != nil {
		c.Error(err)
		return
	}

	// Return OK
	c.JSON(http.StatusOK, common.HTTPResponse{
		Success: true,
		Content: response,
	})
}

/*-----------------------
//       HELPERS
//---------------------*/

var (
	contextUserIDKey = "UserID"

	usernameMinLength = 4
	usernameMaxLength = 32
	passwordMinLength = 8
	passwordMaxLength = 64
	validEmailRegex   = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

func validateUsernameEmailAndPassword(username, email, password string) error {
	if email == "" || username == "" || password == "" {
		return common.ErrAllFieldsRequired
	}

	if !validEmailRegex.MatchString(email) {
		return common.ErrInvalidEmailFormat
	}

	if len(username) < usernameMinLength || len(username) > usernameMaxLength {
		return common.ErrInvalidUsernameLength(usernameMinLength, usernameMaxLength)
	}

	if len(password) < passwordMinLength || len(password) > passwordMaxLength {
		return common.ErrInvalidPasswordLength(passwordMinLength, passwordMaxLength)
	}

	return nil
}

/*--------------------
//       MISC
//-----------------*/

func (h handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, common.HTTPResponse{
		Success: true,
		Content: "service is up and running :)",
	})
}
