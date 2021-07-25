package entity

import (
	"errors"
	"testing"
)

func TestNewUserValid(t *testing.T) {
	_, err := NewUser("validname", "valid@email.com", "picture", "v4lidp#ssword")

	if err != nil {
		t.Errorf("Expect %v Got %v", nil, err)
	}
}

func TestNewUserInvalid(t *testing.T) {
	_, err := NewUser("validname", "invalidEmail", "picture", "v4lidp#ssword")

	expectErr := UserError{reason: "Invalid email"}
	if errors.Is(err, expectErr) == false {
		t.Errorf("Expect %v Got %v", expectErr, err)
	}
}
