package usecase

import (
	"github.com/zackwn/article-api/entity"
)

func NewCreateUserUseCase(repo UserRepository, passHasher PasswordHasher, emailService EmailService, tk TempToken) *CreateUserUseCase {
	return &CreateUserUseCase{userRepository: repo, passwordHasher: passHasher, emailService: emailService, tempToken: tk}
}

type CreateUserUseCase struct {
	userRepository UserRepository
	passwordHasher PasswordHasher
	emailService   EmailService
	tempToken      TempToken
}

type CreateUserDTO struct {
	Name     string
	Email    string
	Password string
	Picture  string
}

func (createUserUseCase CreateUserUseCase) Exec(dto *CreateUserDTO) (*entity.User, UseCaseErr) {
	_, userAlreadyExists := createUserUseCase.userRepository.FindByEmail(dto.Email)
	if userAlreadyExists {
		return nil, ErrUserAlreadyExists{Email: dto.Email}
	}

	hashPassword := createUserUseCase.passwordHasher.HashPassword(dto.Password)

	user, err := entity.NewUser(dto.Name, dto.Email, dto.Picture, hashPassword)
	if err != nil {
		return nil, ErrInvalidEntityField{errMessage: err.Error()}
	}
	err = createUserUseCase.userRepository.Store(user)
	if err != nil {
		return nil, ErrInvalidEntityField{errMessage: err.Error()}
	}
	return user, nil
}
