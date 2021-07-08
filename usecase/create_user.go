package usecase

import (
	"github.com/zackwn/article-api/entity"
)

type ErrUserAlreadyExists struct {
	Email string
}

func (err ErrUserAlreadyExists) Error() string {
	return `user with email: "` + err.Email + `" already exists.`
}

func NewCreateUserUseCase(repo UserRepository, passHasher PasswordHasher) *CreateUserUseCase {
	return &CreateUserUseCase{userRepository: repo, passwordHasher: passHasher}
}

type CreateUserUseCase struct {
	userRepository UserRepository
	passwordHasher PasswordHasher
}

func (createUserUseCase CreateUserUseCase) Exec(name, email, password string) error {
	_, userAlreadyExists := createUserUseCase.userRepository.FindByEmail(email)
	if userAlreadyExists {
		return ErrUserAlreadyExists{Email: email}
	}

	hashPassword := createUserUseCase.passwordHasher.HashPassword(password)

	user, err := entity.NewUser(name, email, hashPassword)
	if err != nil {
		return err
	}

	return createUserUseCase.userRepository.Store(user)
}
