package user

import "time"

// Separate Struct Idea : https://github.com/DremyGit/xwindy-lite/blob/master/models/user.go - Helps in Swagger Doc
// Struct Embedding : https://stackoverflow.com/a/27492025

// User Model
type User struct {
	UserInfoPayload
	Password string `gorm:"size:100;not null;" json:"password" validate:"required,min=8,max=100"`
}

// UserInfoPayload Struct
type UserInfoPayload struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Username  string    `gorm:"size:255;not null;unique" json:"username" validate:"required,min=4,max=30"`
	Email     string    `gorm:"size:100;not null;unique" json:"email" validate:"required,email"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// CreateUserPayload Struct
type CreateUserPayload struct {
	Username string `json:"username" validate:"required,min=4,max=30"`
	Email    string `json:"email" validate:"required,email"`
	Password string `gorm:"size:100;not null;" json:"password" validate:"required,min=8,max=100"`
}

// UpdateUserPayload Struct
type UpdateUserPayload struct {
	Password string `json:"password" validate:"required,min=8,max=100"`
}

// LoginPayload Struct
type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
