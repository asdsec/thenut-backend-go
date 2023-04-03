package token

import "time"

// TokenMaker is an interface managing tokens
type TokenMaker interface {
	// CreateToken creates a new token for a specific username and duration
	CreateToken(username string, duration time.Duration) (string, *TokenPayload, error)

	// VerifyToken checks whether the token is valid or not
	VerifyToken(token string) (*TokenPayload, error)
}
