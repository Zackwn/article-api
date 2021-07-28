package usecase

import "fmt"

func NewForgotPasswordUseCase(userRepo UserRepository, fph ForgotPasswordHandler, emailService EmailService) *ForgotPasswordUseCase {
	return &ForgotPasswordUseCase{userRepository: userRepo, fph: fph, emailService: emailService}
}

type ForgotPasswordUseCase struct {
	userRepository UserRepository
	fph            ForgotPasswordHandler
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
	fphToken, err := forgotPasswordUseCase.fph.Request(user)
	if err != nil {
		return ErrInternalServer{}
	}

	html := fmt.Sprintf(`<p>Your redefine password token: %v</p>`, fphToken)

	err = forgotPasswordUseCase.emailService.Send(user.Email, "Article-API: Redefine password", html)
	if err != nil {
		return ErrInternalServer{err}
	}

	return nil
}
