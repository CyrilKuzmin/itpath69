package progress

import "time"

// ModuleProgress is a module for user
type ModuleProgress struct {
	ID          string    `json:"id" bson:"_id"`
	UserID      string    `json:"user_id"`
	CourseID    string    `json:"course_id"`
	ModuleID    int       `json:"module_id"`
	CreatedAt   time.Time `json:"created_at"`
	OpenedAt    time.Time `json:"opened_at"`
	CompletedAt time.Time `json:"completed_at"`
}
