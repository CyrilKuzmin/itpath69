package users

import (
	"context"
)

type Storage interface {
	CreateUser(ctx context.Context, username, password string) (*User, error)
	CheckUserPassword(ctx context.Context, username, password string) (*User, error)
	GetUser(ctx context.Context, username string) (*User, error)
	// users progress
	UpdateProgress(ctx context.Context, username string, progress map[int]ModuleProgress) error
}
