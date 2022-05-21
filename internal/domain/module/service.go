package module

import (
	"context"

	"go.uber.org/zap"
)

const ModulesPerStage = 4

type Service interface {
	ListOpenModulesMeta(ctx context.Context, courseId string, amount int) ([]ModuleMeta, error)
	// ModulesPreview(ctx context.Context, amount int) ([][]ModuleMeta, error)
	GetModuleByID(ctx context.Context, courseId string, id int) (Module, error)
	// used by ContentManager
	CreateModules(ctx context.Context, modules []Module) error
}

type service struct {
	storage Storage
	log     *zap.Logger
}

func NewService(st Storage, log *zap.Logger) Service {
	return &service{
		storage: st,
		log:     log,
	}
}

func (s *service) ListOpenModulesMeta(ctx context.Context, courseId string, amount int) ([]ModuleMeta, error) {
	modules, err := s.storage.GetModulesMeta(ctx, courseId, amount)
	if err != nil {
		return nil, err
	}
	return modules, nil
}

func (s *service) GetModuleByID(ctx context.Context, courseId string, id int) (Module, error) {
	return s.storage.GetModule(ctx, courseId, id)
}

func (s *service) CreateModules(ctx context.Context, modules []Module) error {
	return s.storage.SaveModules(ctx, modules)
}
