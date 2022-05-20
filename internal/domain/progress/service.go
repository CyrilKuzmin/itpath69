package progress

import (
	"context"
	"sort"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service interface {
	CreateCourseProgress(ctx context.Context, userId, courseId string, total, open int) ([]ModuleProgress, error)
	OpenNewModules(ctx context.Context, userId, courseId string, amount int) error
	MarkModuleAsCompleted(ctx context.Context, userId, courseId string, moduleId int) error
	GetUserProgress(ctx context.Context, userId, courseId string) ([]ModuleProgress, error)
}

type service struct {
	storage Storage
	log     *zap.Logger
}

func NewService(st Storage, log *zap.Logger) Service {
	return &service{st, log}
}

func (s *service) CreateCourseProgress(ctx context.Context, userId, courseId string, modulesTotal, modulesOpen int) ([]ModuleProgress, error) {
	res := make([]ModuleProgress, modulesTotal)
	now := time.Now()
	for i := 0; i < modulesTotal; i++ {
		pr := ModuleProgress{
			ID:        uuid.NewString(),
			UserID:    userId,
			CourseID:  courseId,
			ModuleID:  i,
			CreatedAt: now,
		}
		if modulesOpen > 0 {
			pr.OpenedAt = now
			modulesOpen--
		}
		res[i] = pr
	}
	err := s.storage.CreateCourseProgress(ctx, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *service) MarkModuleAsCompleted(ctx context.Context, userId, courseId string, moduleId int) error {
	pr, err := s.storage.GetModuleProgress(ctx, userId, courseId, moduleId)
	if err != nil {
		return err
	}
	return s.storage.MarkModuleAsCompleted(ctx, pr.ID)
}

func (s *service) OpenNewModules(ctx context.Context, userId, courseId string, amount int) error {
	curr, err := s.GetUserProgress(ctx, userId, courseId)
	toOpenIDs := make([]string, 0)
	if err != nil {
		return err
	}
	for i := range curr {
		if amount == 0 {
			break
		}
		if curr[i].OpenedAt.IsZero() {
			toOpenIDs = append(toOpenIDs, curr[i].ID)
		}
	}
	return s.storage.MarkModulesAsOpen(ctx, toOpenIDs)
}

func (s *service) GetUserProgress(ctx context.Context, userId, courseId string) ([]ModuleProgress, error) {
	res, err := s.storage.GetUserProgress(ctx, userId, courseId)
	if err != nil {
		return nil, err
	}
	sort.Slice(res, func(i, j int) bool { return res[i].ModuleID < res[j].ModuleID })
	return res, nil
}
