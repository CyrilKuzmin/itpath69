package service

import (
	"github.com/CyrilKuzmin/itpath69/internal/domain/comment"
	"github.com/CyrilKuzmin/itpath69/internal/domain/course"
	"github.com/CyrilKuzmin/itpath69/internal/domain/module"
	"github.com/CyrilKuzmin/itpath69/internal/domain/progress"
	"github.com/CyrilKuzmin/itpath69/internal/domain/tests"
)

const formatString = "02 Jan 06 15:04 MST"

func commentToDTO(c *comment.Comment) *CommentDTO {
	return &CommentDTO{
		Id:         c.ID,
		ModuleId:   c.ModuleID,
		PartId:     c.PartID,
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

func countModulesProgress(in []progress.ModuleProgress) (total, opened, completed int) {
	for _, p := range in {
		if !p.OpenedAt.IsZero() {
			opened++
		}
		if !p.CompletedAt.IsZero() {
			completed++
		}
	}
	return len(in), opened, completed
}

func convertModuleProgress(in progress.ModuleProgress) ModuleProgressDTO {
	return ModuleProgressDTO{
		OpenedAt:    in.OpenedAt,
		CompletedAt: in.CompletedAt,
	}
}

func convertModulesProgress(in []progress.ModuleProgress) map[int]ModuleProgressDTO {
	res := make(map[int]ModuleProgressDTO)
	for _, p := range in {
		if !p.OpenedAt.IsZero() {
			res[p.ModuleID] = ModuleProgressDTO{
				OpenedAt:    p.OpenedAt,
				CompletedAt: p.CompletedAt,
			}
		}
	}
	return res
}

func getCurrentUserStage(user *UserDTO, stages []course.Stage) course.Stage {
	total := 0
	for _, st := range stages {
		total += len(st.Modules)
		if user.ModulesOpen > total {
			// not this stage
			continue
		}
		return st
	}
	return stages[0]
}

func isCurrStageCompleted(user *UserDTO, stages []course.Stage) bool {
	total := 0
	for _, st := range stages {
		cmpl := 0
		total += len(st.Modules)
		for _, n := range st.Modules {
			if user.ModulesOpen > total {
				continue
			}
			if !user.Modules[n].CompletedAt.IsZero() {
				cmpl++
			}
			if cmpl >= st.ModulesToComplete {
				return true
			}
		}
	}
	return false
}
