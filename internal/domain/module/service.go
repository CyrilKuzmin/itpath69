package module

import (
	"context"

	"go.uber.org/zap"
)

const ModulesPerStage = 4

type Service interface {
	ListOpenModulesMeta(ctx context.Context, amount int) ([]ModuleMeta, error)
	ModulesPreview(ctx context.Context, amount int) ([][]ModuleMeta, error)
	GetModuleByID(ctx context.Context, id int) (Module, error)
	// used by ContentManager
	CreateModules(ctx context.Context, modules []Module) error
	ModulesTotal() int
}

type service struct {
	storage      Storage
	log          *zap.Logger
	totalModules int
}

func NewService(st Storage, log *zap.Logger) Service {
	return &service{
		storage: st,
		log:     log,
	}
}

func (s *service) ListOpenModulesMeta(ctx context.Context, amount int) ([]ModuleMeta, error) {
	modules, err := s.storage.GetModulesMeta(ctx, amount)
	if err != nil {
		return nil, err
	}
	return modules, nil
}

func (s *service) ModulesPreview(ctx context.Context, amount int) ([][]ModuleMeta, error) {
	// An order has sense is forPreview is true. We want to see last modules on the top.
	// for ex. 1 2 3 4 5 6 7 8 = > 5 6 7 8 1 2 3 4
	metas, err := s.storage.GetModulesMeta(ctx, amount)
	if err != nil {
		return nil, err
	}
	shiftMetas(metas)
	rowsNum := len(metas) / ModulesPerStage
	if len(metas)%ModulesPerStage != 0 {
		rowsNum++
	}
	rows := make([][]ModuleMeta, rowsNum)
	for i := 0; i < rowsNum; i++ {
		row := []ModuleMeta(metas[i*ModulesPerStage : i*ModulesPerStage+ModulesPerStage])
		rows = append(rows, row)
	}
	return rows, nil
}

func (s *service) GetModuleByID(ctx context.Context, id int) (Module, error) {
	return s.storage.GetModule(ctx, id)
}

func (s *service) CreateModules(ctx context.Context, modules []Module) error {
	s.totalModules = len(modules)
	return s.storage.SaveModules(ctx, modules)
}

func (s *service) ModulesTotal() int {
	return s.totalModules
}

func shiftMetas(in []ModuleMeta) {
	for in[len(in)-ModulesPerStage].Id != 1 {
		for k := 0; k < ModulesPerStage; k++ {
			less := in[len(in)-1]
			for i := len(in) - 1; i >= 0; i-- {
				if i == 0 {
					in[i] = less
					continue
				}
				in[i] = in[i-1]
			}
		}
	}
}
