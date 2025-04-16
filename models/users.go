package models

import "time"

type Users struct {
	Username     string    `json:"username" gorm:"unique;not null"`
	Email        string    `json:"email" gorm:"unique;not null"`
	PasswordHash string    `json:"password_hash" gorm:"not null"`
	AvatarUrl    string    `json:"avatar_url" gorm:"default 'https://www.gravatar.com/avatar/"`
	Bio          string    `json:"bio"`
	Role         string    `json:"role" gorm:"default 'viewer'"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
