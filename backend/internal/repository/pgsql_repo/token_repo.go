package pgsql_repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Ilyusha206420/first-day-of-internship/backend/internal/domain"
	"github.com/Ilyusha206420/first-day-of-internship/backend/internal/repository/rows"
	"github.com/jmoiron/sqlx"
)

type userLookup interface {
	FindByID(ctx context.Context, userID int64) (*domain.User, error)
}

type TokenRepo struct {
	db       *sqlx.DB
	userRepo userLookup
}

func NewTokenRepo(db *sqlx.DB, userRepo userLookup) *TokenRepo {
	return &TokenRepo{db: db, userRepo: userRepo}
}

func (tr *TokenRepo) Add(ctx context.Context, token *domain.Token) error {
	query := `
		INSERT INTO tokens (user_id, refresh_token, expires_at)
		VALUES ($1, $2, $3)`
	_, err := tr.db.ExecContext(ctx, query, token.User.ID, token.TokenHash, token.ExpiresAt)
	return err
}

func (tr *TokenRepo) FindByHash(ctx context.Context, tokenHash string) (*domain.Token, error) {
	var trow rows.TokenRow
	query := `SELECT id, user_id, refresh_token, expires_at FROM tokens WHERE refresh_token = $1`

	err := tr.db.GetContext(ctx, &trow, query, tokenHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	user, err := tr.userRepo.FindByID(ctx, trow.UserID)
	if err != nil {
		return nil, err
	}

	return &domain.Token{
		ID:        trow.ID,
		User:      user,
		TokenHash: trow.TokenHash,
		ExpiresAt: trow.ExpiresAt,
	}, nil
}

func (tr *TokenRepo) DeleteByHash(ctx context.Context, tokenHash string) error {
	query := `DELETE FROM tokens WHERE refresh_token = $1`
	_, err := tr.db.ExecContext(ctx, query, tokenHash)
	return err
}

func (tr *TokenRepo) DeleteAllByUserID(ctx context.Context, userID int64) error {
	query := `DELETE FROM TOKENS WHERE user_id = $1`
	_, err := tr.db.ExecContext(ctx, query, userID)
	return err
}
