package usecase

func NewUserLoginUseCase(repo UserRepository, passHasher PasswordHasher, auth AuthProvider) *UserLoginUseCase {
	return &UserLoginUseCase{userRepository: repo, passwordHasher: passHasher, authProvider: auth}
}

type UserLoginUseCase struct {
	userRepository UserRepository
	passwordHasher PasswordHasher
	authProvider   AuthProvider
}

func (userLoginUseCase UserLoginUseCase) Exec(email, password string) (string, UseCaseErr) {
	user, found := userLoginUseCase.userRepository.FindByEmail(email)
	if !found {
		return "", &ErrUserDoNotExists{}
	}
	passValid := userLoginUseCase.passwordHasher.CompareHashAndPassword(user.Password, password)
	if !passValid {
		return "", ErrWrongPassword{}
	}
	accessToken := userLoginUseCase.authProvider.Sign(user.ID)
	return accessToken, nil
}
