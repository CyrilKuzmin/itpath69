package course

import "context"

type Storage interface {
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
