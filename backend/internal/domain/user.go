package domain

import (
	"context"
	"errors"

	"github.com/dhenkes/argon2id"
)

type User struct {
	ID           int64
	Name         string
	Password     string
	RefreshToken string
}

type UserRepository interface {
	FindByName(ctx context.Context, name string) (*User, error)
	ChangePassword(ctx context.Context, user *User, newPassword string) error
	Save(ctx context.Context, user *User) error
}

type TokensRepository interface {
	Add(ctx context.Context, user_id int64) error
}

func NewUser(name, password string) (*User, error) {
	if len(password) < 8 {
		return nil, errors.New("password must be at least 8 characters long")
	} else if len(password) > 64 {
		return nil, errors.New("password must be at most 64 characters long")
	}

	hash, err := argon2id.Hash(password, argon2id.DefaultParams())
	if err != nil {
		return nil, err
	}
	return &User{Name: name, Password: hash}, nil
}

func (u *User) ValidatePassword(password string) (bool, error) {
	ok, err := argon2id.Verify(password, u.Password)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (u *User) ChangePassword(newPassword string) error {
	if len(newPassword) < 8 {
		return errors.New("password must be at least 8 characters long")
	} else if len(newPassword) > 64 {
		return errors.New("password must be at most 64 characters long")
	}
	hash, err := argon2id.Hash(newPassword, argon2id.DefaultParams())
	if err != nil {
		return errors.New("failed to hash new password")
	}
	u.Password = hash
	return nil
}
