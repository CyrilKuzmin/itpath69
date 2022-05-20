package users

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const DefaultCourseID = "default"

type Service interface {
	CreateUser(ctx context.Context, username, password string) (*User, error)
	CheckUserPassword(ctx context.Context, username, password string) error
	GetUserByName(ctx context.Context, username string) (*User, error)
}

type service struct {
	storage Storage
	log     *zap.Logger
}

func NewService(st Storage, log *zap.Logger) Service {
	return &service{st, log}
}

func (s *service) CreateUser(ctx context.Context, username, password string) (*User, error) {
	pwMd5 := md5.Sum([]byte(password))
	user := &User{
		Id:            uuid.New().String(),
		Username:      username,
		PasswordHash:  hex.EncodeToString(pwMd5[:]),
		CreatedAt:     time.Now(),
		CurrentCourse: DefaultCourseID,
		CurrentStage:  1,
	}
	err := s.storage.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (s *service) CheckUserPassword(ctx context.Context, username, password string) error {
	pwMd5 := md5.Sum([]byte(password))
	hash := hex.EncodeToString(pwMd5[:])
	return s.storage.CheckUserPassword(ctx, username, hash)
}

func (s *service) GetUserByName(ctx context.Context, username string) (*User, error) {
	return s.storage.GetUserByName(ctx, username)
}
