package common

import (
	"time"
)

/*---------------------------------------------------------------------------
// Models are the representation of the database schema. They are used in the Service & Repository Layers.
// They are probably the most important part of the app.
------------------------*/

var AllModels = []interface{}{
	&User{},
	&UserDetail{},
	&UserPost{},
}

type Users []User

type User struct {
	ID        int    `gorm:"primaryKey"`
	Username  string `gorm:"unique;not null"`
	Email     string `gorm:"unique;not null"`
	Password  string `gorm:"not null"`
	IsAdmin   bool
	Details   UserDetail
	Posts     UserPosts `gorm:"foreignKey:UserID;references:ID"`
	Deleted   bool
	CreatedAt time.Time
	UpdatedAt time.Time

	// DTOs
	NewPassword string `gorm:"-"`
}

type UserDetail struct {
	ID        int    `gorm:"primaryKey"`
	UserID    int    `gorm:"unique;not null"`
	FirstName string `gorm:"not null"`
	LastName  string `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserPosts []UserPost

type UserPost struct {
	ID     int    `gorm:"primaryKey"`
	Title  string `gorm:"not null"`
	Body   string `gorm:"type:text"`
	UserID int    `gorm:"not null"`
}

/*---------------------------------------------------------------------------
// Particular Models are a key part of the application, they work as business
// objects and contain some of the logic of the app.
//----------------------*/

/*-------------------
//       AUTH
//-----------------*/

func (u *User) GenerateTokenString(a AuthI) (string, error) {
	return a.GenerateToken(u.ID, u.Username, u.Email, u.GetRole())
}

/*----------------
//     USERS
//--------------*/

func (u *User) GetRole() Role {
	if u.IsAdmin {
		return AdminRole
	}
	return UserRole
}

func (u *User) HashPassword(salt string) {
	u.Password = Hash(u.Password, salt)
}

func (u *User) PasswordMatches(password, salt string) bool {
	return u.Password == Hash(password, salt)
}

func (u *User) OverwriteFields(username, email, password string) {
	if username != "" {
		u.Username = username
	}
	if email != "" {
		u.Email = email
	}
	if password != "" {
		u.Password = password
	}
}

func (u *User) OverwriteDetails(firstName, lastName *string) {
	if firstName != nil {
		u.Details.FirstName = *firstName
	}
	if lastName != nil {
		u.Details.LastName = *lastName
	}
}

/*---------------------------------------------------------------------------
// When the Service layer calls the Repository layer, the output is a Model.
// Here we transform those Models into Response Models, returned on our Custom Responses
// to the Transport layer.
------------------------*/

/*-------------------
//      USERS
//-----------------*/

func (u User) ToResponseModel() ResponseUser {
	return ResponseUser{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		IsAdmin:   u.IsAdmin,
		Details:   u.Details.ToResponseModel(),
		Posts:     u.Posts.ToResponseModel(),
		Deleted:   u.Deleted,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func (u Users) ToResponseModel() []ResponseUser {
	users := []ResponseUser{}
	for _, user := range u {
		users = append(users, user.ToResponseModel())
	}
	return users
}

func (u UserDetail) ToResponseModel() ResponseUserDetail {
	return ResponseUserDetail{
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}

/*-------------------
//      POSTS
//-----------------*/

func (p UserPost) ToResponseModel() ResponseUserPost {
	return ResponseUserPost{
		ID:    p.ID,
		Title: p.Title,
		Body:  p.Body,
	}
}

func (p UserPosts) ToResponseModel() []ResponseUserPost {
	posts := []ResponseUserPost{}
	for _, post := range p {
		posts = append(posts, post.ToResponseModel())
	}
	return posts
}
