package models

import "time"

type Token struct {
	ID        uint
	UserId    uint
	Token     string
	ExpiredAt time.Time
	CreatedAt time.Time
}

func NewToken(user User, token string, expiredAt time.Time) *Token {
	return &Token{
		UserId:    user.ID,
		Token:     token,
		ExpiredAt: expiredAt,
		CreatedAt: time.Now(),
	}
}
