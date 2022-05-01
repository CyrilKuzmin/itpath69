// Models defines all used data models
package models

import "time"

// Stage contains several modules. Id defines the oder of stages as well
type Stage struct {
	Id      int      `json:"id,omitempty"`
	Modules []Module `json:"modules,omitempty"`
}

// Module contains metadata of module and list of its Parts (below)
type Module struct {
	Id          int      `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Parts       []Part   `json:"parts,omitempty"`
}

// Part contains all valuable Data. Probably it will be MD
type Part struct {
	Id   int    `json:"id,omitempty"`
	Data string `json:"data,omitempty"`
}

// User info and its opened stages
type User struct {
	Id           string       `json:"id,omitempty"`
	Username     string       `json:"username,omitempty"`
	PasswordHash string       `json:"password_hash,omitempty"`
	CreatedAt    time.Time    `json:"created_at,omitempty"`
	Stages       []UsersStage `json:"stages,omitempty"`
}

// UsersStage is a stage for user. It contains additional metadata. Id - uuid. Created from Stage
type UsersStage struct {
	Id          string           `json:"id,omitempty"`
	OrderNum    int              `json:"order_num,omitempty"`
	CreatedAt   time.Time        `json:"created_at,omitempty"`
	CompletedAt time.Time        `json:"completed_at,omitempty"`
	Modules     []ModuleProgress `json:"modules,omitempty"`
}

// ModuleProgress is a module for user. It contains additional metadata. Id - uuid. Created from Module
type ModuleProgress struct {
	Id          string               `json:"id,omitempty"`
	OrderNum    int                  `json:"order_num,omitempty"`
	CreatedAt   time.Time            `json:"created_at,omitempty"`
	CompletedAt time.Time            `json:"completed_at,omitempty"`
	Name        string               `json:"name,omitempty"`
	Description string               `json:"description,omitempty"`
	Tags        []string             `json:"tags,omitempty"`
	Parts       []ModuleProgressPart `json:"parts,omitempty"`
}

// ModuleProgressPart is a module's part for user. It contains additional metadata and comments. Id - uuid. Created from Part
type ModuleProgressPart struct {
	Id       string        `json:"id,omitempty"`
	IsDone   bool          `json:"is_done,omitempty"`
	Comments []PartComment `json:"comments,omitempty"`
	Part     *Part         `json:"part,omitempty"`
}

// PartComment is a users comments for Part data
type PartComment struct {
	Id        string    `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Text      string    `json:"text,omitempty"`
}
