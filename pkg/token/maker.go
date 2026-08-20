package token

import (
	"time"
)

type Maker interface {
	CreateToken(userId int64, duration time.Duration, tokenType TokenType) (string, *Payload, error)
	VerifyToken(tokenStr string, tokenType TokenType) (*Payload, error)
}

