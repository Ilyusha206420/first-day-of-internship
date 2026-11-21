package domain

import (
	"context"
	"errors"
	"uuid"

	"github.com/dhenkes/argon2id"
)

type User struct {
	ID           int64
	PublicID     uuid.UUID
	Name         string
	PasswordHash string
}

type UserRepository interface {
	FindByName(ctx context.Context, name string) (*User, error)
	FindByPublicID(ctx context.Context, publicID uuid.UUID) (*User, error)
	ChangePassword(ctx context.Context, userID int64, newPasswordHash string) error
	Save(ctx context.Context, user *User) error
}

func NewUser(name, password string) (*User, error) {
	hash, err := validateAndHashPassword(password)
	if err != nil {
		return nil, err
	}
	return &User{Name: name, PasswordHash: hash}, nil
}

func (u *User) ValidatePassword(password string) (bool, error) {
	return argon2id.Verify(password, u.PasswordHash)
}

func (u *User) ChangePassword(newPassword string) error {
	hash, err := validateAndHashPassword(newPassword)
	if err != nil {
		return err
	}
	u.PasswordHash = hash
	return nil
}

func validateAndHashPassword(password string) (string, error) {
	if len(password) < 8 {
		return "", errors.New("password must be at least 8 characters long")
	}
	if len(password) > 64 {
		return "", errors.New("password must be at most 64 characters long")
	}
	return argon2id.Hash(password, argon2id.DefaultParams())
}
