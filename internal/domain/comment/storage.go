package comment

import "context"

type Storage interface {
	GetCommentByID(ctx context.Context, user, id string) (*Comment, error)
	CreateComment(ctx context.Context, c *Comment) error
	UpdateComment(ctx context.Context, c *Comment) error
	DeleteCommentByID(ctx context.Context, id string) error
	ListCommentsByModule(ctx context.Context, user string, module int) ([]*Comment, error)
}
