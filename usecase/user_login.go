package usecase

type ErrUserDoNotExists struct {
	Email string
}

func (err ErrUserDoNotExists) Error() string {
	return `user with email: "` + err.Email + `" do not exists.`
}

func NewUserLoginUseCase(repo UserRepository, passHasher PasswordHasher, auth AuthProvider) *UserLoginUseCase {
	return &UserLoginUseCase{userRepository: repo, passwordHasher: passHasher, authProvider: auth}
}

type UserLoginUseCase struct {
	userRepository UserRepository
	passwordHasher PasswordHasher
	authProvider   AuthProvider
}

func (userLoginUseCase UserLoginUseCase) Exec(email, password string) (string, error) {
	user, found := userLoginUseCase.userRepository.FindByEmail(email)
	if !found {
		return "", &ErrUserDoNotExists{Email: email}
	}
	err := userLoginUseCase.passwordHasher.CompareHashAndPassword(user.Password, password)
	if err != nil {
		return "", err
	}
	accessToken := userLoginUseCase.authProvider.Sign(user.ID)
	return accessToken, nil
}
