package postgres

import (
	"go-membership/app/models"
	"time"
)

type TokenRepository struct {
	connection
}

const TokensTable = "tokens"

func NewTokenRepository() *TokenRepository {
	return &TokenRepository{
		NewConnection("postgres"),
	}
}

func (tr *TokenRepository) Insert(token *models.Token) {
	tr.open()
	defer tr.close()
	tr.db.Create(token)
}

func (tr *TokenRepository) SelectUnExpiredByUserId(userId int) *models.Token {
	tr.open()
	defer tr.close()
	var token models.Token
	tr.db.Table(TokensTable).Where("user_id = ? ", userId).Where("expired_at > ?", time.Now()).Limit(1).Find(&token)
	return &token
}
