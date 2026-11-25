package sqlite_repo

import (
	"context"

	"github.com/Ilyusha206420/first-day-of-internship/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (*UserRepository) FindByName(ctx context.Context, name string) (*domain.User, error) {
	return nil, nil
}

func (*UserRepository) Add(ctx context.Context, u *domain.User) error {
	return nil
}

func (*UserRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	return nil, nil
}
