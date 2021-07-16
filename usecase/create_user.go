package usecase

import (
	"github.com/zackwn/article-api/entity"
)

func NewCreateUserUseCase(repo UserRepository, passHasher PasswordHasher) *CreateUserUseCase {
	return &CreateUserUseCase{userRepository: repo, passwordHasher: passHasher}
}

type CreateUserUseCase struct {
	userRepository UserRepository
	passwordHasher PasswordHasher
}

func (createUserUseCase CreateUserUseCase) Exec(name, email, password string) UseCaseErr {
	_, userAlreadyExists := createUserUseCase.userRepository.FindByEmail(email)
	if userAlreadyExists {
		return ErrUserAlreadyExists{Email: email}
	}

	hashPassword := createUserUseCase.passwordHasher.HashPassword(password)

	user, err := entity.NewUser(name, email, hashPassword)
	if err != nil {
		return ErrInvalidEntityField{errMessage: err.Error()}
	}
	err = createUserUseCase.userRepository.Store(user)
	if err != nil {
		return ErrInvalidEntityField{errMessage: err.Error()}
	}
	return nil
}
