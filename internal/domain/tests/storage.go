package tests

import "context"

type Storage interface {
	// Tests
	GetTestByID(ctx context.Context, id string) (*Test, error)
	GetTestsByUser(ctx context.Context, userId string) ([]*Test, error)
	SaveTest(ctx context.Context, test *Test) error
	MarkTestExpired(ctx context.Context, id string) error
	// Questions
	GetModuleQuestions(ctx context.Context, moduleId int, amount int) ([]*Question, error)
	// Content Manager method
	SaveQuestions(ctx context.Context, qs []Question) error
}
