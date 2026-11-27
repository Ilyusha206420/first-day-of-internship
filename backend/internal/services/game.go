package services

import (
	"sync"

	"github.com/Ilyusha206420/first-day-of-internship/backend/internal/domain"
)

type GameNodeRepository interface {
}

type GameService struct {
	mu        sync.RWMutex
	userNodes map[int64]*domain.Node
	NodeRepo  GameNodeRepository
}

func NewGameService(nodeRepo GameNodeRepository) *GameService {
	return &GameService{
		userNodes: make(map[int64]*domain.Node),
		NodeRepo:  nodeRepo,
	}
}

func (gs *GameService) Start(u *domain.User) error {
	return nil
}

func (gs *GameService) PlayerMove() error {
	return nil
}
