package domain

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type Token struct {
	ID        int64
	User      *User
	TokenHash string
	ExpiresAt time.Time
}

func NewToken(user *User, ttl time.Duration) (*Token, string, error) {
	rawBytes := make([]byte, 32)
	if _, err := rand.Read(rawBytes); err != nil {
		return nil, "", err
	}

	rawToken := hex.EncodeToString(rawBytes)
	hash := HashToken(rawToken)

	return &Token{
		User:      user,
		TokenHash: hash,
		ExpiresAt: time.Now().UTC().Add(ttl),
	}, rawToken, nil
}

func HashToken(rawToken string) string {
	sum := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(sum[:])
}

func (t *Token) IsExpired() bool {
	return time.Now().UTC().After(t.ExpiresAt)
}
