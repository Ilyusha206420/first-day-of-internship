package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/Ilyusha206420/first-day-of-internship/backend/internal/domain"
)

var (
	InvalidCredentials = errors.New("Invalid username or password")
)

type AuthService struct {
	userRepo   domain.UserRepository
	tokenRepo  domain.TokenRepository
	jwtSecret  []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type AccessClaims struct {
	jwt.RegisteredClaims
}

func NewAuthService(
	userRepo domain.UserRepository,
	tokenRepo domain.TokenRepository,
	jwtSecret string,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		jwtSecret:  []byte(jwtSecret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (as *AuthService) generateAccessToken(u *domain.User) (string, error) {
	now := time.Now().UTC()

	claims := AccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   u.PublicID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(as.accessTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(as.jwtSecret)
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

	refreshSession, rawRefreshToken, err := domain.NewToken(u.ID, as.refreshTTL)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	if err := as.tokenRepo.Add(ctx, refreshSession); err != nil {
		return nil, fmt.Errorf("save refresh token: %w", err)
	}

	accessToken, err := as.generateAccessToken(u)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: rawRefreshToken,
	}, nil
}

func (as *AuthService) AuthWithCookie(ctx context.Context, refreshToken string) (*TokenPair, error) {
	return nil, nil
}
