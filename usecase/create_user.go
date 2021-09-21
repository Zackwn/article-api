package usecase

import (
	"fmt"

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

	token, err := createUserUseCase.tempToken.New(user)
	if err != nil {
		return ErrInternalServer{err}
	}
	html := fmt.Sprintf("<p>Your account verifier token: %v</p>", token)
	createUserUseCase.emailService.Send(user.Email, "Verify Account", html)

	return nil
}
