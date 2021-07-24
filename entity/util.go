package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func GenerateID() string {
	return uuid.NewString()
}

type Date struct {
	time.Time
}

func (date Date) MarshalJSON() ([]byte, error) {
	stamp := fmt.Sprintf("\"%s\"", date.Format(time.RFC3339))
	return []byte(stamp), nil
}
