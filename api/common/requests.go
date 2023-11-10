package common

type AllRequests interface {
	SignupRequest |
		LoginRequest |
		CreateUserRequest |
		GetUserRequest |
		UpdateUserRequest |
		DeleteUserRequest |
		SearchUsersRequest |
		ChangePasswordRequest |
		CreateUserPostRequest
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

/*--------------
//    LOGIN
//------------*/

type LoginRequest struct {
	UsernameOrEmail string `json:"username_or_email"`
	Password        string `json:"password"`
}

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

/*--------------------
//     GET USER
//------------------*/

type GetUserRequest struct {
	UserID int `json:"user_id"`
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

/*--------------------
//    DELETE USER
//------------------*/

type DeleteUserRequest struct {
	UserID int `json:"user_id"`
}

/*--------------------
//    SEARCH USERS
//------------------*/

type SearchUsersRequest struct {
	Username string `json:"username"`
	Page     int    `json:"page"`
	PerPage  int    `json:"per_page"`
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

/*------------------------
//    CREATE USER POST
//----------------------*/

type CreateUserPostRequest struct {
	UserID int    `json:"user_id"`
	Title  string `json:"title"`
	Body   string `json:"body"`
}
