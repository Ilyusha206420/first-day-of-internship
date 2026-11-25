package pgsql_repo

import (
	"context"
	"database/sql"
	"errors"
	"uuid"

	"github.com/Ilyusha206420/first-day-of-internship/backend/internal/domain"
	"github.com/Ilyusha206420/first-day-of-internship/backend/internal/repository/rows"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) Add(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (name, pass) 
		VALUES ($1, $2) 
		RETURNING id, public_id`

	return ur.db.QueryRowxContext(ctx, query, user.Name, user.PasswordHash).Scan(&user.ID, &user.PublicID)
}

func (ur *UserRepository) FindByID(ctx context.Context, user_id int64) (*domain.User, error) {
	var urow rows.UserRow
	query := `SELECT id, public_id, name, pass FROM users WHERE id = $1`

	err := ur.db.GetContext(ctx, &urow, query, user_id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &domain.User{
		ID:           urow.ID,
		PublicID:     urow.PublicID,
		Name:         urow.Name,
		PasswordHash: urow.Password,
	}, nil
}

func (ur *UserRepository) FindByName(ctx context.Context, name string) (*domain.User, error) {
	var urow rows.UserRow
	query := `SELECT id, public_id, name, pass FROM users WHERE name = $1`

	err := ur.db.GetContext(ctx, &urow, query, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &domain.User{
		ID:           urow.ID,
		PublicID:     urow.PublicID,
		Name:         urow.Name,
		PasswordHash: urow.Password,
	}, nil
}

func (ur *UserRepository) FindByPublicID(ctx context.Context, publicID uuid.UUID) (*domain.User, error) {
	var urow rows.UserRow
	query := `SELECT id, public_id, name, pass FROM users WHERE public_id = $1`

	err := ur.db.GetContext(ctx, &urow, query, publicID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &domain.User{
		ID:           urow.ID,
		PublicID:     urow.PublicID,
		Name:         urow.Name,
		PasswordHash: urow.Password,
	}, nil
}

func (ur *UserRepository) ChangePassword(ctx context.Context, user_id int64, newPasswordHash string) error {
	query := `UPDATE users SET pass = $1 WHERE id = $2`
	_, err := ur.db.ExecContext(ctx, query, newPasswordHash, user_id)
	return err
}
