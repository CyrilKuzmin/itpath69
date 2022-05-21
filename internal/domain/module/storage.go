package module

import "context"

type Storage interface {
	GetModulesMeta(ctx context.Context, courseId string, amount int) ([]ModuleMeta, error)
	GetModule(ctx context.Context, courseId string, id int) (Module, error)
	// used by ContentManager
	SaveModules(ctx context.Context, modules []Module) error
}
