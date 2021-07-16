package usecase

import "net/http"

type ErrUserAlreadyExists struct {
	Email string
}

func (err ErrUserAlreadyExists) Error() string {
	return `user with email: "` + err.Email + `" already exists`
}

func (ErrUserAlreadyExists) HttpStatus() int {
	return http.StatusBadRequest
}

type ErrInvalidAccessToken struct{}

func (ErrInvalidAccessToken) Error() string {
	return "invalid access token"
}

func (ErrInvalidAccessToken) HttpStatus() int {
	return http.StatusUnauthorized
}

type ErrUserDoNotExists struct{}

func (err ErrUserDoNotExists) Error() string {
	return `user do not exists`
}

func (ErrUserDoNotExists) HttpStatus() int {
	return http.StatusBadRequest
}

type ErrWrongPassword struct{}

func (ErrWrongPassword) Error() string {
	return "wrong password"
}

func (ErrWrongPassword) HttpStatus() int {
	return http.StatusUnauthorized
}

type ErrInvalidEntityField struct {
	errMessage string
}

func (err ErrInvalidEntityField) Error() string {
	return err.errMessage
}

func (ErrInvalidEntityField) HttpStatus() int {
	return http.StatusBadRequest
}

type ErrForbiddenUserAction struct{}

func (ErrForbiddenUserAction) Error() string {
	return "Forbidden"
}

func (ErrForbiddenUserAction) HttpStatus() int {
	return http.StatusForbidden
}
