package module

import (
	"html/template"

	"github.com/CyrilKuzmin/itpath69/internal/domain/comment"
)

// Module contains metadata of module and list of its Parts (below)
type Module struct {
	Id   int        `json:"id" bson:"_id"`
	Meta ModuleMeta `json:"meta"`
	Data []Part     `json:"parts"`
}

type ModuleMeta struct {
	Id                  int      `json:"id"`
	Name                string   `json:"name"`
	Description         string   `json:"description"`
	Tags                []string `json:"tags"`
	Logo                string   `json:"logo"`
	TestQuestionsAmount int      `json:"test_questions_amounts"`
	TestPassThreshold   float64  `json:"test_pass_threshold"`
}

// Part contains all valuable Data. Probably it will be MD or HTML
type Part struct {
	Id   int    `json:"id"`
	Data string `json:"data"`
}

// ModulePartDTO is used for rendering
type ModulePartDTO struct {
	Id       int
	ModuleId int // comment form rendering bug
	Data     template.HTML
	Comments []*comment.CommentDTO
}

type ModuleDTO struct {
	ModuleMeta
	IsCompleted bool `json:"completed"`
	Data        []ModulePartDTO
}
