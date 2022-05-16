package users

import "time"

// User info and its opened modules
type User struct {
	Id           string                 `json:"id" bson:"_id"`
	Username     string                 `json:"username"`
	PasswordHash string                 `json:"password_hash"`
	CreatedAt    time.Time              `json:"created_at"`
	Modules      map[int]ModuleProgress `json:"modules"`
}

// ModuleProgress is a module for user. It contains additional metadata. Id - uuid. Created from Module
type ModuleProgress struct {
	CreatedAt   time.Time `json:"created_at"`
	CompletedAt time.Time `json:"completed_at"`
}

type UserDTO struct {
	*User
	ModulesTotal     int `json:"modules_total"`
	ModulesOpened    int `json:"modules_opened"`
	ModulesCompleted int `json:"modules_completed"`
}
