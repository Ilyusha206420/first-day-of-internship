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
