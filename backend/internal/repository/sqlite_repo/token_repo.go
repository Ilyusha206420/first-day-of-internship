package sqlite_repo

import (
	"context"

	"github.com/Ilyusha206420/first-day-of-internship/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type TokenRepo struct {
	db *sqlx.DB
}

func NewTokenRepo(db *sqlx.DB) *TokenRepo {
	return &TokenRepo{db: db}
}

func (*TokenRepo) Add(ctx context.Context, token *domain.Token) error {
	return nil
}

func (*TokenRepo) DeleteAllByUserID(ctx context.Context, id int64) error {
	return nil
}

func (*TokenRepo) DeleteByHash(ctx context.Context, hash string) error {
	return nil
}

func (*TokenRepo) FindByHash(ctx context.Context, hash string) (*domain.Token, error) {
	return nil, nil
}
