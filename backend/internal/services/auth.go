package services

import (
	"context"
	"errors"

	"github.com/Ilyusha206420/first-day-of-internship/backend/internal/domain"
)

var (
	InvalidCredentials = errors.New("Invalid username or password")
)

type AuthService struct {
	userRepo domain.UserRepository
}

func NewAuthService(userRepo domain.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (as *AuthService) AuthWithLoginPassword(ctx context.Context, username string, password string) (*domain.User, error) {
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

	return u, nil
}
