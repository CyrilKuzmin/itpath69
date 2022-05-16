package tests

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const DefaultQuestionsAmount = 3
const DefaultExpiraionTime = 24 * time.Hour

type Service interface {
	CreateNewTest(ctx context.Context, userId string, moduleId, amount int) (*Test, error)
	GetTestByID(ctx context.Context, id string, hideAnswers bool) (*Test, error)
	ListTestsByUser(ctx context.Context, userId string) ([]*Test, error)
	CheckTest(ctx context.Context, id string, userAnswers []*Question) (float32, error)
	MarkTestExpired(ctx context.Context, id string) error
	SaveQuestions(ctx context.Context, qs []Question) error
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

func (s *service) CreateNewTest(ctx context.Context, userId string, moduleId, amount int) (*Test, error) {
	// moduleId == 0 - specific case, need to get questions for all opened modules
	qs, err := s.storage.GetModuleQuestions(ctx, moduleId, amount)
	if err != nil {
		return nil, err
	}
	test := Test{
		Id:        uuid.NewString(),
		UserId:    userId,
		ModuleId:  moduleId,
		CreatedAt: time.Now(),
		ExpiredAt: time.Now().Add(DefaultExpiraionTime),
		Questions: qs,
	}
	err = s.storage.SaveTest(ctx, &test)
	if err != nil {
		return nil, err
	}
	for qId := range test.Questions {
		for aId := range test.Questions[qId].Answers {
			test.Questions[qId].Answers[aId].IsCorrect = false
		}
	}
	fmt.Println("new test ID", test.Id)
	return &test, nil
}

func (s *service) GetTestByID(ctx context.Context, id string, hideAnswers bool) (*Test, error) {
	test, err := s.storage.GetTestByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if hideAnswers {
		for qId := range test.Questions {
			for aId := range test.Questions[qId].Answers {
				test.Questions[qId].Answers[aId].IsCorrect = false
			}
		}
	}
	return test, nil
}

func (s *service) ListTestsByUser(ctx context.Context, userId string) ([]*Test, error) {
	return s.storage.ListTestsByUser(ctx, userId)
}

func (s *service) CheckTest(ctx context.Context, id string, userAnswers []*Question) (float32, error) {
	t, err := s.storage.GetTestByID(ctx, id)
	if err != nil {
		return 0, err
	}
	res := float32(0)
	for _, tq := range t.Questions {
		for _, ua := range userAnswers {
			if tq.Id == ua.Id {
				res += checkQuestion(tq, ua)
			}
		}
	}
	score := res / float32(len(t.Questions))
	return score, nil
}

func (s *service) MarkTestExpired(ctx context.Context, id string) error {
	return s.storage.MarkTestExpired(ctx, id)
}

func (s *service) SaveQuestions(ctx context.Context, qs []Question) error {
	return s.storage.SaveQuestions(ctx, qs)
}

func checkQuestion(orig, user *Question) float32 {
	switch orig.QuestionType {
	case SingleAnswer:
		return checkSingleAnswer(orig, user)
	case MultiChoose:
		return checkMultiChoose(orig, user)
	default:
		return 0
	}
}

func checkSingleAnswer(orig, user *Question) float32 {
	for _, a := range orig.Answers {
		for _, b := range user.Answers {
			if a.Text == b.Text {
				if a.IsCorrect == b.IsCorrect {
					return 1.0
				} else {
					return 0.0
				}
			}
		}
	}
	return 0.0
}

func checkMultiChoose(orig, user *Question) float32 {
	correctAnswers := 0
	for _, a := range orig.Answers {
		if a.IsCorrect {
			correctAnswers++
		}
	}
	userCorrect := 0
	for _, a := range orig.Answers {
		for _, b := range user.Answers {
			if a.Text == b.Text {
				if a.IsCorrect == true && b.IsCorrect == true {
					userCorrect += 1
				}
			}
		}
	}
	res := (float32(userCorrect) / float32(correctAnswers))
	return res
}
