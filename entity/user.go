package entity

import (
	"net/mail"
	"time"
)

type UserPermission int

const (
	CreateArticlesPermission UserPermission = 1 << iota
	DeleteArticlesPermission

	DefaultPermission UserPermission = CreateArticlesPermission
)

func NewUser(name, email, password string) (*User, error) {
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
	user.Name = name
	user.Password = password
	user.Email = email
	user.Permission = DefaultPermission
	user.CreatedAt = time.Now()
	return user, nil
}

type UserError struct {
	reason string
}

func (err UserError) Error() string {
	return err.reason
}

type User struct {
	ID        string    `json:"id" bson:"id"`
	Name      string    `json:"name" bson:"name"`
	Password  string    `json:"password" bson:"password"`
	Email     string    `json:"email" bson:"email"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`

	Permission UserPermission `json:"-" bson:"permission"`
}

func (user *User) AppendPermission(permission UserPermission) {
	user.Permission = user.Permission | permission
}

func (user *User) HasPermission(permission UserPermission) bool {
	return user.Permission&permission != 0
}
