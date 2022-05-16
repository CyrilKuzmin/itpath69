package course

import "time"

type Course struct {
	ID                  string   `json:"id"`
	Name                string   `json:"name"`
	Owners              []string `json:"owners"`
	IsPrivate           bool     `json:"is_private"`
	Stages              []Stage  `json:"stages"`
	TestsExpirationTime time.Duration
}

type Stage struct {
	ID                int   `json:"id"`
	ModulesToComplete int   `json:"modules_to_complete"`
	Modules           []int `json:"modules"` // IDs of modules
}
