package security

import (
	"errors"
	"os"
	"testing"
)

var passwordHasher *PasswordHasher

func setup() {
	passwordHasher = NewPasswordHasher()
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func TestHashPassword(t *testing.T) {
	password := "1678f6d1ydvy1"
	hash := passwordHasher.HashPassword(password)

	hashExpectLength := 60
	if len(hash) != hashExpectLength {
		t.Errorf("Expect %v Got %v", hashExpectLength, len(hash))
	}
}

func TestCompareHashAndPassword(t *testing.T) {
	password := "d17gd8f16df1wm"
	hash := passwordHasher.HashPassword(password)

	err := passwordHasher.CompareHashAndPassword(hash, password)
	if err != nil {
		t.Errorf("Expect %v Got %v", nil, err)
	}
}

func TestCompareHashAndPasswordErr(t *testing.T) {
	password := "d1hu912dg7861e"
	hash := passwordHasher.HashPassword(password)

	err := passwordHasher.CompareHashAndPassword(hash, "wrong password")

	expectErr := ErrWrongPassword{}
	if errors.Is(err, expectErr) == false {
		t.Errorf("Expect ErrWrongPassword Got %v", err)
	}
}
