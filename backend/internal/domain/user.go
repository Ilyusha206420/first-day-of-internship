package domain

import (
	"context"
	"errors"

	"github.com/dhenkes/argon2id"
)

type User struct {
	ID       int64
	Name     string
	Password string
}

type UserRepository interface {
	FindByName(ctx context.Context, name string) (*User, error)
	ChangePassword(ctx context.Context, user *User, newPassword string) error
}

func (u *User) ValidatePassword(password string) (bool, error) {
	hash, err := argon2id.Hash(password, argon2id.DefaultParams)
	if err != nil {
		return false, err
	}
	return u.Password == hash, nil
}

func (u *User) ChangePassword(newPassword string) error {
	if len(newPassword) < 8 {
		return errors.New("password must be at least 8 characters long")
	} else if len(newPassword) > 64 {
		return errors.New("password must be at most 64 characters long")
	}
	hash, err := argon2id.Hash(newPassword, argon2id.DefaultParams)
	if err != nil {
		return errors.New("failed to hash new password")
	}
	u.Password = hash
	return nil
}
