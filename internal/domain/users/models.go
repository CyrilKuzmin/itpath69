package users

import "time"

// User info and its opened modules
type User struct {
	Id            string    `json:"id" bson:"_id"`
	Username      string    `json:"username"`
	PasswordHash  string    `json:"password_hash"`
	CreatedAt     time.Time `json:"created_at"`
	CurrentCourse string    `json:"current_course"` // ID of current course
	CurrentStage  int       `json:"current_stage"`
}
