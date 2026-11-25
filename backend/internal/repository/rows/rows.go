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

type NodeRow struct {
	ID        int64 `db:"id"`
	MethodID  int64 `db:"method_id"`
	ContentID int64 `db:"content_id"`
}

type EdgeRow struct {
	SrcID int64 `db:"src_id"`
	TgtID int64 `db:"tgt_id"`
}

type ContentRow struct {
	ID   int64  `db:"id"`
	Data string `db:"data"`
}

type MethodsRow struct {
	ID   int64  `db:"id"`
	Name string `db:"method_name"`
}
