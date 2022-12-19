package users

import (
	"context"
)

type Storage interface {
	CreateUser(ctx context.Context, user *User) error
	CheckUserPassword(ctx context.Context, username, passwordHash string) error
	GetUserByName(ctx context.Context, username string) (*User, error)
}
