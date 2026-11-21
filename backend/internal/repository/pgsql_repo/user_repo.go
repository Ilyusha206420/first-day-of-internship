package pgsql_repo

import (
	"context"
	"database/sql"
	"errors"
	"uuid"

	"github.com/Ilyusha206420/first-day-of-internship/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type PGUserRepository struct {
	db *sqlx.DB
}

type userRow struct {
	ID       int64     `db:"id"`
	PublicID uuid.UUID `db:"public_id"`
	Name     string    `db:"name"`
	Password string    `db:"pass"`
}

func NewUserRepository(db *sqlx.DB) *PGUserRepository {
	return &PGUserRepository{db: db}
}

func (ur *PGUserRepository) Add(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (name, pass) 
		VALUES ($1, $2) 
		RETURNING id, public_id`

	return ur.db.QueryRowxContext(ctx, query, user.Name, user.PasswordHash).Scan(&user.ID, &user.PublicID)
}

func (ur *PGUserRepository) FindByName(ctx context.Context, name string) (*domain.User, error) {
	var urow userRow
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

func (ur *PGUserRepository) FindByPublicID(ctx context.Context, publicID uuid.UUID) (*domain.User, error) {
	var urow userRow
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

func (ur *PGUserRepository) ChangePassword(ctx context.Context, user *domain.User, newPasswordHash string) error {
	query := `UPDATE users SET pass = $1 WHERE id = $2`
	_, err := ur.db.ExecContext(ctx, query, newPasswordHash, user.ID)
	return err
}
