package entity

import "net/mail"

type UserPermission int

const (
	CreateArticles UserPermission = 1 << iota
	DeleteArticles

	DefaultPermission UserPermission = CreateArticles
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
	return user, nil
}

type UserError struct {
	reason string
}

func (err UserError) Error() string {
	return err.reason
}

type User struct {
	ID       string `json:"id" bson:"id"`
	Name     string `json:"name" bson:"name"`
	Password string `json:"password" bson:"password"`
	Email    string `json:"email" bson:"email"`

	Permission UserPermission `json:"-" bson:"permission"`
}

func (user *User) AppendPermission(permission UserPermission) {
	user.Permission = user.Permission | permission
}

func (user *User) HasPermission(permission UserPermission) bool {
	return user.Permission&permission != 0
}
