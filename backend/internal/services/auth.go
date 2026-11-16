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

func (s *AuthService) ValidateUserCredentials(username, password string) (bool, error) {
	user, err := s.UserRepo.FindByName(context.Background(), username)
	if err != nil || user == nil {
		return false, err
	}
	return user.ValidatePassword(password)
}
