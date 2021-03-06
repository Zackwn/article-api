package entity

import (
	"net/mail"
	"time"
)

type UserPermission int

const (
	CreateArticlesPermission UserPermission = 1 << iota
	DeleteArticlesPermission

	UnverifiedPermission UserPermission = 0
	VerifiedPermisson    UserPermission = CreateArticlesPermission
)

func NewUser(name, email, picture, password string) (*User, error) {
	// validation
	if len(password) <= 6 {
		return nil, UserError{reason: "Invalid password"}
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, UserError{reason: "Invalid email"}
	}
	if len(name) == 0 {
		return nil, UserError{reason: "Invalid name"}
	}

	user := new(User)
	user.ID = GenerateID()
	user.Verified = false
	user.Permission = UnverifiedPermission
	user.Name = name
	user.Password = password
	user.Email = email
	user.Picture = picture
	user.CreatedAt = Date{time.Now()}
	return user, nil
}

type UserError struct {
	reason string
}

func (err UserError) Error() string {
	return err.reason
}

type User struct {
	ID        string `json:"id" bson:"id"`
	Name      string `json:"name" bson:"name"`
	Password  string `json:"password,omitempty" bson:"password"`
	Email     string `json:"email" bson:"email"`
	Picture   string `json:"picture" bson:"picture"`
	Verified  bool   `json:"verified" bson:"verified"`
	CreatedAt Date   `json:"created_at" bson:"created_at"`

	Permission UserPermission `json:"-" bson:"permission"`
}

func (user *User) AppendPermission(permission UserPermission) {
	user.Permission = user.Permission | permission
}

func (user *User) HasPermission(permission UserPermission) bool {
	return user.Permission&permission != 0
}
