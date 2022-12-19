package course

import (
	"context"
	"sort"

	"go.uber.org/zap"
)

type Service interface {
	CreateCourse(ctx context.Context, c *Course) error
	UpdateCourse(ctx context.Context, c *Course) error
	DeleteCourse(ctx context.Context, id string) error
	GetCourseByID(ctx context.Context, id string) (*Course, error)
	ListCourses(ctx context.Context) ([]string, error)
	ListCoursesByOwner(ctx context.Context, userId string) ([]string, error)
	MakePrivate(ctx context.Context, id string) error
	MakePublic(ctx context.Context, id string) error
	Publish(ctx context.Context, id string) error
	AddOwner(ctx context.Context, id, userId string) error
}

type service struct {
	storage Storage
	log     *zap.Logger
}

func NewService(st Storage, log *zap.Logger) Service {
	return &service{st, log}
}

func (s *service) CreateCourse(ctx context.Context, c *Course) error {
	return s.storage.CreateCourse(ctx, c)
}
func (s *service) UpdateCourse(ctx context.Context, c *Course) error {
	return s.storage.UpdateCourse(ctx, c)
}
func (s *service) DeleteCourse(ctx context.Context, id string) error {
	return s.storage.DeleteCourse(ctx, id)
}
func (s *service) GetCourseByID(ctx context.Context, id string) (*Course, error) {
	c, err := s.storage.GetCourseByID(ctx, id)
	if err != nil {
		return nil, err
	}
	sort.Slice(c.Stages, func(i, j int) bool { return c.Stages[i].ID < c.Stages[j].ID })
	return c, nil
}
func (s *service) ListCourses(ctx context.Context) ([]string, error) {
	return s.storage.ListCourses(ctx)
}
func (s *service) ListCoursesByOwner(ctx context.Context, userId string) ([]string, error) {
	return s.storage.ListCoursesByOwner(ctx, userId)
}
func (s *service) MakePrivate(ctx context.Context, id string) error {
	return s.storage.MakePrivate(ctx, id)
}
func (s *service) MakePublic(ctx context.Context, id string) error {
	return s.storage.MakePublic(ctx, id)
}
func (s *service) Publish(ctx context.Context, id string) error {
	return s.storage.Publish(ctx, id)
}
func (s *service) AddOwner(ctx context.Context, id, userId string) error {
	return s.storage.AddOwner(ctx, id, userId)
}
