package pgsql_repo

import (
	"context"

	"github.com/Ilyusha206420/first-day-of-internship/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	ID       int64  `db:"id"`
	Name     string `db:"name"`
	Password string `db:"password"`
	db       *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) FindByName(ctx context.Context, name string) (*domain.User, error) {
	err := ur.db.GetContext(ctx, ur, "SELECT id, name, password FROM users WHERE name = $1", name)
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:       ur.ID,
		Name:     ur.Name,
		Password: ur.Password,
	}, nil
}

func (ur *UserRepository) ChangePassword(ctx context.Context, user *domain.User, newPassword string) error {
	_, err := ur.db.ExecContext(ctx, "UPDATE users SET password = $1 WHERE id = $2", newPassword, user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) Save(ctx context.Context, user *domain.User) error {
	_, err := ur.db.ExecContext(ctx, "INSERT INTO users (name, password) VALUES ($1, $2)", user.Name, user.Password)
	if err != nil {
		return err
	}
	return nil
}
