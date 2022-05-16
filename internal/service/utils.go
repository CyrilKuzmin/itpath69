package service

import (
	"github.com/CyrilKuzmin/itpath69/internal/domain/comment"
	"github.com/CyrilKuzmin/itpath69/internal/domain/module"
	"github.com/CyrilKuzmin/itpath69/internal/domain/tests"
)

const formatString = "02 Jan 06 15:04 MST"

func commentToDTO(c *comment.Comment) *CommentDTO {
	return &CommentDTO{
		Id:         c.Id,
		ModuleId:   c.ModuleId,
		PartId:     c.PartId,
		ModifiedAt: c.ModifiedAt.Format(formatString),
		Text:       c.Text,
	}
}

func moduleMetaToDTO(meta *module.ModuleMeta) *ModuleMetaDTO {
	return &ModuleMetaDTO{
		Id:          meta.Id,
		Name:        meta.Name,
		Description: meta.Description,
		Tags:        meta.Tags,
		Logo:        meta.Logo,
	}
}

func testToDTO(test *tests.Test) *TestDTO {
	return &TestDTO{
		Id:        test.Id,
		UserId:    test.UserId,
		ModuleId:  test.ModuleId,
		CreatedAt: test.CreatedAt.Format(formatString),
		ExpiredAt: test.ExpiredAt.Format(formatString),
		Questions: questionToDTO(test.Questions),
	}
}

func questionToDTO(in []*tests.Question) []*QuestionDTO {
	res := make([]*QuestionDTO, len(in))
	for i, a := range in {
		res[i] = &QuestionDTO{
			Id:           a.Id,
			QuestionText: a.QuestionText,
			ImageURL:     a.ImageURL,
			QuestionType: int(a.QuestionType),
			Answers:      answersToDTO(a.Answers),
		}
	}
	return res
}

func answersToDTO(in []tests.Answer) []AnswerDTO {
	res := make([]AnswerDTO, len(in))
	for i, a := range in {
		res[i] = AnswerDTO{
			Text:      a.Text,
			IsCorrect: a.IsCorrect,
		}
	}
	return res
}

func qDtoToModel(in []*QuestionDTO) []*tests.Question {
	res := make([]*tests.Question, len(in))
	for i, q := range in {
		res[i] = &tests.Question{
			Id:           q.Id,
			QuestionType: tests.QType(q.QuestionType),
			QuestionText: q.QuestionText,
			Answers:      aDtoToModel(q.Answers),
		}
	}
	return res
}

func aDtoToModel(in []AnswerDTO) []tests.Answer {
	res := make([]tests.Answer, len(in))
	for i, a := range in {
		res[i] = tests.Answer{
			Text:      a.Text,
			IsCorrect: a.IsCorrect,
		}
	}
	return res
}