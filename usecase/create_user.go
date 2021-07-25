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

type CreateUserDTO struct {
	Name     string
	Email    string
	Password string
	Picture  string
}

func (createUserUseCase CreateUserUseCase) Exec(dto *CreateUserDTO) UseCaseErr {
	_, userAlreadyExists := createUserUseCase.userRepository.FindByEmail(dto.Email)
	if userAlreadyExists {
		return ErrUserAlreadyExists{Email: dto.Email}
	}

	hashPassword := createUserUseCase.passwordHasher.HashPassword(dto.Password)

	user, err := entity.NewUser(dto.Name, dto.Email, dto.Picture, hashPassword)
	if err != nil {
		return ErrInvalidEntityField{errMessage: err.Error()}
	}
	err = createUserUseCase.userRepository.Store(user)
	if err != nil {
		return ErrInvalidEntityField{errMessage: err.Error()}
	}
	return nil
}
