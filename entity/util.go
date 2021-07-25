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
	// d := date.Format(time.RFC3339)
	d := date.Format("01/02/2006 15:04")
	stamp := fmt.Sprintf("\"%s\"", d)
	return []byte(stamp), nil
}
