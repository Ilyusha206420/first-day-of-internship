package pgsql_repo

import (
	"context"
	"database/sql"
	"errors"
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
	var trow tokenRow
	query := `SELECT id, user_id, refresh_token, expires_at FROM tokens WHERE refresh_token = $1`

	err := tr.db.GetContext(ctx, &trow, query, tokenHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &domain.Token{
		ID:        trow.ID,
		UserID:    trow.UserID,
		TokenHash: trow.TokenHash,
		ExpiresAt: trow.ExpiresAt,
	}, nil
}

func (tr *PGTokenRepo) DeleteByHash(ctx context.Context, tokenHash string) error {
	query := `DELETE FROM tokens WHERE refresh_token = $1`
	_, err := tr.db.ExecContext(ctx, query, tokenHash)
	return err
}

func (tr *PGTokenRepo) DeleteAllByUserID(ctx context.Context, userID int64) error {
	query := `DELETE FROM TOKENS WHERE user_id = $1`
	_, err := tr.db.ExecContext(ctx, query, userID)
	return err
}
