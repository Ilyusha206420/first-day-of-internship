package sqlite_repo

import (
	"context"

	"github.com/Ilyusha206420/first-day-of-internship/backend/internal/domain"
	"github.com/Ilyusha206420/first-day-of-internship/backend/internal/repository/rows"
	"github.com/jmoiron/sqlx"
)

type ContentRepository struct {
	db *sqlx.DB
}

func NewContentRepo(db *sqlx.DB) *ContentRepository {
	return &ContentRepository{db: db}
}

func (cr *ContentRepository) FindByData(ctx context.Context, data string) (*domain.Content, error) {
	var crow rows.ContentRow
	query := `SELECT id, data FROM contents WHERE data = ?`

	err := cr.db.GetContext(ctx, &crow, query, data)
	if err != nil {
		return nil, err
	}

	return &domain.Content{
		ID:   crow.ID,
		Data: crow.Data,
	}, nil
}

func (cr *ContentRepository) Add(data string) (int64, error) {
	return 0, nil
}
