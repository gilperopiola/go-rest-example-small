package common

import (
	"time"
)

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

func (r *LoginRequest) ToUserModel() User {
	user := User{Password: r.Password}

	if validEmailRegex.MatchString(r.UsernameOrEmail) {
		user.Email = r.UsernameOrEmail
	} else {
		user.Username = r.UsernameOrEmail
	}

	return user
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

func (r *GetUserRequest) ToUserModel() User {
	return User{ID: r.UserID}
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
func (r *DeleteUserRequest) ToUserModel() User {
	return User{ID: r.UserID}
}

func (r *SearchUsersRequest) ToUserModel() User {
	return User{Username: r.Username}
}

func (r *ChangePasswordRequest) ToUserModel() User {
	return User{
		ID:       r.UserID,
		Password: r.OldPassword,
	}
}

func (r *CreateUserPostRequest) ToUserPostModel() UserPost {
	return UserPost{
		UserID: r.UserID,
		Title:  r.Title,
		Body:   r.Body,
	}
}
