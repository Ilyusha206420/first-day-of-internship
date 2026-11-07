package domain

import "context"

type Content struct {
	ID   int64
	Data string
}

type ContentRepository interface {
	GetContentByID(ctx context.Context, id int64) (*Content, error)
}

type Node struct {
	ID        int64
	Method    int64
	ContentID int64
}

type NodeRepository interface {
	GetNodeByID(ctx context.Context, id int64) (*Node, error)
}

type Edge struct {
	SrcID int64
	TgtID int64
}

type EdgeRepository interface {
	GetTgtsBySrcId(ctx context.Context, src_id int64) ([]*Edge, error)
}
