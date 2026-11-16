package pgsql_repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Ilyusha206420/first-day-of-internship/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

type userRow struct {
	ID       int64  `db:"id"`
	Name     string `db:"name"`
	Password string `db:"pass"`
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) FindByName(ctx context.Context, name string) (*domain.User, error) {
	var urow userRow
	err := ur.db.GetContext(ctx, &urow, "SELECT id, name, pass FROM users WHERE name = $1", name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &domain.User{
		ID:       urow.ID,
		Name:     urow.Name,
		Password: urow.Password,
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
	_, err := ur.db.ExecContext(ctx, "INSERT INTO users (name, pass) VALUES ($1, $2)", user.Name, user.Password)
	if err != nil {
		return err
	}
	return nil
}
