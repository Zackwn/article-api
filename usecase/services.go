package usecase

import (
	"github.com/zackwn/article-api/entity"
	"github.com/zackwn/article-api/services/security"
)

type UserRepository interface {
	FindByID(ID string) (*entity.User, bool)
	FindByEmail(email string) (*entity.User, bool)

	Store(user *entity.User) error
}

type ArticleRepository interface {
	All() []*entity.Article
	FindByID(ID string) (*entity.Article, bool)
	FindAllByAuthor(authorID string) []*entity.Article

	Store(article *entity.Article) error
}

type PasswordHasher interface {
	HashPassword(password string) string
	CompareHashAndPassword(hashPassword, password string) bool
}

type AuthProvider interface {
	Sign(userID string) string
	Verify(tokenString string) (*security.TokenPayload, bool)
}

type ForgotPasswordHandler interface {
	Request(user *entity.User) (token string, err error)
	Validate(token string) (userid string, valid bool)
}

type EmailService interface {
	Send(to string, subject string, html string) error
}
