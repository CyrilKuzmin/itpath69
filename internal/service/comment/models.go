package comment

import "time"

// Comment is a users comments for Part data
type Comment struct {
	ID         string    `json:"id" bson:"_id"` // UUID
	UserID     string    `json:"user_id"`
	CourseID   string    `json:"course_id"`
	ModuleID   int       `json:"module_id"`
	PartID     int       `json:"part_id"`
	ModifiedAt time.Time `json:"modified_at"`
	Text       string    `json:"text"`
}
