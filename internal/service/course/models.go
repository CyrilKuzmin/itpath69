package course

import "time"

type Course struct {
	ID                  string        `json:"id" bson:"_id"`
	IsPublished         bool          `json:"is_published"`
	IsPrivate           bool          `json:"is_private"`
	Name                string        `json:"name"`
	Owners              []string      `json:"owners"`
	Stages              []Stage       `json:"stages"`
	TotalModules        int           `json:"total_modules"`
	TestsExpirationTime time.Duration `json:"tests_expiration_time"`
}

type Stage struct {
	ID                int   `json:"id"`
	ModulesToComplete int   `json:"modules_to_complete"`
	Modules           []int `json:"modules"` // IDs of modules
}
