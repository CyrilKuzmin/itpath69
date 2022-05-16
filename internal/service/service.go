package service

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"

	"github.com/CyrilKuzmin/itpath69/internal/domain/comment"
	"github.com/CyrilKuzmin/itpath69/internal/domain/module"
	"github.com/CyrilKuzmin/itpath69/internal/domain/tests"
	"github.com/CyrilKuzmin/itpath69/internal/domain/users"
	"go.uber.org/zap"
)

type Service interface {
	// User
	CreateUser(ctx context.Context, username, password string) (*users.User, error)
	GetUserByName(ctx context.Context, username string) (*users.UserDTO, error)
	CheckUserPassword(ctx context.Context, username, password string) error

	// Modules
	GetModuleForUser(ctx context.Context, user *users.UserDTO, moduleId int) (*module.ModuleDTO, error)
	ModulesPreview(ctx context.Context, user *users.UserDTO, amount int) ([]module.ModuleDTO, error)

	// Tests
	CreateNewTest(ctx context.Context, userId string, moduleId int) (*tests.Test, error)
	GetTestByID(ctx context.Context, id string, hideAnswers bool) (*tests.Test, error)
	ListTestsByUser(ctx context.Context, userId string) ([]*tests.Test, error)
	CheckTest(ctx context.Context, username string, userData io.Reader) (*tests.UserResult, error)

	//Comments
	CreateComment(ctx context.Context, username, text string, module, part int) (*comment.CommentDTO, error)
	ListCommentsByModule(ctx context.Context, username string, module int) ([]*comment.CommentDTO, error)
	UpdateComment(ctx context.Context, username, id, text string) (*comment.CommentDTO, error)
	DeleteCommentByID(ctx context.Context, username, id string) error
}

// Storage
type Storage interface {
	comment.Storage
	module.Storage
	tests.Storage
	users.Storage
}

type service struct {
	us users.Service
	ms module.Service
	cs comment.Service
	ts tests.Service
}

func NewService(log *zap.Logger, s Storage) Service {
	// init echo server with its dependencies
	us := users.NewService(s, log)
	ms := module.NewService(s, log)
	cs := comment.NewService(s, log)
	ts := tests.NewService(s, log)
	return &service{us, ms, cs, ts}
}

// simple functions for rendering

// GetUserByName
func (s *service) GetUserByName(ctx context.Context, username string) (*users.UserDTO, error) {
	return s.us.GetUserByName(ctx, username)
}

func (s *service) ListTestsByUser(ctx context.Context, userId string) ([]*tests.Test, error) {
	return s.ts.ListTestsByUser(ctx, userId)
}

func (s *service) ListCommentsByModule(ctx context.Context, username string, module int) ([]*comment.CommentDTO, error) {
	// check if IDs are valid and allowed
	user, err := s.us.GetUserByName(ctx, username)
	if err != nil {
		return nil, err
	}
	if module > user.ModulesOpened {
		return nil, errModuleNotAllowed(module)
	}
	return s.cs.ListCommentsByModule(ctx, username, module)
}

func (s *service) CreateComment(ctx context.Context, username, text string, module, part int) (*comment.CommentDTO, error) {
	// check if IDs are valid and allowed
	user, err := s.us.GetUserByName(ctx, username)
	if err != nil {
		return nil, err
	}
	if module > user.ModulesOpened {
		return nil, errModuleNotAllowed(module)
	}
	m, err := s.ms.GetModuleByID(ctx, module)
	if err != nil {
		return nil, err
	}
	if part > len(m.Data) {
		return nil, errPartNotExists(part)
	}
	return s.cs.CreateComment(ctx, username, text, module, part)
}

func (s *service) UpdateComment(ctx context.Context, username, id, text string) (*comment.CommentDTO, error) {
	return s.cs.UpdateComment(ctx, username, id, text)
}

func (s *service) DeleteCommentByID(ctx context.Context, username, id string) error {
	return s.cs.DeleteCommentByID(ctx, username, id)
}

