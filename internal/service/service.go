package service

import (
	"context"
	"encoding/json"
	"errors"
	"html/template"
	"io"

	"github.com/CyrilKuzmin/itpath69/internal/service/comment"
	"github.com/CyrilKuzmin/itpath69/internal/service/course"
	"github.com/CyrilKuzmin/itpath69/internal/service/module"
	"github.com/CyrilKuzmin/itpath69/internal/service/progress"
	"github.com/CyrilKuzmin/itpath69/internal/service/tests"
	"github.com/CyrilKuzmin/itpath69/internal/service/users"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type Service interface {
	// User
	CreateUser(ctx context.Context, username, password string) (*UserDTO, error)
	GetUserByName(ctx context.Context, username string) (*UserDTO, error)
	CheckUserPassword(ctx context.Context, username, password string) error

	// Modules
	GetModuleForUser(ctx context.Context, user *UserDTO, moduleId int) (*ModuleDTO, error)
	ModulesPreview(ctx context.Context, user *UserDTO) ([]ModuleDTO, error)

	// Tests
	CreateTest(ctx context.Context, userId string, moduleId int) (*TestDTO, error)
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
	progress.Storage
	course.Storage
}

type service struct {
	us  users.Service
	ms  module.Service
	cs  comment.Service
	ts  tests.Service
	ps  progress.Service
	crs course.Service
}

func NewService(log *zap.Logger, s Storage) Service {
	// init echo server with its dependencies
	us := users.NewService(s, log)
	ms := module.NewService(s, log)
	cs := comment.NewService(s, log)
	ts := tests.NewService(s, log)
	ps := progress.NewService(s, log)
	crs := course.NewService(s, log)
	return &service{us, ms, cs, ts, ps, crs}
}

// simple functions for rendering

// GetUserByName
func (s *service) GetUserByName(ctx context.Context, username string) (*UserDTO, error) {
	user, err := s.us.GetUserByName(ctx, username)
	if err != nil {
		return nil, err
	}
	pr, err := s.ps.GetUserProgress(ctx, user.Id, user.CurrentCourse)
	if err != nil {
		return nil, err
	}
	total, opened, completed := countModulesProgress(pr)
	return &UserDTO{
		User:             user,
		ModulesOpen:      opened,
		ModulesCompleted: completed,
		ModulesTotal:     total,
		Modules:          convertModulesProgress(pr),
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
	if module > user.ModulesOpen {
		return nil, errModuleNotAllowed(module)
	}
	comments, err := s.cs.ListCommentsByModule(ctx, user.Id, user.CurrentCourse, module)
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
	if module > user.ModulesOpen {
		return nil, errModuleNotAllowed(module)
	}
	m, err := s.ms.GetModuleByID(ctx, user.CurrentCourse, module)
	if err != nil {
		return nil, err
	}
	if part > len(m.Data) {
		return nil, errPartNotExists(part)
	}
	newComment := comment.CreateCommentArgs{
		UserID:   user.Id,
		CourseID: user.CurrentCourse,
		ModuleID: module,
		PartID:   part,
		Text:     text,
	}
	c, err := s.cs.CreateComment(ctx, newComment)
	if err != nil {
		return nil, err
	}
	return commentToDTO(c), nil
}

func (s *service) UpdateComment(ctx context.Context, username, id, text string) (*CommentDTO, error) {
	user, err := s.GetUserByName(ctx, username)
	if err != nil {
		return nil, err
	}
	c, err := s.cs.UpdateComment(ctx, user.Id, id, text)
	if err != nil {
		return nil, err
	}
	return commentToDTO(c), err
}

func (s *service) DeleteCommentByID(ctx context.Context, username, id string) error {
	user, err := s.GetUserByName(ctx, username)
	if err != nil {
		return err
	}
	return s.cs.DeleteCommentByID(ctx, user.Id, id)
}

func (s *service) CreateUser(ctx context.Context, username, password string) (*UserDTO, error) {
	user, err := s.us.CreateUser(ctx, username, password)
	if err != nil {
		return nil, err
	}
	cc, err := s.crs.GetCourseByID(ctx, user.CurrentCourse)
	if err != nil {
		// if no default course in DB just return empty user
		if !errors.Is(err, mongo.ErrNoDocuments) {
			return nil, err
		}
		return &UserDTO{
			User:             user,
			ModulesOpen:      0,
			ModulesCompleted: 0,
			ModulesTotal:     0,
		}, nil
	}

	modulesToOpen := len(cc.Stages[0].Modules)
	pr, err := s.ps.CreateCourseProgress(ctx, user.Id, user.CurrentCourse, cc.TotalModules, modulesToOpen) // total amount of modules
	if err != nil {
		return nil, err
	}
	total, opened, completed := countModulesProgress(pr)
	return &UserDTO{
		User:             user,
		ModulesOpen:      opened,
		ModulesCompleted: completed,
		ModulesTotal:     total,
		Modules:          convertModulesProgress(pr),
	}, nil

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
	mod, err := s.ms.GetModuleByID(ctx, user.CurrentCourse, userTestData.ModuleId)
	var isPassed bool
	if float64(score) >= mod.Meta.TestPassThreshold {
		pr, err := s.ps.MarkModuleAsCompleted(ctx, user.Id, user.CurrentCourse, userTestData.ModuleId)
		if err != nil {
			return nil, err
		}
		user.ModulesCompleted++
		user.Modules[userTestData.ModuleId] = convertModuleProgress(*pr)
		isPassed = true
	}
	// add logic for opening new modules here
	cc, err := s.crs.GetCourseByID(ctx, user.CurrentCourse)
	if err != nil {
		return nil, err
	}
	if user.ModulesOpen < cc.TotalModules { // if we have closed modules
		// check stages
		if isCurrStageCompleted(user, cc.Stages) {
			modulesToOpen := len(getCurrentUserStage(user, cc.Stages).Modules)
			err = s.ps.OpenNewModules(ctx, user.Id, user.CurrentCourse, modulesToOpen)
			if err != nil {
				return nil, err
			}
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

func (s *service) CreateTest(ctx context.Context, username string, moduleId int) (*TestDTO, error) {
	user, err := s.GetUserByName(ctx, username)
	if err != nil {
		return nil, err
	}
	if moduleId > user.ModulesOpen {
		return nil, errModuleNotAllowed(moduleId)
	}
	cc, err := s.crs.GetCourseByID(ctx, user.CurrentCourse)
	if err != nil {
		return nil, err
	}
	module, err := s.ms.GetModuleByID(ctx, user.CurrentCourse, moduleId)
	if err != nil {
		return nil, err
	}
	newTestArgs := tests.CreateTestArgs{
		UserID:          user.Id,
		CourseID:        user.CurrentCourse,
		ModuleID:        moduleId,
		QuestionsAmount: module.Meta.TestQuestionsAmount,
		ExpirateAfter:   cc.TestsExpirationTime,
	}
	test, err := s.ts.CreateTest(ctx, newTestArgs)
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
	mod, err := s.ms.GetModuleByID(ctx, user.CurrentCourse, moduleId)
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

func (s *service) ModulesPreview(ctx context.Context, user *UserDTO) ([]ModuleDTO, error) {
	modules, err := s.ms.ListOpenModulesMeta(ctx, user.CurrentCourse, user.ModulesOpen)
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
