package usecase

import (
	"github.com/zackwn/article-api/entity"
	"github.com/zackwn/article-api/security"
)

type UserRepository interface {
	FindByID(ID string) (*entity.User, bool)
	FindByEmail(email string) (*entity.User, bool)

	Store(user *entity.User) error
}

type PasswordHasher interface {
	HashPassword(password string) string
	CompareHashAndPassword(hashPassword, password string) error
}

type AuthProvider interface {
	Sign(userID string) string
	Verify(tokenString string) (*security.TokenPayload, bool)
}
