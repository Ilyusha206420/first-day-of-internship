package domain

import (
	"context"
	"time"
)

type Token struct {
	ID        int64
	UserID    int64
	TokenHash string
	ExpiresAt time.Time
}

type TokenRepository interface {
	Add(ctx context.Context, token *Token) error
	FindByHash(ctx context.Context, tokenHash string) (*Token, error)
	DeleteByHash(ctx context.Context, tokenHash string) error
	DeleteAllByUserID(ctx context.Context, userID int64) error
}
