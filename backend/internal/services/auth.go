package services

import (
	"context"

	"github.com/Ilyusha206420/first-day-of-internship/backend/internal/domain"
)

type AuthService struct {
	UserRepo domain.UserRepository
}

func NewAuthService(userRepo domain.UserRepository) *AuthService {
	return &AuthService{UserRepo: userRepo}
}

func (as *AuthService) AuthWithPassword(ctx context.Context, username string, password string) (bool, error) {
	u, err := as.UserRepo.FindByName(ctx, username)
	if err != nil {
		return false, err
	}

	return u.ValidatePassword(password)
}
