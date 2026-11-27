package sqlite_repo

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
		VALUES (?, ?, ?)`
	_, err := tr.db.ExecContext(ctx, query, token.User.ID, token.TokenHash, token.ExpiresAt)
	return err
}

func (tr *TokenRepo) DeleteAllByUserID(ctx context.Context, id int64) error {
	query := `DELETE FROM tokens WHERE user_id = ?`
	_, err := tr.db.ExecContext(ctx, query, id)
	return err
}

func (tr *TokenRepo) DeleteByHash(ctx context.Context, hash string) error {
	query := `DELETE FROM tokens WHERE refresh_token = ?`
	_, err := tr.db.ExecContext(ctx, query, hash)
	return err
}

func (tr *TokenRepo) FindByHash(ctx context.Context, hash string) (*domain.Token, error) {
	var trow rows.TokenRow
	query := `SELECT id, user_id, refresh_token, expires_at FROM tokens WHERE refresh_token = ?`

	err := tr.db.GetContext(ctx, &trow, query, hash)
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
