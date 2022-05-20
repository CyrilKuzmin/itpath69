package service

import (
	"html/template"
	"time"

	"github.com/CyrilKuzmin/itpath69/internal/domain/users"
)

type UserDTO struct {
	*users.User
	ModulesTotal     int                       `json:"modules_total"`
	ModulesOpen      int                       `json:"modules_open"`
	ModulesCompleted int                       `json:"modules_completed"`
	Modules          map[int]ModuleProgressDTO `json:"modules"`
}

type CommentDTO struct {
	Id         string `json:"id" bson:"_id"` // UUID
	ModuleId   int    `json:"module_id"`
	PartId     int    `json:"part_id"`
	ModifiedAt string `json:"modified_at"`
	Text       string `json:"text"`
}

// ModulePartDTO is used for rendering
type ModulePartDTO struct {
	Id       int           `json:"id"`
	ModuleId int           `json:"module_id"` // comment form rendering bug
	Data     template.HTML `json:"data"`
	Comments []*CommentDTO `json:"comments"`
}

type ModuleMetaDTO struct {
	Id          int      `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Logo        string   `json:"logo"`
	// TestQuestionsAmount int      `json:"test_questions_amounts"`
	// TestPassThreshold   float64  `json:"test_pass_threshold"`
}

type ModuleDTO struct {
	*ModuleMetaDTO
	IsCompleted bool            `json:"completed"`
	Data        []ModulePartDTO `json:"data"`
}

type TestDTO struct {
	Id        string         `json:"id" bson:"_id"`
	UserId    string         `json:"user_id"`
	ModuleId  int            `json:"module_id"`
	CreatedAt string         `json:"created_at"`
	ExpiredAt string         `json:"expired_at"`
	Questions []*QuestionDTO `json:"questions"`
}

type QuestionDTO struct {
	Id           string      `json:"id"`
	QuestionText string      `json:"question_text"`
	ImageURL     string      `json:"image_url"`
	QuestionType int         `json:"question_type"`
	Answers      []AnswerDTO `json:"answers"`
}

type AnswerDTO struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

type TestResultDTO struct {
	Score    float32 `json:"score"`
	IsPassed bool    `json:"is_passed"`
}

type ModuleProgressDTO struct {
	OpenedAt    time.Time `json:"opened_at"`
	CompletedAt time.Time `json:"completed_at"`
}
