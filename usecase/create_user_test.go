package usecase

import (
	"errors"
	"log"
	"testing"
)

var createUserUseCase *CreateUserUseCase

func setupCreateUserUseCase(userRepo UserRepository, passHasher PasswordHasher) {
	createUserUseCase = NewCreateUserUseCase(userRepo, passHasher)
}

func TestCreateUser(t *testing.T) {
	email := "testmail@email.com"
	err := createUserUseCase.Exec("testname", email, "d81fdw8fd81df81")
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
	createUserUseCase.Exec("testname", email, "d#8wvadw6dv7")
	err := createUserUseCase.Exec("testname", email, "dawbd8awy")
	if errors.Is(err, ErrUserAlreadyExists{Email: email}) == false {
		t.Errorf("Expect err to be ErrUserAlreadyExists Got %v", err)
	}
}
