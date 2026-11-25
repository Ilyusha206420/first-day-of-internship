package rows

import (
	"time"
	"uuid"
)

type TokenRow struct {
	ID        int64     `db:"id"`
	UserID    int64     `db:"user_id"`
	TokenHash string    `db:"refresh_token"`
	ExpiresAt time.Time `db:"expires_at"`
}

type UserRow struct {
	ID       int64     `db:"id"`
	PublicID uuid.UUID `db:"public_id"`
	Name     string    `db:"name"`
	Password string    `db:"pass"`
}
