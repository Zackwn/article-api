package usecase

import (
	"errors"
	"log"
	"testing"
)

var createUserUseCase *CreateUserUseCase

func setupCreateUserUseCase(userRepo UserRepository, passHasher PasswordHasher, emailService EmailService, tk TempToken) {
	createUserUseCase = NewCreateUserUseCase(userRepo, passHasher, emailService, tk)
}

func TestCreateUser(t *testing.T) {
	email := "testmail@email.com"
	dto := &CreateUserDTO{
		Name:     "testname",
		Email:    email,
		Password: "d81fdw8fd81df81",
		Picture:  "picture",
	}
	_, err := createUserUseCase.Exec(dto)
	if err != nil {
		log.Fatal(err)
	}
	_, found := createUserUseCase.userRepository.FindByEmail(email)
	if found != true {
		t.Errorf("Expect %v Got %v", "true", found)
	}
}

func TestCreateUserErr(t *testing.T) {
	email := "testexecerr@email.com"
	dto := &CreateUserDTO{
		Name:     "testname",
		Email:    email,
		Password: "d#8wvadw6dv7",
		Picture:  "picture",
	}
	createUserUseCase.Exec(dto)
	_, err := createUserUseCase.Exec(dto)
	if errors.Is(err, ErrUserAlreadyExists{Email: email}) == false {
		t.Errorf("Expect err to be ErrUserAlreadyExists Got %v", err)
	}
}
