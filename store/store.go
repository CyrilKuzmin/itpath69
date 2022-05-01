// Store defines the main interface for the storage

package store

import (
	"context"
	"fmt"

	"github.com/CyrilKuzmin/itpath69/models"
)

type Store interface {
	SaveUser(ctx context.Context, username, password string) error
	GetUser(ctx context.Context, username, password string) (*models.User, error)
	Close(ctx context.Context)
}

type StoreError error

func ErrInternal(err error) StoreError {
	return fmt.Errorf("internal storage error %w", err)
}

func ErrUserAlreadyExists(username string) StoreError {
	return fmt.Errorf("user %v already exists", username)
}

func ErrUserNotFound(username string) StoreError {
	return fmt.Errorf("user %v not found", username)
}
