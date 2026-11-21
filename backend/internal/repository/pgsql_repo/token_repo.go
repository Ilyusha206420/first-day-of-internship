package pgsql_repo

import (
	"context"
	"time"

	"github.com/Ilyusha206420/first-day-of-internship/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type PGTokenRepo struct {
	db *sqlx.DB
}

type tokenRow struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	TokenHash string    `db:"refresh_token"`
	ExpiresAt time.Time `db:"expires_at"`
}

func (tr *PGTokenRepo) Add(ctx context.Context, token *domain.Token) error {
	query := `
		INSERT INTO tokens (user_id, refresh_token, expires_at)
		VALUES ($1, $2, $3)`
	_, err := tr.db.ExecContext(ctx, query, token.UserID, token.TokenHash, token.ExpiresAt)
	return err
}

func (tr *PGTokenRepo) FindByHash(ctx context.Context, tokenHash string) (*domain.Token, error) {
	return nil, nil
}

func (tr *PGTokenRepo) DeleteByHash(ctx context.Context, tokenHash string) error {
	return nil
}

func (tr *PGTokenRepo) DeleteAllByUserID(ctx context.Context, userID int64) error {
	return nil
}
