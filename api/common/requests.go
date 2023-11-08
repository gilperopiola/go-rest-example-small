package common

import (
	"regexp"
	"time"
)

type All interface {
	*SignupRequest |
		*LoginRequest |
		*CreateUserRequest |
		*GetUserRequest |
		*UpdateUserRequest |
		*DeleteUserRequest |
		*SearchUsersRequest |
		*ChangePasswordRequest |
		*CreateUserPostRequest
}

/*---------------
//    SIGNUP
//-------------*/

type SignupRequest struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	RepeatPassword string `json:"repeat_password"`

	// User Detail
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (r *SignupRequest) ToUserModel() User {
	return User{
		Username: r.Username,
		Email:    r.Email,
		Password: r.Password,
		Deleted:  false,
		Details: UserDetail{
			FirstName: r.FirstName,
			LastName:  r.LastName,
		},
		IsAdmin:   false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

/*--------------
//    LOGIN
//------------*/

type LoginRequest struct {
	UsernameOrEmail string `json:"username_or_email"`
	Password        string `json:"password"`
}

func (r *LoginRequest) ToUserModel() User {
	user := User{Password: r.Password}

	if validEmailRegex.MatchString(r.UsernameOrEmail) {
		user.Email = r.UsernameOrEmail
	} else {
		user.Username = r.UsernameOrEmail
	}

	return user
}

var (
	contextUserIDKey = "UserID"
	pathUserIDKey    = "user_id"

	validEmailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

/*---------------------
//    CREATE USER
--------------------*/

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`

	// User Detail
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (r *CreateUserRequest) ToUserModel() User {
	return User{
		Email:    r.Email,
		Username: r.Username,
		Password: r.Password,
		Deleted:  false,
		Details: UserDetail{
			FirstName: r.FirstName,
			LastName:  r.LastName,
		},
		IsAdmin:   r.IsAdmin,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

/*--------------------
//     GET USER
//------------------*/

type GetUserRequest struct {
	UserID int `json:"user_id"`
}

func (r *GetUserRequest) ToUserModel() User {
	return User{ID: r.UserID}
}

/*--------------------
//    UPDATE USER
//------------------*/

type UpdateUserRequest struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`

	// User Detail
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
}

func (r *UpdateUserRequest) ToUserModel() User {
	firstName, lastName := "", ""
	if r.FirstName != nil {
		firstName = *r.FirstName
	}
	if r.LastName != nil {
		lastName = *r.LastName
	}

	return User{
		ID:       r.UserID,
		Username: r.Username,
		Email:    r.Email,
		Details: UserDetail{
			FirstName: firstName,
			LastName:  lastName,
		},
	}
}

/*--------------------
//    DELETE USER
//------------------*/

type DeleteUserRequest struct {
	UserID int `json:"user_id"`
}

func (r *DeleteUserRequest) ToUserModel() User {
	return User{ID: r.UserID}
}

/*--------------------
//    SEARCH USERS
//------------------*/

type SearchUsersRequest struct {
	Username string `json:"username"`
	Page     int    `json:"page"`
	PerPage  int    `json:"per_page"`
}

func (r *SearchUsersRequest) ToUserModel() User {
	return User{Username: r.Username}
}

/*-----------------------
//    CHANGE PASSWORD
//---------------------*/

type ChangePasswordRequest struct {
	UserID         int    `json:"user_id"`
	OldPassword    string `json:"old_password"`
	NewPassword    string `json:"new_password"`
	RepeatPassword string `json:"repeat_password"`
}

func (r *ChangePasswordRequest) ToUserModel() User {
	return User{
		ID:       r.UserID,
		Password: r.OldPassword,
	}
}

/*------------------------
//    CREATE USER POST
//----------------------*/

type CreateUserPostRequest struct {
	UserID int    `json:"user_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}

func (r *CreateUserPostRequest) ToUserPostModel() UserPost {
	return UserPost{
		UserID: r.UserID,
		Title:  r.Title,
		Body:   r.Body,
	}
}
