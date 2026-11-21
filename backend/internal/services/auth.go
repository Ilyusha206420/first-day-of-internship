package services

import (
	"context"
	"errors"

	"github.com/Ilyusha206420/first-day-of-internship/backend/internal/domain"
)

var (
	InvalidCredentials = errors.New("Invalid username or password")
)

type SessionManager interface {
	CreateSession(ctx context.Context, user *domain.User) (*TokenPair, error)
}

type AuthService struct {
	userRepo       domain.UserRepository
	sessionManager SessionManager
}

func NewAuthService(userRepo domain.UserRepository, sessionManager SessionManager) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		sessionManager: sessionManager,
	}
}

func (as *AuthService) AuthWithLoginPassword(ctx context.Context, username string, password string) (*TokenPair, error) {
	u, err := as.userRepo.FindByName(ctx, username)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, InvalidCredentials
	}

	valid, err := u.ValidatePassword(password)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, InvalidCredentials
	}

	return as.sessionManager.CreateSession(ctx, u)
}
