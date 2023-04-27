package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// All types of errors returned by the VerifyToken function
var (
	ErrExpiredToken = errors.New("token: token has expired")
	ErrInvalidToken = errors.New("token: token is invalid")
)

const At = 12

// TokenPayload contains the payload data of the token
type TokenPayload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

type Footer struct {
	ExpiredAt time.Time `json:"exp"`
}

// NewPayload creates a new token payload with a specific username and duration
func NewPayload(username string, duration time.Duration) (*TokenPayload, *Footer, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, nil, err
	}

	payload := &TokenPayload{
		ID:        tokenID,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	footer := &Footer{
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, footer, nil
}

// Validate checks whether the token is valid or not
func (payload *TokenPayload) Validate() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
