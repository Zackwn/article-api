package usecase

func NewChangePasswordUseCase(userRepo UserRepository, fph ForgotPasswordHandler, passHasher PasswordHasher) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{userRepository: userRepo, fph: fph, passwordHasher: passHasher}
}

type ChangePasswordUseCase struct {
	userRepository UserRepository
	fph            ForgotPasswordHandler
	passwordHasher PasswordHasher
}

type ChangePasswordDTO struct {
	Token       string
	NewPassword string
}

func (changePasswordUseCase ChangePasswordUseCase) Exec(dto *ChangePasswordDTO) UseCaseErr {
	userID, valid := changePasswordUseCase.fph.Validate(dto.Token)
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
