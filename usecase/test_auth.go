package usecase

import "fmt"

func NewTestAuthUseCase(auth AuthProvider) *TestAuthUseCase {
	return &TestAuthUseCase{authProvider: auth}
}

type TestAuthUseCase struct {
	authProvider AuthProvider
}

func (testAuthUseCase TestAuthUseCase) Exec(accessToken string) error {
	fmt.Println("access_token", accessToken)
	payload, valid := testAuthUseCase.authProvider.Verify(accessToken)
	if valid {
		fmt.Println("Valid payload:", payload)
	}
	return nil
}
