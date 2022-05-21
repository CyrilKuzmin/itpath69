package comment

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service interface {
	CreateComment(ctx context.Context, args CreateCommentArgs) (*Comment, error)
	UpdateComment(ctx context.Context, userId, id, text string) (*Comment, error)
	DeleteCommentByID(ctx context.Context, userId, id string) error
	ListCommentsByModule(ctx context.Context, userId, courseID string, module int) ([]*Comment, error)
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

type CreateCommentArgs struct {
	UserID   string
	CourseID string
	ModuleID int
	PartID   int
	Text     string
}

func (s *service) CreateComment(ctx context.Context, args CreateCommentArgs) (*Comment, error) {
	c := &Comment{
		ID:         uuid.NewString(),
		UserID:     args.UserID,
		CourseID:   args.CourseID,
		ModuleID:   args.ModuleID,
		PartID:     args.PartID,
		ModifiedAt: time.Now(),
		Text:       args.Text,
	}
	err := s.storage.CreateComment(ctx, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
func (s *service) UpdateComment(ctx context.Context, userId, id, text string) (*Comment, error) {
	c, err := s.storage.GetCommentByID(ctx, userId, id)
	if err != nil {
		return nil, err
	}
	c.Text = text
	c.ModifiedAt = time.Now()
	err = s.storage.UpdateComment(ctx, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
func (s *service) DeleteCommentByID(ctx context.Context, userId, id string) error {
	return s.storage.DeleteCommentByID(ctx, id)
}
func (s *service) ListCommentsByModule(ctx context.Context, userId, courseID string, module int) ([]*Comment, error) {
	return s.storage.ListCommentsByModule(ctx, userId, courseID, module)

}
