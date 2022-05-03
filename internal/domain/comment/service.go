package comment

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service interface {
	CreateComment(ctx context.Context, user, text string, module, part int) (*Comment, error)
	UpdateComment(ctx context.Context, user, id, text string) error
	DeleteCommentByID(ctx context.Context, user, id string) error
	ListCommentsByModule(ctx context.Context, user string, module int) ([]*Comment, error)
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

func (s *service) CreateComment(ctx context.Context, user, text string, module, part int) (*Comment, error) {
	c := &Comment{
		Id:         uuid.NewString(),
		User:       user,
		ModuleId:   module,
		PartId:     part,
		ModifiedAt: time.Now(),
		Text:       text,
	}
	err := s.storage.CreateComment(ctx, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
func (s *service) UpdateComment(ctx context.Context, user, id, text string) error {
	c, err := s.storage.GetCommentByID(ctx, user, id)
	if err != nil {
		return err
	}
	c.Text = text
	c.ModifiedAt = time.Now()
	return s.storage.UpdateComment(ctx, c)
}
func (s *service) DeleteCommentByID(ctx context.Context, user, id string) error {
	return s.storage.DeleteCommentByID(ctx, id)
}
func (s *service) ListCommentsByModule(ctx context.Context, user string, module int) ([]*Comment, error) {
	return s.storage.ListCommentsByModule(ctx, user, module)
}
