package usecase

import (
	"fmt"

	"github.com/zackwn/article-api/entity"
)

func NewSendVerifyAccoutTokenUseCase(userRepo UserRepository, auth AuthProvider, tp TempToken, emailService EmailService) *SendVerifyAccountTokenUseCase {
	return &SendVerifyAccountTokenUseCase{userRepository: userRepo, authProvider: auth, tempToken: tp, emailService: emailService}
}

type SendVerifyAccountTokenUseCase struct {
	userRepository UserRepository
	authProvider   AuthProvider
	tempToken      TempToken
	emailService   EmailService
}

type SendVerifyAccountTokenDTO struct {
	AccessToken string
	UserID      string
}

func (sendVerifyAccountTokenUseCase SendVerifyAccountTokenUseCase) Exec(dto *SendVerifyAccountTokenDTO) UseCaseErr {
	var user *entity.User
	if dto.UserID == "" {
		payload, valid := sendVerifyAccountTokenUseCase.authProvider.Verify(dto.AccessToken)
		if !valid {
			return ErrInvalidAccessToken{}
		}
		dto.UserID = payload.UserID
	}
	user, found := sendVerifyAccountTokenUseCase.userRepository.FindByID(dto.UserID)
	if !found {
		return ErrUserDoNotExists{}
	}
	token, err := sendVerifyAccountTokenUseCase.tempToken.New(user)
	if err != nil {
		return ErrInternalServer{err}
	}
	html := fmt.Sprintf("<p>Your account verifier token: %v</p>", token)
	err = sendVerifyAccountTokenUseCase.emailService.Send(user.Email, "Verify Account", html)
	if err != nil {
		return ErrInternalServer{err}
	}
	return nil
}
