package comment

import (
	"context"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const formatString = "02 Jan 06 15:04 MST"

type Service interface {
	CreateComment(ctx context.Context, user, text string, module, part int) (*CommentDTO, error)
	UpdateComment(ctx context.Context, user, id, text string) (*CommentDTO, error)
	DeleteCommentByID(ctx context.Context, user, id string) error
	ListCommentsByModule(ctx context.Context, user string, module int) ([]*CommentDTO, error)
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

func (s *service) CreateComment(ctx context.Context, user, text string, module, part int) (*CommentDTO, error) {
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
	return getCommentDTO(c), nil
}
func (s *service) UpdateComment(ctx context.Context, user, id, text string) (*CommentDTO, error) {
	c, err := s.storage.GetCommentByID(ctx, user, id)
	if err != nil {
		return nil, err
	}
	c.Text = text
	c.ModifiedAt = time.Now()
	err = s.storage.UpdateComment(ctx, c)
	if err != nil {
		return nil, err
	}
	return getCommentDTO(c), nil
}
func (s *service) DeleteCommentByID(ctx context.Context, user, id string) error {
	return s.storage.DeleteCommentByID(ctx, id)
}
func (s *service) ListCommentsByModule(ctx context.Context, user string, module int) ([]*CommentDTO, error) {
	comments, err := s.storage.ListCommentsByModule(ctx, user, module)
	if err != nil {
		return nil, err
	}
	res := make([]*CommentDTO, len(comments))
	for i, c := range comments {
		res[i] = getCommentDTO(c)
	}
	return res, nil
}

func getCommentDTO(c *Comment) *CommentDTO {
	return &CommentDTO{
		Id:         c.Id,
		ModuleId:   c.ModuleId,
		PartId:     c.PartId,
		ModifiedAt: c.ModifiedAt.Format(formatString),
		Text:       c.Text,
	}
}
