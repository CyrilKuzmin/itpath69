package progress

import "context"

type Storage interface {
	CreateCourseProgress(ctx context.Context, data []ModuleProgress) error
	MarkModuleAsCompleted(ctx context.Context, id string) error
	MarkModulesAsOpen(ctx context.Context, toOpenIDs []string) error
	GetUserProgress(ctx context.Context, userId, courseId string) ([]ModuleProgress, error)
	GetModuleProgress(ctx context.Context, userId, courseId string, moduleId int) (ModuleProgress, error)
}
