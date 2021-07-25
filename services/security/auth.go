package security

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

func NewAuthProvider() *AuthProvider {
	authProvider := new(AuthProvider)
	err := godotenv.Load()
	if err != nil {
		log.Panic(err)
	} else if os.Getenv("JWT_SECRET") == "" {
		log.Panic(errors.New("missing env `JWT_SECRET`"))
	}
	authProvider.secret = []byte(os.Getenv("JWT_SECRET"))
	return authProvider
}

type AuthProvider struct {
	secret []byte
}

type TokenPayload struct {
	UserID string
}

type TokenClaims struct {
	UserID string
	jwt.StandardClaims
}

func (auth AuthProvider) Sign(userID string) string {
	tokenClaims := TokenClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	tokenString, err := token.SignedString(auth.secret)
	if err != nil {
		log.Panic(err)
	}
	return tokenString
}

func (auth AuthProvider) Verify(tokenString string) (*TokenPayload, bool) {
	var tokenClaims TokenClaims
	_, err := jwt.ParseWithClaims(
		tokenString,
		&tokenClaims,
		func(token *jwt.Token) (interface{}, error) {
			return auth.secret, nil
		},
	)
	if err != nil {
		return nil, false
	} else if tokenClaims.UserID == "" {
		return nil, false
	}
	payload := TokenPayload{
		UserID: tokenClaims.UserID,
	}
	return &payload, true
}
