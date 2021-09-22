package usecase

import (
	"fmt"

	"github.com/zackwn/article-api/entity"
)

func NewVerifyAccountUseCase(userRepo UserRepository, tempToken TempToken) *VerifyAccountUseCase {
	return &VerifyAccountUseCase{tempToken: tempToken, userRepository: userRepo}
}

type VerifyAccountUseCase struct {
	tempToken      TempToken
	userRepository UserRepository
}

type VerifyAccountDTO struct {
	Token string
}

func (verifyAccountUseCase VerifyAccountUseCase) Exec(dto *VerifyAccountDTO) UseCaseErr {
	userID, valid := verifyAccountUseCase.tempToken.Validate(dto.Token)
	if !valid {
		return ErrInvalidVerifyAccountRequest{}
	}
	user, found := verifyAccountUseCase.userRepository.FindByID(userID)
	if !found {
		return ErrUserDoNotExists{}
	}
	// verify user account
	user.Verified = true
	user.AppendPermission(entity.VerifiedPermisson)
	fmt.Println("verify user:", user)
	err := verifyAccountUseCase.userRepository.Update(user)
	if err != nil {
		return ErrInternalServer{err}
	}
	return nil
}
