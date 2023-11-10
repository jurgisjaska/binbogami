package service

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/labstack/gommon/random"
)

func CreateSalt(email *string) string {
	s := fmt.Sprintf(
		"%s%s%s",
		email,
		time.Now().Format(time.RFC3339Nano),
		random.String(16),
	)

	hash := sha256.New()
	hash.Write([]byte(s))

	return fmt.Sprintf("%x", hash.Sum(nil))
}
