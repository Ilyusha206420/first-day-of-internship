package internal

import (
	"context"
	"database/sql"
)

type Content struct {
	ID   int64  `db:"id" json:"id"`
	Data string `db:"data" json:"data"`
}

type ContentRepository interface {
	GetContentByID(ctx context.Context, id int64) (*Content, error)
}

type Node struct {
	ID        int64         `db:"id" json:"id"`
	Method    int64         `db:"method" json:"method"`
	ContentID sql.NullInt64 `db:"content_id" json:"content_id"`
}

type NodeRepository interface {
	GetNodeByID(ctx context.Context, id int64) (*Node, error)
}

type Edge struct {
	SrcID int64 `db:"src_id" json:"src_id"`
	TgtID int64 `db:"tgt_id" json:"tgt_id"`
}

type EdgeRepository interface {
	GetTgtsBySrcId(ctx context.Context, src_id int64) ([]*Edge, error)
}

type User struct {
	ID int64 `db:"id" json:"id"`
}

type UserRepository interface {
	GetUserByID(ctx context.Context, id int64) (*User, error)
}

type Session struct {
	ID     int64 `db:"id" json:"id"`
	UserID int64 `db:"user_id" json:"user_id"`
}

type SessionRepository interface {
	GetSessionByID(ctx context.Context, id int64) (*Session, error)
}

type Dialog struct {
	SessionID int64 `db:"session_id" json:"session_id"`
	ContentID int64 `db:"content_id" json:"content_id"`
}

type DialogRepository interface {
	GetDialogsBySessionID(ctx context.Context, session_id int64) ([]*Dialog, error)
}