func (s *service) CreateUser(ctx context.Context, username, password string) (*users.User, error) {
	return s.us.CreateUser(ctx, username, password)
}
func (s *service) CheckUserPassword(ctx context.Context, username, password string) error {
	return s.us.CheckUserPassword(ctx, username, password)
}

func (s *service) CheckTest(ctx context.Context, username string, userData io.Reader) (*tests.UserResult, error) {
	userTestData := &tests.Test{}
	err := json.NewDecoder(userData).Decode(&userTestData)
	if err != nil {
		return nil, err
	}
	user, err := s.us.GetUserByName(ctx, username)
	if err != nil {
		return nil, err
	}
	score, err := s.ts.CheckTest(ctx, userTestData.Id, userTestData.Questions)
	if err != nil {
		return nil, err
	}
	mod, err := s.ms.GetModuleByID(ctx, userTestData.ModuleId)
	var isPassed bool
	if float64(score) >= mod.Meta.TestPassThreshold {
		err = s.us.MarkModuleAsCompleted(ctx, username, userTestData.ModuleId)
		if err != nil {
			return nil, err
		}
		isPassed = true
	}
	if user.ModulesOpened < 8 { // fix it
		err = s.us.OpenNewModules(ctx, username, 4)
		if err != nil {
			return nil, err
		}
	}
	fmt.Printf("marking test %v for user %v as expired\n", userTestData.Id, username)
	err = s.ts.MarkTestExpired(ctx, userTestData.Id)
	if err != nil {
		return nil, err
	}
	return &tests.UserResult{
		Score:    score,
		IsPassed: isPassed,
	}, nil
}

func (s *service) CreateNewTest(ctx context.Context, username string, moduleId int) (*tests.Test, error) {
	user, err := s.us.GetUserByName(ctx, username)
	if err != nil {
		return nil, err
	}
	if moduleId > user.ModulesOpened {
		return nil, errModuleNotAllowed(moduleId)
	}
	module, err := s.ms.GetModuleByID(ctx, moduleId)
	if err != nil {
		return nil, err
	}
	return s.ts.CreateNewTest(ctx, user.Id, moduleId, module.Meta.TestQuestionsAmount)
}

func (s *service) GetTestByID(ctx context.Context, id string, hideAnswers bool) (*tests.Test, error) {
	return s.ts.GetTestByID(ctx, id, hideAnswers)
}

func (s *service) GetModuleForUser(ctx context.Context, user *users.UserDTO, moduleId int) (*module.ModuleDTO, error) {
	if moduleId > len(user.Modules) {
		return nil, errModuleNotAllowed(moduleId)
	}
	// load module
	mod, err := s.ms.GetModuleByID(ctx, moduleId)
	if err != nil {
		return nil, err
	}
	// list comments for module
	cmts, err := s.cs.ListCommentsByModule(ctx, user.Username, moduleId)
	if err != nil {
		return nil, err
	}
	comments := make(map[int][]*comment.CommentDTO)
	for _, c := range cmts {
		comments[c.PartId] = append(comments[c.PartId], c)
	}
	// need to convert string into template.HTML and add comments
	data := make([]module.ModulePartDTO, len(mod.Data))
	for i, p := range mod.Data {
		data[i] = module.ModulePartDTO{
			Id:       p.Id,
			Data:     template.HTML(p.Data),
			Comments: comments[p.Id],
			ModuleId: mod.Id,
		}
	}
	return &module.ModuleDTO{
		ModuleMeta:  mod.Meta,
		IsCompleted: !user.Modules[moduleId].CompletedAt.IsZero(),
		Data:        data,
	}, nil
}

func (s *service) ModulesPreview(ctx context.Context, user *users.UserDTO, amount int) ([]module.ModuleDTO, error) {
	modules, err := s.ms.ListOpenModulesMeta(ctx, amount)
	if err != nil {
		return nil, err
	}
	res := make([]module.ModuleDTO, len(modules))
	for i, m := range modules {
		res[i].ModuleMeta = m
		if !user.Modules[m.Id].CompletedAt.IsZero() {
			res[i].IsCompleted = true
		}
	}
	return res, nil
}
