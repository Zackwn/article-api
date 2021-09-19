package forgotpasswordhandler

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/zackwn/article-api/entity"
)

func NewForgotPasswordHandler(redis *redis.Client) *ForgotPasswordHandler {
	return &ForgotPasswordHandler{redis: redis, prefix: "FPH:"}
}

type ForgotPasswordHandler struct {
	prefix string
	redis  *redis.Client
}

func (fph ForgotPasswordHandler) Request(user *entity.User) (string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	token := uuid.NewString()
	key := fph.prefix + token
	value := user.ID
	expiration := 10 * time.Minute
	err := fph.redis.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return "", err
	}
	return token, nil
}

func (fph ForgotPasswordHandler) Validate(token string) (string, bool) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	key := fph.prefix + token
	result := fph.redis.Get(ctx, key)
	err := result.Err()
	if err != nil {
		return "", false
	}
	return result.Val(), true
}
