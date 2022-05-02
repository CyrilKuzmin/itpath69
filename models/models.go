// Models defines all used data models
package models

import "time"

// Stage contains several modules. Id defines the oder of stages as well
type Stage struct {
	Id      int      `json:"id,omitempty" bson:"_id"`
	Modules []Module `json:"modules,omitempty"`
}

// Module contains metadata of module and list of its Parts (below)
type Module struct {
	Id   int        `json:"id,omitempty" bson:"_id"`
	Meta ModuleMeta `json:"meta,omitempty"`
	Data []Part     `json:"parts,omitempty"`
}

type ModuleMeta struct {
	Id          int      `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Logo        string   `json:"logo,omitempty"`
	Completed   bool     `json:"completed"` // will be fullfiled for rendering
}

// Part contains all valuable Data. Probably it will be MD or HTML
type Part struct {
	Id       int       `json:"id,omitempty"`
	Data     string    `json:"data,omitempty"`
	Comments []Comment `json:"comments"` // will be fullfiled for rendering
}

// User info and its opened modules
type User struct {
	Id           string                 `json:"id,omitempty"`
	Username     string                 `json:"username,omitempty"`
	PasswordHash string                 `json:"password_hash,omitempty"`
	CreatedAt    time.Time              `json:"created_at,omitempty"`
	Modules      map[int]ModuleProgress `json:"modules,omitempty"`
}

// ModuleProgress is a module for user. It contains additional metadata. Id - uuid. Created from Module
type ModuleProgress struct {
	CreatedAt   time.Time `json:"created_at,omitempty"`
	CompletedAt time.Time `json:"completed_at,omitempty"`
}

// Comment is a users comments for Part data
type Comment struct {
	Id        string    `json:"id,omitempty"`
	ModuleId  string    `json:"module_id,omitempty"`
	PartId    string    `json:"part_id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Text      string    `json:"text,omitempty"`
}
