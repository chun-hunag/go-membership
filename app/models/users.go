package models

import (
	"time"
)

type Users struct {
	ID            uint
	Name          string
	Email         string
	Password      string
	RememberToken string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewUsers(name, email, password string) *Users {
	return &Users{
		Name:      name,
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
	}
}
