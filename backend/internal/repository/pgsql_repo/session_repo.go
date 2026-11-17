package pgsql_repo

import "github.com/jmoiron/sqlx"

type SessionRepository struct {
	db *sqlx.DB
}

type SessionRow struct {
	ID           int64  `db:"id"`
	UserID       int64  `db:"user_id"`
	RefreshToken string `db:"refresh_token"`
	AccessToken  string `db:"access_token"`
}

func NewSessionRepository(db *sqlx.DB) *SessionRepository {
	return &SessionRepository{db: db}
}
