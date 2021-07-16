package security

import (
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

func TestCompareHashAndPasswordValid(t *testing.T) {
	password := "d17gd8f16df1wm"
	hash := passwordHasher.HashPassword(password)

	passValid := passwordHasher.CompareHashAndPassword(hash, password)
	if !passValid {
		t.Errorf("Expect %v Got %v", "true", passValid)
	}
}

func TestCompareHashAndPasswordInvalid(t *testing.T) {
	password := "d1hu912dg7861e"
	hash := passwordHasher.HashPassword(password)

	passValid := passwordHasher.CompareHashAndPassword(hash, "wrong password")

	if passValid {
		t.Errorf("Expect %v Got %v", "false", passValid)
	}
}
