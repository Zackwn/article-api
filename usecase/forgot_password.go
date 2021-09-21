package usecase

import "fmt"

func NewForgotPasswordUseCase(userRepo UserRepository, tk TempToken, emailService EmailService) *ForgotPasswordUseCase {
	return &ForgotPasswordUseCase{userRepository: userRepo, tempToken: tk, emailService: emailService}
}

type ForgotPasswordUseCase struct {
	userRepository UserRepository
	tempToken      TempToken
	emailService   EmailService
}

type ForgotPasswordDTO struct {
	Email string
}

func (forgotPasswordUseCase ForgotPasswordUseCase) Exec(dto ForgotPasswordDTO) UseCaseErr {
	user, found := forgotPasswordUseCase.userRepository.FindByEmail(dto.Email)
	if !found {
		return ErrUserDoNotExists{}
	}
	token, err := forgotPasswordUseCase.tempToken.New(user)
	if err != nil {
		return ErrInternalServer{}
	}

	html := fmt.Sprintf(`<p>Your redefine password token: %v</p>`, token)

	err = forgotPasswordUseCase.emailService.Send(user.Email, "Article-API: Redefine password", html)
	if err != nil {
		return ErrInternalServer{err}
	}

	return nil
}
