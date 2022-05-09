package users

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/CyrilKuzmin/itpath69/internal/domain/module"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service interface {
	CreateUser(ctx context.Context, username, password string) (*User, error)
	CheckUserPassword(ctx context.Context, username, password string) error
	GetUserByName(ctx context.Context, username string) (*UserDTO, error)
	// users progress
	OpenNewModules(ctx context.Context, username string) error
	MarkModuleAsCompleted(ctx context.Context, username string, moduleId int) error
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
		Id:           uuid.New().String(),
		Username:     username,
		PasswordHash: hex.EncodeToString(pwMd5[:]),
		CreatedAt:    time.Now(),
		Modules:      map[int]ModuleProgress{},
	}
	err := s.storage.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}
	currTime := time.Now()
	for i := 1; i <= module.ModulesPerStage; i++ {
		user.Modules[i] = ModuleProgress{CreatedAt: currTime}
	}
	err = s.storage.UpdateProgress(ctx, username, user.Modules)
	return user, err
}

func (s *service) CheckUserPassword(ctx context.Context, username, password string) error {
	pwMd5 := md5.Sum([]byte(password))
	hash := hex.EncodeToString(pwMd5[:])
	return s.storage.CheckUserPassword(ctx, username, hash)
}

func (s *service) GetUserByName(ctx context.Context, username string) (*UserDTO, error) {
	user, err := s.storage.GetUserByName(ctx, username)
	if err != nil {
		return nil, err
	}
	opened := len(user.Modules)
	completed := 0
	for _, mp := range user.Modules {
		if !mp.CompletedAt.IsZero() {
			completed++
		}
	}
	return &UserDTO{
		User:             user,
		ModulesOpened:    opened,
		ModulesCompleted: completed,
	}, nil
}

func (s *service) OpenNewModules(ctx context.Context, username string) error {
	user, _ := s.storage.GetUserByName(ctx, username)
	currTime := time.Now()
	for i := len(user.Modules); i <= len(user.Modules)+module.ModulesPerStage; i++ {
		if _, found := user.Modules[i]; found {
			continue
		}
		user.Modules[i] = ModuleProgress{CreatedAt: currTime}
	}
	return s.storage.UpdateProgress(ctx, username, user.Modules)
}

func (s *service) MarkModuleAsCompleted(ctx context.Context, username string, moduleId int) error {
	currTime := time.Now()
	user, _ := s.storage.GetUserByName(ctx, username)
	opened := len(user.Modules)
	created := user.Modules[moduleId].CreatedAt
	user.Modules[moduleId] = ModuleProgress{CreatedAt: created, CompletedAt: currTime}
	completedOnStage := 0
	for i := opened - module.ModulesPerStage + 1; i <= opened; i++ {
		if !user.Modules[i].CompletedAt.IsZero() {
			completedOnStage++
		}
	}
	if completedOnStage > 2 {
		for i := len(user.Modules) + 1; i <= opened+module.ModulesPerStage; i++ {
			user.Modules[i] = ModuleProgress{CreatedAt: currTime}
		}
	}
	return s.storage.UpdateProgress(ctx, username, user.Modules)
}
