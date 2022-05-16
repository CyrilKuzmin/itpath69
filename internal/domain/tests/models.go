package tests

import "time"

type QType uint8

const (
	SingleAnswer QType = 1 << iota
	MultiChoose
)

type Test struct {
	Id        string      `json:"id" bson:"_id"`
	UserId    string      `json:"user_id"`
	ModuleId  int         `json:"module_id"`
	CreatedAt time.Time   `json:"created_at"`
	ExpiredAt time.Time   `json:"expired_at"`
	Questions []*Question `json:"questions"`
}

type Question struct {
	Id           string   `json:"id" bson:"_id"` // uuid generated
	ModuleId     int      `json:"module_id"`
	QuestionText string   `json:"question_text"`
	ImageURL     string   `json:"image_url"`
	QuestionType QType    `json:"question_type"`
	Answers      []Answer `json:"answers"`
}

type Answer struct {
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

type UserResult struct {
	Score    float32 `json:"score"`
	IsPassed bool    `json:"is_passed"`
}
