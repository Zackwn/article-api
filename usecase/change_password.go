package usecase

func NewChangePasswordUseCase(userRepo UserRepository, tk TempToken, passHasher PasswordHasher) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{userRepository: userRepo, tempToken: tk, passwordHasher: passHasher}
}

type ChangePasswordUseCase struct {
	userRepository UserRepository
	tempToken      TempToken
	passwordHasher PasswordHasher
}

type ChangePasswordDTO struct {
	Token       string
	NewPassword string
}

func (changePasswordUseCase ChangePasswordUseCase) Exec(dto *ChangePasswordDTO) UseCaseErr {
	userID, valid := changePasswordUseCase.tempToken.Validate(dto.Token)
	if !valid {
		return ErrInvalidChangePasswordRequest{}
	}
	user, found := changePasswordUseCase.userRepository.FindByID(userID)
	if !found {
		return ErrInvalidChangePasswordRequest{}
	}
	newPasswordHash := changePasswordUseCase.passwordHasher.HashPassword(dto.NewPassword)
	user.Password = newPasswordHash
	err := changePasswordUseCase.userRepository.Update(user)
	if err != nil {
		return ErrInternalServer{err}
	}
	return nil
}
