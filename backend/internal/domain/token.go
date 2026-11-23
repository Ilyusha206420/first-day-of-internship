package domain

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type Token struct {
	ID        int64
	UserID    int64
	TokenHash string
	ExpiresAt time.Time
}

func NewToken(userID int64, ttl time.Duration) (*Token, string, error) {
	rawBytes := make([]byte, 32)
	if _, err := rand.Read(rawBytes); err != nil {
		return nil, "", err
	}

	rawToken := hex.EncodeToString(rawBytes)
	hash := HashToken(rawToken)

	return &Token{
		UserID:    userID,
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
