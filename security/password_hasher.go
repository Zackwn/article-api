package security

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func NewPasswordHasher() *PasswordHasher {
	return new(PasswordHasher)
}

type ErrWrongPassword struct{}

func (ErrWrongPassword) Error() string {
	return "wrong password."
}

type PasswordHasher struct{}

func (PasswordHasher) HashPassword(password string) string {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	return string(hashPassword)
}

func (PasswordHasher) CompareHashAndPassword(hashPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	if err != nil {
		return ErrWrongPassword{}
	}
	return nil
}
