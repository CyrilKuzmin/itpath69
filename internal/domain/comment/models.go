package comment

import "time"

// Comment is a users comments for Part data
type Comment struct {
	Id        string    `json:"id,omitempty"`
	ModuleId  string    `json:"module_id,omitempty"`
	PartId    string    `json:"part_id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Text      string    `json:"text,omitempty"`
}
