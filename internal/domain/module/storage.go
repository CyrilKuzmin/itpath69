package module

import "context"

type Storage interface {
	GetModulesMeta(ctx context.Context, amount int) ([]ModuleMeta, error)
	GetModule(ctx context.Context, id int) (Module, error)
	// used by ContentManager
	SaveModules(ctx context.Context, modules []Module) error
}
