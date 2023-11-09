package common

import "time"

// These aren't the HTTP Responses that the API will return, but the responses that the Service Layer
// returns to the Transport Layer.

type AllResponses interface {
	SignupResponse |
		LoginResponse |
		CreateUserResponse |
		GetUserResponse |
		UpdateUserResponse |
		DeleteUserResponse |
		SearchUsersResponse |
		ChangePasswordResponse |
		CreateUserPostResponse
}

/*-------------------
//      AUTH
//-----------------*/

type SignupResponse struct {
	User ResponseUser `json:"user"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

/*--------------------
//      USERS
//------------------*/

type CreateUserResponse struct {
	User ResponseUser `json:"user"`
}

type GetUserResponse struct {
	User ResponseUser `json:"user"`
}

type UpdateUserResponse struct {
	User ResponseUser `json:"user"`
}

type DeleteUserResponse struct {
	User ResponseUser `json:"user"`
}

type SearchUsersResponse struct {
	Users   []ResponseUser `json:"users"`
	Page    int            `json:"page"`
	PerPage int            `json:"per_page"`
}

type ChangePasswordResponse struct {
	User ResponseUser `json:"user"`
}

type CreateUserPostResponse struct {
	UserPost ResponseUserPost `json:"user_post"`
}

/*-----------------------
//    RESPONSE MODELS
//---------------------*/

type ResponseUser struct {
	ID        int                `json:"id"`
	Username  string             `json:"username"`
	Email     string             `json:"email"`
	IsAdmin   bool               `json:"is_admin,omitempty"`
	Details   ResponseUserDetail `json:"details"`
	Posts     []ResponseUserPost `json:"posts"`
	Deleted   bool               `json:"deleted,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty"`
}

type ResponseUserDetail struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type ResponseUserPost struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Body  string `json:"body"`
}
