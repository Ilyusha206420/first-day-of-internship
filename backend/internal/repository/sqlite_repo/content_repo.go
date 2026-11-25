package sqlite_repo

import (
	"github.com/Ilyusha206420/first-day-of-internship/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type ContentRepository struct {
	DB *sqlx.DB
}

func (cr *ContentRepository) FindByData(data string) (*domain.Content, error) {
	return nil, nil
}
