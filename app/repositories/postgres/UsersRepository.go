package postgres

import (
	"go-membership/app/models"
)

type UsersRepository struct {
	connection
}

func NewUsersRepository() *UsersRepository {
	return &UsersRepository{
		NewConnection("postgres"),
	}
}

func (ur *UsersRepository) Insert(users *models.Users) {
	ur.open()
	defer ur.close()
	ur.db.Create(users)
}
