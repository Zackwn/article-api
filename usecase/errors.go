package usecase

type ErrUserAlreadyExists struct {
	Email string
}

func (err ErrUserAlreadyExists) Error() string {
	return `user with email: "` + err.Email + `" already exists.`
}

type ErrInvalidAccessToken struct{}

func (ErrInvalidAccessToken) Error() string {
	return "invalid access token"
}
