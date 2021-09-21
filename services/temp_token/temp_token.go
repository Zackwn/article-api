package temptoken

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/zackwn/article-api/entity"
)

func NewTempToken(redis *redis.Client) *TempToken {
	return &TempToken{redis: redis, prefix: "TK:"}
}

type TempToken struct {
	prefix string
	redis  *redis.Client
}

func (tk TempToken) New(user *entity.User) (string, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	token := uuid.NewString()
	key := tk.prefix + token
	value := user.ID
	expiration := 10 * time.Minute
	err := tk.redis.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return "", err
	}
	return token, nil
}

func (tk TempToken) Validate(token string) (string, bool) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	key := tk.prefix + token
	result := tk.redis.Get(ctx, key)
	err := result.Err()
	if err != nil {
		return "", false
	}
	return result.Val(), true
}
