package domain

import "errors"

type User struct {
	ID       int64
	Name     string
	Password string //TODO: change to hash
}

type UserRepository interface {
	FindByName(name string) (*User, error)               //TODO: Context
	ChangePassword(user *User, newPassword string) error //TODO: Context
}

func (u *User) ValidatePassword(password string) bool {
	return u.Password == password //TODO: change to hash
}

func (u *User) ChangePassword(newPassword string) error {
	if len(newPassword) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	u.Password = newPassword //TODO: change to hash
	return nil
}
