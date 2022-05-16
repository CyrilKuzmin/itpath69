package comment

import "time"

// Comment is a users comments for Part data
type Comment struct {
	Id         string    `json:"id" bson:"_id"` // UUID
	User       string    `json:"user"`
	ModuleId   int       `json:"module_id"`
	PartId     int       `json:"part_id"`
	ModifiedAt time.Time `json:"modified_at"`
	Text       string    `json:"text"`
}
