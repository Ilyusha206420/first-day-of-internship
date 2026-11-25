package sqlite_repo

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

func (ur *UserRepository) FindByName(ctx context.Context, name string) (*domain.User, error) {
	var urow rows.UserRow
	query := `SELECT id, public_id, name, pass FROM users WHERE name = ?`

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

func (ur *UserRepository) Add(ctx context.Context, u *domain.User) error {
	query := `
		INSERT INTO users (public_id, name, pass) 
		VALUES (?, ?, ?) 
		RETURNING id`

	pubID := uuid.New()

	return ur.db.QueryRowxContext(ctx, query, pubID, u.Name, u.PasswordHash).Scan(&u.ID)
}

func (ur *UserRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	var urow rows.UserRow
	query := `SELECT id, public_id, name, pass FROM users WHERE id = ?`

	err := ur.db.GetContext(ctx, &urow, query, id)
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
