package service

import (
	"context"
	"encoding/json"
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
	GetUserByName(ctx context.Context, username string) (*UserDTO, error)
	CheckUserPassword(ctx context.Context, username, password string) error

	// Modules
	GetModuleForUser(ctx context.Context, user *UserDTO, moduleId int) (*ModuleDTO, error)
	ModulesPreview(ctx context.Context, user *UserDTO, amount int) ([]ModuleDTO, error)

	// Tests
	CreateNewTest(ctx context.Context, userId string, moduleId int) (*TestDTO, error)
	GetTestByID(ctx context.Context, id string, hideAnswers bool) (*TestDTO, error)
	ListTestsByUserID(ctx context.Context, userId string) ([]*TestDTO, error)
	ListTestsByUsername(ctx context.Context, username string) ([]*TestDTO, error)
	CheckTest(ctx context.Context, username string, userData io.Reader) (*TestResultDTO, error)

	//Comments
	CreateComment(ctx context.Context, username, text string, module, part int) (*CommentDTO, error)
	ListCommentsByModule(ctx context.Context, username string, module int) ([]*CommentDTO, error)
	UpdateComment(ctx context.Context, username, id, text string) (*CommentDTO, error)
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
func (s *service) GetUserByName(ctx context.Context, username string) (*UserDTO, error) {
	user, err := s.us.GetUserByName(ctx, username)
	if err != nil {
		return nil, err
	}
	opened := len(user.Modules)
	completed := 0
	for _, mp := range user.Modules {
		if !mp.CompletedAt.IsZero() {
			completed++
		}
	}
	return &UserDTO{
		User:             user,
		ModulesOpened:    opened,
		ModulesCompleted: completed,
		ModulesTotal:     8, // will be fixed when course will be implemented
	}, nil
}

func (s *service) ListTestsByUserID(ctx context.Context, userId string) ([]*TestDTO, error) {
	tests, err := s.ts.ListTestsByUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	res := make([]*TestDTO, len(tests))
	for i, t := range tests {
		res[i] = testToDTO(t)
	}
	return res, nil
}

func (s *service) ListTestsByUsername(ctx context.Context, username string) ([]*TestDTO, error) {
	user, err := s.us.GetUserByName(ctx, username)
	if err != nil {
		return nil, err
	}
	return s.ListTestsByUserID(ctx, user.Id)
}

func (s *service) ListCommentsByModule(ctx context.Context, username string, module int) ([]*CommentDTO, error) {
	// check if IDs are valid and allowed
	user, err := s.GetUserByName(ctx, username)
	if err != nil {
		return nil, err
	}
	if module > user.ModulesOpened {
		return nil, errModuleNotAllowed(module)
	}
	comments, err := s.cs.ListCommentsByModule(ctx, username, module)
	if err != nil {
		return nil, err
	}
	res := make([]*CommentDTO, len(comments))
	for i, c := range comments {
		res[i] = commentToDTO(c)
	}
	return res, nil
}

func (s *service) CreateComment(ctx context.Context, username, text string, module, part int) (*CommentDTO, error) {
	// check if IDs are valid and allowed
	user, err := s.GetUserByName(ctx, username)
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
	c, err := s.cs.CreateComment(ctx, username, text, module, part)
	if err != nil {
		return nil, err
	}
	return commentToDTO(c), nil
}

func (s *service) UpdateComment(ctx context.Context, username, id, text string) (*CommentDTO, error) {
	c, err := s.cs.UpdateComment(ctx, username, id, text)
	if err != nil {
		return nil, err
	}
	return commentToDTO(c), err
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

func (s *service) CheckTest(ctx context.Context, username string, userData io.Reader) (*TestResultDTO, error) {
	userTestData := &TestDTO{}
	err := json.NewDecoder(userData).Decode(&userTestData)
	if err != nil {
		return nil, err
	}
	user, err := s.GetUserByName(ctx, username)
	if err != nil {
		return nil, err
	}
	score, err := s.ts.CheckTest(ctx, userTestData.Id, qDtoToModel(userTestData.Questions))
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
	// add logic for opening new modules here
	if user.ModulesOpened < 8 { // fix it
		err = s.us.OpenNewModules(ctx, username, 4)
		if err != nil {
			return nil, err
		}
	}
	err = s.ts.MarkTestExpired(ctx, userTestData.Id)
	if err != nil {
		return nil, err
	}
	return &TestResultDTO{
		Score:    score,
		IsPassed: isPassed,
	}, nil
}

func (s *service) CreateNewTest(ctx context.Context, username string, moduleId int) (*TestDTO, error) {
	user, err := s.GetUserByName(ctx, username)
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
	test, err := s.ts.CreateNewTest(ctx, user.Id, moduleId, module.Meta.TestQuestionsAmount)
	if err != nil {
		return nil, err
	}
	return testToDTO(test), err
}

func (s *service) GetTestByID(ctx context.Context, id string, hideAnswers bool) (*TestDTO, error) {
	test, err := s.ts.GetTestByID(ctx, id, hideAnswers)
	if err != nil {
		return nil, err
	}
	return testToDTO(test), err
}

func (s *service) GetModuleForUser(ctx context.Context, user *UserDTO, moduleId int) (*ModuleDTO, error) {
	if moduleId > len(user.Modules) {
		return nil, errModuleNotAllowed(moduleId)
	}
	// load module
	mod, err := s.ms.GetModuleByID(ctx, moduleId)
	if err != nil {
		return nil, err
	}
	// list comments for module
	cmts, err := s.ListCommentsByModule(ctx, user.Username, moduleId)
	if err != nil {
		return nil, err
	}
	comments := make(map[int][]*CommentDTO)
	for _, c := range cmts {
		comments[c.PartId] = append(comments[c.PartId], c)
	}
	// need to convert string into template.HTML and add comments
	data := make([]ModulePartDTO, len(mod.Data))
	for i, p := range mod.Data {
		data[i] = ModulePartDTO{
			Id:       p.Id,
			Data:     template.HTML(p.Data),
			Comments: comments[p.Id],
			ModuleId: mod.Id,
		}
	}
	return &ModuleDTO{
		ModuleMetaDTO: moduleMetaToDTO(&mod.Meta),
		IsCompleted:   !user.Modules[moduleId].CompletedAt.IsZero(),
		Data:          data,
	}, nil
}

func (s *service) ModulesPreview(ctx context.Context, user *UserDTO, amount int) ([]ModuleDTO, error) {
	modules, err := s.ms.ListOpenModulesMeta(ctx, amount)
	if err != nil {
		return nil, err
	}
	res := make([]ModuleDTO, len(modules))
	for i, m := range modules {
		res[i].ModuleMetaDTO = moduleMetaToDTO(&m)
		if !user.Modules[m.Id].CompletedAt.IsZero() {
			res[i].IsCompleted = true
		}
	}
	return res, nil
}
