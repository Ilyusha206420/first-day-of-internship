package services

import (
	"context"
	"errors"
	"fmt"
	"time"
	"uuid"

	"github.com/Ilyusha206420/first-day-of-internship/backend/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

var (
	InvalidToken = errors.New("Invalid or expired token")
)

type SessionUserRepository interface {
	FindByID(ctx context.Context, user_id int64) (*domain.User, error)
}

type SessionTokenRepository interface {
}

type AccessClaims struct {
	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type SessionService struct {
	userRepo   SessionUserRepository
	tokenRepo  SessionTokenRepository
	jwtSecret  []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewSessionService(
	userRepo SessionUserRepository,
	tokenRepo SessionTokenRepository,
	jwtSecret string,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) *SessionService {
	return &SessionService{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		jwtSecret:  []byte(jwtSecret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (ss *SessionService) CreateSession(ctx context.Context, u *domain.User) (*TokenPair, error) {
	accessToken, err := ss.generateAccessToken(u)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refreshSession, rawRefreshToken, err := domain.NewToken(u.ID, ss.refreshTTL)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	if err := ss.tokenRepo.Add(ctx, refreshSession); err != nil {
		return nil, fmt.Errorf("save refresh session: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: rawRefreshToken,
	}, nil
}

func (ss *SessionService) RefreshAccessToken(ctx context.Context, rawRefreshToken string) (string, error) {
	if rawRefreshToken == "" {
		return "", InvalidToken
	}

	tokenHash := domain.HashToken(rawRefreshToken)
	session, err := ss.tokenRepo.FindByHash(ctx, tokenHash)
	if err != nil {
		return "", fmt.Errorf("find refresh session: %w", err)
	}
	if session == nil {
		return "", InvalidToken
	}

	if session.IsExpired() {
		_ = ss.tokenRepo.DeleteByHash(ctx, tokenHash)
		return "", InvalidToken
	}

	u, err := ss.userRepo.FindByID(ctx, session.UserID)
	if err != nil {
		return "", fmt.Errorf("find user by id: %w", err)
	}
	if u == nil {
		return "", InvalidToken
	}

	accessToken, err := ss.generateAccessToken(u)
	if err != nil {
		return "", fmt.Errorf("generate access token: %w", err)
	}

	return accessToken, nil
}

func (ss *SessionService) ValidateAccessToken(tokenString string) (uuid.UUID, error) {
	if tokenString == "" {
		return uuid.Nil(), InvalidToken
	}

	var claims AccessClaims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return ss.jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return uuid.Nil(), InvalidToken
	}

	publicID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil(), InvalidToken
	}

	return publicID, nil
}

func (ss *SessionService) RevokeSession(ctx context.Context, rawRefreshToken string) error {
	if rawRefreshToken == "" {
		return nil
	}
	tokenHash := domain.HashToken(rawRefreshToken)
	return ss.tokenRepo.DeleteByHash(ctx, tokenHash)
}

func (ss *SessionService) RevokeAllSessions(ctx context.Context, userID int64) error {
	return ss.tokenRepo.DeleteAllByUserID(ctx, userID)
}

func (ss *SessionService) generateAccessToken(u *domain.User) (string, error) {
	now := time.Now().UTC()

	claims := AccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   u.PublicID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ss.accessTTL)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(ss.jwtSecret)
}
