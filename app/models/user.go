package models

import (
	"time"
)

type User struct {
	ID            uint
	Name          string
	Email         string
	Password      string
	RememberToken string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Tokens        []Token `gorm:"foreignKey:user_id"`
}

func NewUser(name, email, password string) *User {
	return &User{
		Name:      name,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
	}
}
